package controller

import (
	"context"
	"time"

	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl *GatewayControllerImpl) PasswordReset(ctx context.Context, code string, password string) error {
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

		passwordHash, err := impl.Password.GenerateHashFromPassword(password)
		if err != nil {
			impl.Logger.Error("hashing error", slog.Any("error", err))
			return nil, err
		}

		u.PasswordHash = passwordHash
		u.PasswordHashAlgorithm = impl.Password.AlgorithmName()
		u.EmailVerificationCode = "" // Remove email active code so it cannot be used agian.
		u.EmailVerificationExpiry = time.Now()
		u.ModifiedAt = time.Now()

		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("update error", slog.Any("err", err))
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
