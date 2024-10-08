package controller

import (
	"context"
	"time"

	"log/slog"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl *OfferControllerImpl) UpdateByID(ctx context.Context, ns *domain.Offer) (*domain.Offer, error) {
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
		// Extract from our session the following data.
		urole := sessCtx.Value(constants.SessionUserRole).(int8)
		// uid := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		// uname := sessCtx.Value(constants.SessionUserName).(string)

		switch urole { // Security.
		case u_d.UserRoleRoot:
			impl.Logger.Debug("access granted")
		default:
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Fetch the original store.
		os, err := impl.OfferStorer.GetByID(sessCtx, ns.ID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			return nil, httperror.NewForBadRequestWithSingleField("id", "receipt type does not exist")
		}

		os.ModifiedAt = time.Now()
		os.Status = ns.Status
		os.BusinessFunction = ns.BusinessFunction
		os.ServiceType = ns.ServiceType

		// Save to the database the modified store.
		if err := impl.OfferStorer.UpdateByID(sessCtx, os); err != nil {
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

	return res.(*domain.Offer), nil
}
