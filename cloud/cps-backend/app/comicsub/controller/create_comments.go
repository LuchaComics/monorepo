package controller

import (
	"context"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
)

func (impl *ComicSubmissionControllerImpl) CreateComment(ctx context.Context, submissionID primitive.ObjectID, content string) (*submission_s.ComicSubmission, error) {
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
		// Fetch the original submission.
		s, err := impl.ComicSubmissionStorer.GetByID(ctx, submissionID)
		if err != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", err))
			return nil, err
		}
		if s == nil {
			return nil, nil
		}

		// Create our comment.
		comment := &submission_s.ComicSubmissionComment{
			ID:               primitive.NewObjectID(),
			Content:          content,
			StoreID:          ctx.Value(constants.SessionUserStoreID).(primitive.ObjectID),
			CreatedByUserID:  ctx.Value(constants.SessionUserID).(primitive.ObjectID),
			CreatedByName:    ctx.Value(constants.SessionUserName).(string),
			CreatedAt:        time.Now(),
			ModifiedByUserID: ctx.Value(constants.SessionUserID).(primitive.ObjectID),
			ModifiedByName:   ctx.Value(constants.SessionUserName).(string),
			ModifiedAt:       time.Now(),
		}

		// Add our comment to the comments.
		s.ModifiedByUserID = ctx.Value(constants.SessionUserID).(primitive.ObjectID)
		s.ModifiedAt = time.Now()
		s.Comments = append(s.Comments, comment)

		// Save to the database the modified submission.
		if err := impl.ComicSubmissionStorer.UpdateByID(ctx, s); err != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", err))
			return nil, err
		}

		return s, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*submission_s.ComicSubmission), nil
}
