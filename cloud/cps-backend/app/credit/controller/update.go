package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *CreditControllerImpl) UpdateByID(ctx context.Context, ns *domain.Credit) (*domain.Credit, error) {
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
		oid := sessCtx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
		oname := sessCtx.Value(constants.SessionUserStoreName).(string)
		otz := sessCtx.Value(constants.SessionUserStoreTimezone).(string)

		switch urole { // Security.
		case u_d.UserRoleRoot:
			impl.Logger.Debug("access granted")
		default:
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Fetch the original store.
		os, err := impl.CreditStorer.GetByID(sessCtx, ns.ID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			return nil, httperror.NewForBadRequestWithSingleField("id", "credit does not exist")
		}

		// Get the user.
		user, err := impl.UserStorer.GetByID(sessCtx, ns.UserID)
		if err != nil {
			impl.Logger.Error("database fetch error", slog.Any("error", err))
			return nil, err
		}
		if user == nil {
			impl.Logger.Error("user does not exist", slog.Any("user_id", ns.UserID))
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "user does not exist for this value")
		}

		// Get the office.
		offer, err := impl.OfferStorer.GetByID(sessCtx, ns.OfferID)
		if err != nil {
			impl.Logger.Error("database fetch error", slog.Any("error", err))
			return nil, err
		}
		if offer == nil {
			impl.Logger.Error("offer does not exist", slog.Any("offer_id", ns.OfferID))
			return nil, httperror.NewForBadRequestWithSingleField("offer_id", "offer does not exist for this value")
		}

		// Modify our original store.
		os.StoreID = oid
		os.StoreName = oname
		os.StoreTimezone = otz
		os.ModifiedAt = time.Now()
		os.Status = ns.Status
		os.UserID = user.ID
		os.UserName = user.Name
		os.UserLexicalName = user.LexicalName
		os.OfferID = offer.ID
		os.OfferName = offer.Name
		os.OfferServiceType = offer.ServiceType
		os.ModifiedAt = time.Now()

		// Save to the database the modified store.
		if err := impl.CreditStorer.UpdateByID(sessCtx, os); err != nil {
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

	return res.(*domain.Credit), nil
}
