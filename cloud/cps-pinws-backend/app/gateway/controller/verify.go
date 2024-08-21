package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/mongo"

	gateway_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/gateway/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Verify(ctx context.Context, code string) (*gateway_s.VerifyResponseIDO, error) {
	impl.Kmutex.Lock(code)
	defer func() {
		impl.Kmutex.Unlock(code)
	}()

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
		res := &gateway_s.VerifyResponseIDO{}

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByVerificationCode(sessCtx, code)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("code", "does not exist")
		}

		//TODO: Handle expiry dates.

		// Verify the user.
		u.WasEmailVerified = true
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("update error", slog.Any("err", err))
			return nil, err
		}

		//
		// Send notification based on user role
		//

		switch u.Role {
		case user_s.UserRoleRetailer:
			{
				o, err := impl.TenantStorer.GetByID(sessCtx, u.TenantID)
				if err != nil {
					impl.Logger.Error("database error", slog.Any("err", err))
					return nil, err
				}
				if o == nil {
					impl.Logger.Warn("store does not exist validation error")
					return nil, httperror.NewForBadRequestWithSingleField("code", "does not exist")
				}

				// If nothing was returned then proceed to send the notification and then
				// track that it was already sent.
				if o.PendingReviewEmailSent == false {
					// Send email to root staff to inform them of a pending review.
					res, err := impl.UserStorer.ListAllRootStaff(sessCtx)
					if err != nil {
						impl.Logger.Error("database error", slog.Any("err", err))
						return nil, err
					}
					var emails []string
					for _, rootUser := range res.Results {
						emails = append(emails, rootUser.Email)
					}
					if err := impl.TemplatedEmailer.SendNewStoreEmailToStaff(emails, u.TenantID.Hex()); err != nil {
						impl.Logger.Error("failed sending verification email with error", slog.Any("err", err))
						return nil, err
					}

					// Keep track of this verification sent so we don't send it again.
					o.PendingReviewEmailSent = true
					if err := impl.TenantStorer.UpdateByID(sessCtx, o); err != nil {
						impl.Logger.Error("database error", slog.Any("err", err))
						return nil, err
					}
				}
				res.Message = "Thank you for verifying. Your application has been forward to CPS_PINWS staff, and will be processed within 1 business day."
				impl.Logger.Debug("business user verified")
				break
			}
		case user_s.UserRoleCustomer:
			{
				res.Message = "Thank you for verifying. You may log in now to get started!"
				impl.Logger.Debug("customer user verified")
				break
			}
		default:
			{
				res.Message = "Thank you for verifying. You may log in now to get started!"
				impl.Logger.Debug("unknown user verified")
				break
			}
		}
		res.UserRole = u.Role

		return res, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	impl.Logger.Debug("user verified", slog.Any("response", res))

	return res.(*gateway_s.VerifyResponseIDO), nil
}
