package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
)

func (impl *ComicSubmissionControllerImpl) SetCustomer(ctx context.Context, submissionID primitive.ObjectID, customerID primitive.ObjectID) (*submission_s.ComicSubmission, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Fetch the original submission.
		os, err := impl.ComicSubmissionStorer.GetByID(sessCtx, submissionID)
		if err != nil {
			impl.Logger.Error("database get by id error",
				slog.Any("customerID", customerID),
				slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			return nil, nil
		}

		if !customerID.IsZero() {
			customer, err := impl.UserStorer.GetByID(sessCtx, customerID)
			if err != nil {
				impl.Logger.Error("get customer user error", slog.Any("error", err))
				return nil, err
			}

			// Modify our original submission.
			os.ModifiedAt = time.Now()
			os.CustomerID = customer.ID
			os.CustomerFirstName = customer.FirstName
			os.CustomerLastName = customer.LastName
		}

		// Save to the database the modified submission.
		if err := impl.ComicSubmissionStorer.UpdateByID(sessCtx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		return os, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*submission_s.ComicSubmission), nil
}
