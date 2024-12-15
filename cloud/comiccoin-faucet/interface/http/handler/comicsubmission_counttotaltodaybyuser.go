package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
)

type ComicSubmissionCountTotalCreatedTodayByUserHTTPHandler struct {
	logger   *slog.Logger
	dbClient *mongo.Client
	service  *service.ComicSubmissionCountTotalCreatedTodayByUserService
}

func NewComicSubmissionCountTotalCreatedTodayByUserHTTPHandler(
	logger *slog.Logger,
	dbClient *mongo.Client,
	service *service.ComicSubmissionCountTotalCreatedTodayByUserService,
) *ComicSubmissionCountTotalCreatedTodayByUserHTTPHandler {
	return &ComicSubmissionCountTotalCreatedTodayByUserHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *ComicSubmissionCountTotalCreatedTodayByUserHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	////
	//// Start the transaction.
	////

	session, err := h.dbClient.StartSession()
	if err != nil {
		h.logger.Error("start session error",
			slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		resp, err := h.service.Execute(sessCtx, userID)
		if err != nil {
			// httperror.ResponseError(w, err)
			return nil, err
		}
		return resp, nil
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		h.logger.Error("session failed error",
			slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}

	resp := result.(*service.ComicSubmissionCountTotalCreatedTodayByUserServiceResponseIDO)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
