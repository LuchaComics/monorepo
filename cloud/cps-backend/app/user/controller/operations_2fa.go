package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	u_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

type UserOperationChangeTwoFactorAuthenticationRequest struct {
	UserID     primitive.ObjectID `bson:"user_id" json:"user_id"`
	OTPEnabled bool               `bson:"otp_enabled" json:"otp_enabled"`
}

func (impl *UserControllerImpl) validateOperationChangeTwoFactorAuthenticationRequest(ctx context.Context, dirtyData *UserOperationChangeTwoFactorAuthenticationRequest) error {
	e := make(map[string]string)

	if dirtyData.UserID.IsZero() {
		e["user_id"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *UserControllerImpl) ChangeTwoFactorAuthentication(ctx context.Context, req *UserOperationChangeTwoFactorAuthenticationRequest) error {
	//
	// Get variables from our user authenticated session.
	//

	tid, _ := ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID)
	role, _ := ctx.Value(constants.SessionUserRole).(int8)
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName, _ := ctx.Value(constants.SessionUserName).(string)
	// ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	switch role {
	case u_s.UserRoleRoot:
		break
	default:
		impl.Logger.Error("you do not have permission to change password")
		return httperror.NewForForbiddenWithSingleField("message", "you do not have permission to change password")
	}

	//
	// Perform our validation and return validation error on any issues detected.
	//

	if err := impl.validateOperationChangeTwoFactorAuthenticationRequest(ctx, req); err != nil {
		impl.Logger.Error("validation error", slog.Any("error", err))
		return err
	}

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

		//
		// Fetch the original user.
		//

		u, err := impl.UserStorer.GetByID(sessCtx, req.UserID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if u == nil {
			return nil, httperror.NewForBadRequestWithSingleField("user_id", "user does not exist")
		}

		// Defensive Code: Tenancy protection
		if u.StoreID != tid && role != u_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("security", "you do not belong to this organization")
		}

		if req.OTPEnabled {
			// CASE 1 OF 2: Added 2FA
			// The following case will force the user to set 2FA on their
			// next first login to load up the wizard and help them through
			// the 2FA setup process.

			// Update user record.
			u.OTPEnabled = true
			u.OTPVerified = false
			u.OTPValidated = false
			u.OTPSecret = ""
			u.OTPAuthURL = ""

		} else {
			// CASE 2 OF 2: Remove 2FA
			// The following code will disable all 2FA details.

			// Update user record.
			u.OTPEnabled = false
			u.OTPVerified = false
			u.OTPValidated = false
			u.OTPSecret = ""
			u.OTPAuthURL = ""
		}

		// Update user record.
		u.ModifiedAt = time.Now()
		u.ModifiedByUserID = userID
		// u.ModifiedByUserName = userName
		// u.ModifiedFromIPAddress = ipAddress

		impl.Logger.Error("record",
			slog.Any("req", req.OTPEnabled),
			slog.Any("u", u.OTPEnabled))

		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("update user error", slog.Any("err", err))
			return nil, err
		}

		return u, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}
