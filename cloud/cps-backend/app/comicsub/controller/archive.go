package controller

import (
	"context"
	"log/slog"
	"time"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl *ComicSubmissionControllerImpl) ArchiveByID(ctx context.Context, id primitive.ObjectID) (*domain.ComicSubmission, error) {
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
		os, err := impl.ComicSubmissionStorer.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			return nil, nil
		}

		// Modify our original submission.
		os.ModifiedAt = time.Now()
		os.Status = domain.StatusArchived

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

	return res.(*domain.ComicSubmission), nil
}
