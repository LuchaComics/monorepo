package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type ProfileChangePasswordRequestIDO struct {
	OldPassword         string `json:"old_password"`
	NewPassword         string `json:"new_password"`
	NewPasswordRepeated string `json:"new_password_repeated"`
}

func ValidateProfileChangePassworRequest(dirtyData *ProfileChangePasswordRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.OldPassword == "" {
		e["old_password"] = "missing value"
	}
	if dirtyData.NewPassword == "" {
		e["new_password"] = "missing value"
	}
	if dirtyData.NewPasswordRepeated == "" {
		e["new_password_repeated"] = "missing value"
	}
	if dirtyData.NewPasswordRepeated != dirtyData.NewPassword {
		e["new_password"] = "does not match"
		e["new_password_repeated"] = "does not match"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *GatewayControllerImpl) ProfileChangePassword(ctx context.Context, req *ProfileChangePasswordRequestIDO) error {
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
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		if err := ValidateProfileChangePassworRequest(req); err != nil {
			impl.Logger.Warn("user validation failed", slog.Any("err", err))
			return nil, err
		}

		// Verify the inputted password and hashed password match.
		if passwordMatch, _ := impl.Password.ComparePasswordAndHash(req.OldPassword, u.PasswordHash); passwordMatch == false {
			impl.Logger.Warn("password check validation error")
			return nil, httperror.NewForBadRequestWithSingleField("old_password", "old password do not match with record of existing password")
		}

		passwordHash, err := impl.Password.GenerateHashFromPassword(req.NewPassword)
		if err != nil {
			impl.Logger.Error("hashing error", slog.Any("error", err))
			return nil, err
		}
		u.PasswordHash = passwordHash
		u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("user update by id error", slog.Any("error", err))
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
