package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *OfferControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		d, err := impl.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if d == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, httperror.NewForBadRequestWithSingleField("id", "receipt type does not exist")
		}
		d.Status = s_d.StatusArchived
		d.ModifiedAt = time.Now()

		// Save to the database the modified store.
		if err := impl.OfferStorer.UpdateByID(sessCtx, d); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		return nil, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}
