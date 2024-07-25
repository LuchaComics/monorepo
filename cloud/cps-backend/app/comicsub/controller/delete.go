package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl *ComicSubmissionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
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
		// STEP 1: Lookup the record or error.
		submission, err := impl.GetByID(sessCtx, id)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if submission == nil {
			impl.Logger.Error("database returns nothing from get by id")
			return nil, err
		}

		// STEP 2: Delete from remote storage.
		impl.Logger.Debug("Will delete previous uploaded findings form",
			slog.String("path", submission.FindingsFormObjectKey))

		// Delete previous record from remote storage.
		if err := impl.S3.DeleteByKeys(sessCtx, []string{submission.FindingsFormObjectKey}); err != nil {
			impl.Logger.Warn("object delete by keys error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}

		impl.Logger.Debug("Will delete previous uploaded label",
			slog.String("path", submission.LabelObjectKey))

		// Delete previous record from remote storage.
		if err := impl.S3.DeleteByKeys(sessCtx, []string{submission.LabelObjectKey}); err != nil {
			impl.Logger.Warn("object delete by keys error", slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}

		// STEP 3: Delete from database.
		if err := impl.ComicSubmissionStorer.DeleteByID(sessCtx, id); err != nil {
			impl.Logger.Error("database delete by id error", slog.Any("error", err))
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
