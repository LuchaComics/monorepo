package controller

import (
	"context"
	"strings"

	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl *GatewayControllerImpl) ForgotPassword(ctx context.Context, email string) error {
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
		// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
		email = strings.ToLower(email)

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByEmail(sessCtx, email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error", slog.String("email", email))
			return nil, httperror.NewForBadRequestWithSingleField("email", "does not exist")
		}

		// Generate unique token and save it to the user record.
		u.EmailVerificationCode = impl.UUID.NewUUID()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Warn("user update by id failed", slog.Any("error", err))
			return nil, err
		}

		// Send password reset email.
		return impl.TemplatedEmailer.SendForgotPasswordEmail(email, u.EmailVerificationCode, u.FirstName), nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}
