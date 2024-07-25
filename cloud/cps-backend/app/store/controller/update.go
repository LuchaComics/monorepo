package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *StoreControllerImpl) UpdateByID(ctx context.Context, ns *domain.Store) (*domain.Store, error) {
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
		// Fetch the original store.
		os, err := impl.StoreStorer.GetByID(ctx, ns.ID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			impl.Logger.Error("store does not exist error",
				slog.Any("store_id", ns.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "store does not exist")
		}

		// Extract from our session the following data.
		userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
		userStoreID := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
		userRole := ctx.Value(constants.SessionUserRole).(int8)
		userName := ctx.Value(constants.SessionUserName).(string)

		// If user is not administrator nor belongs to the store then error.
		if userRole != user_d.UserRoleRoot && os.ID != userStoreID {
			impl.Logger.Error("authenticated user is not staff role nor belongs to the store error",
				slog.Any("userRole", userRole),
				slog.Any("userStoreID", userStoreID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this store")
		}

		// Store previous value before we do any modifications so we can run logic
		// based on change of certain values.
		previousStatus := os.Status

		// Modify our original store.
		os.ModifiedAt = time.Now()
		os.ModifiedByUserID = userID
		os.ModifiedByUserName = userName
		os.Type = ns.Type
		os.Status = ns.Status
		os.Name = ns.Name
		os.WebsiteURL = ns.WebsiteURL
		os.EstimatedSubmissionsPerMonth = ns.EstimatedSubmissionsPerMonth
		os.HasOtherGradingService = ns.HasOtherGradingService
		os.OtherGradingServiceName = ns.OtherGradingServiceName
		os.RequestWelcomePackage = ns.RequestWelcomePackage
		os.HowLongStoreOperating = ns.HowLongStoreOperating
		os.GradingComicsExperience = ns.GradingComicsExperience
		os.RetailPartnershipReason = ns.RetailPartnershipReason
		os.CPSPartnershipReason = ns.CPSPartnershipReason
		os.Level = ns.Level
		os.SpecialCollection = ns.SpecialCollection

		// Save to the database the modified store.
		if err := impl.StoreStorer.UpdateByID(ctx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		// Send notifications in the background.
		if previousStatus != ns.Status && ns.Status == s_d.StoreActiveStatus {
			impl.Logger.Debug("store became active, sending email to retailer staff")
			go func(m *s_d.Store) {
				res, err := impl.UserStorer.ListAllRetailerStaffForStoreID(context.Background(), m.ID)
				if err != nil {
					impl.Logger.Error("list store error", slog.Any("error", err))
					return
				}
				var retailerEmails []string
				for _, u := range res.Results {
					retailerEmails = append(retailerEmails, u.Email)
				}
				if err := impl.TemplatedEmailer.SendRetailerStoreActiveEmailToRetailers(retailerEmails, m.Name); err != nil {
					impl.Logger.Error("failed sending templated error", slog.Any("error", err))
					return
				}
			}(os)
		}

		////
		//// Transaction successfully completed.
		////

		return os, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	org := res.(*domain.Store)

	// Run the following update tasks in the background regardless of this
	// function completing.
	go func(o *domain.Store) {
		impl.updateRelatedUsersInBackground(o)
	}(org)
	go func(o *domain.Store) {
		impl.updateRelatedComicSubmissionsInBackground(o)
	}(org)
	go func(o *domain.Store) {
		impl.updateRelatedAttachmentsInBackground(o)
	}(org)
	go func(o *domain.Store) {
		impl.updateRelatedCreditsInBackground(o)
	}(org)
	go func(o *domain.Store) {
		impl.updateRelatedReceiptsInBackground(o)
	}(org)
	go func(o *domain.Store) {
		impl.updateRelateUserPurchasesInBackground(o)
	}(org)

	return org, nil
}
