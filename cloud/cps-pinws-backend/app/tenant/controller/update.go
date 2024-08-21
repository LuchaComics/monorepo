package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (impl *TenantControllerImpl) UpdateByID(ctx context.Context, ns *domain.Tenant) (*domain.Tenant, error) {
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
		// Fetch the original tenant.
		os, err := impl.TenantStorer.GetByID(ctx, ns.ID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if os == nil {
			impl.Logger.Error("tenant does not exist error",
				slog.Any("tenant_id", ns.ID))
			return nil, httperror.NewForBadRequestWithSingleField("message", "tenant does not exist")
		}

		// Extract from our session the following data.
		userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
		userTenantID := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
		userRole := ctx.Value(constants.SessionUserRole).(int8)
		userName := ctx.Value(constants.SessionUserName).(string)

		// If user is not administrator nor belongs to the tenant then error.
		if userRole != user_d.UserRoleRoot && os.ID != userTenantID {
			impl.Logger.Error("authenticated user is not staff role nor belongs to the tenant error",
				slog.Any("userRole", userRole),
				slog.Any("userTenantID", userTenantID))
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not belong to this tenant")
		}

		// Tenant previous value before we do any modifications so we can run logic
		// based on change of certain values.
		previousStatus := os.Status

		// Modify our original tenant.
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
		os.HowLongTenantOperating = ns.HowLongTenantOperating
		os.GradingComicsExperience = ns.GradingComicsExperience
		os.RetailPartnershipReason = ns.RetailPartnershipReason
		os.CPS_PINWSPartnershipReason = ns.CPS_PINWSPartnershipReason
		os.Level = ns.Level
		os.SpecialCollection = ns.SpecialCollection

		// Save to the database the modified tenant.
		if err := impl.TenantStorer.UpdateByID(ctx, os); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		// Send notifications in the background.
		if previousStatus != ns.Status && ns.Status == s_d.TenantActiveStatus {
			impl.Logger.Debug("tenant became active, sending email to retailer staff")
			go func(m *s_d.Tenant) {
				res, err := impl.UserStorer.ListAllRetailerStaffForTenantID(context.Background(), m.ID)
				if err != nil {
					impl.Logger.Error("list tenant error", slog.Any("error", err))
					return
				}
				var retailerEmails []string
				for _, u := range res.Results {
					retailerEmails = append(retailerEmails, u.Email)
				}
				// if err := impl.TemplatedEmailer.SendRetailerTenantActiveEmailToRetailers(retailerEmails, m.Name); err != nil {
				// 	impl.Logger.Error("failed sending templated error", slog.Any("error", err))
				// 	return
				// }
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

	org := res.(*domain.Tenant)

	// Run the following update tasks in the background regardless of this
	// function completing.
	go func(o *domain.Tenant) {
		impl.updateRelatedUsersInBackground(o)
	}(org)
	go func(o *domain.Tenant) {
		impl.updateRelatedAttachmentsInBackground(o)
	}(org)

	return org, nil
}
