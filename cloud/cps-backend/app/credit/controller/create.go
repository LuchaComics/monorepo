package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

type CreditCreateRequest struct {
	UserID           primitive.ObjectID `bson:"user_id" json:"user_id"`
	BusinessFunction int8               `bson:"business_function" json:"business_function"`
	OfferID          primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	NumberOfCredits  int                `bson:"number_of_credits" json:"number_of_credits"`
}

func (impl *CreditControllerImpl) Create(ctx context.Context, req *CreditCreateRequest) error {
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
		// Extract from our session the following data.
		urole := sessCtx.Value(constants.SessionUserRole).(int8)
		// uid := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		// uname := sessCtx.Value(constants.SessionUserName).(string)
		// oid := sessCtx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
		// oname := sessCtx.Value(constants.SessionUserStoreName).(string)

		switch urole { // Security.
		case u_d.UserRoleRoot:
			impl.Logger.Debug("access granted")
		default:
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Get the user.
		user, err := impl.UserStorer.GetByID(sessCtx, req.UserID)
		if err != nil {
			impl.Logger.Error("database fetch error", slog.Any("error", err))
			return nil, err
		}
		if user == nil {
			impl.Logger.Error("user does not exist", slog.Any("user_id", req.UserID))
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "user does not exist for this value")
		}

		// Get the office.
		offer, err := impl.OfferStorer.GetByID(sessCtx, req.OfferID)
		if err != nil {
			impl.Logger.Error("database fetch error", slog.Any("error", err))
			return nil, err
		}
		if offer == nil {
			impl.Logger.Error("offer does not exist", slog.Any("offer_id", req.OfferID))
			return nil, httperror.NewForBadRequestWithSingleField("offer_id", "offer does not exist for this value")
		}

		// Add defaults / meta / etimpl.
		m := &s_d.Credit{
			StoreID:          user.StoreID,
			StoreName:        user.StoreName,
			StoreTimezone:    user.StoreTimezone,
			CreatedAt:        time.Now(),
			ModifiedAt:       time.Now(),
			Status:           s_d.StatusActive,
			BusinessFunction: req.BusinessFunction,
			UserID:           user.ID,
			UserName:         user.Name,
			UserLexicalName:  user.LexicalName,
			OfferID:          offer.ID,
			OfferName:        offer.Name,
			OfferServiceType: offer.ServiceType,
		}

		for no := 0; no < req.NumberOfCredits; no++ {
			m.ID = primitive.NewObjectID() // Generate new ID every time.

			impl.Logger.Debug("attaching credit",
				slog.Any("id", m.ID),
				slog.Int("no", no))

			// Save to our database.
			if err := impl.CreditStorer.Create(sessCtx, m); err != nil {
				impl.Logger.Error("database create error", slog.Any("error", err))
				return nil, err
			}
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
