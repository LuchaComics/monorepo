package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
)

type ComicSubmissionTotalCoinsAwardedHTTPHandler struct {
	logger   *slog.Logger
	dbClient *mongo.Client
	service  *service.ComicSubmissionTotalCoinsAwardedService
}

func NewComicSubmissionTotalCoinsAwardedHTTPHandler(
	logger *slog.Logger,
	dbClient *mongo.Client,
	service *service.ComicSubmissionTotalCoinsAwardedService,
) *ComicSubmissionTotalCoinsAwardedHTTPHandler {
	return &ComicSubmissionTotalCoinsAwardedHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *ComicSubmissionTotalCoinsAwardedHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	tenantID, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)

	// Here is where you extract url parameters.
	query := r.URL.Query()
	userIDStr := query.Get("user_id")
	if userIDStr != "" {
		userIDParam, err := primitive.ObjectIDFromHex(userIDStr)
		if err != nil {
			httperror.ResponseError(w, err)
			return
		}
		userID = userIDParam
	}

	filter := &domain.ComicSubmissionFilter{
		TenantID: tenantID,
		UserID:   userID,
	}

	status := query.Get("status")
	if status != "" {
		statusInt, _ := strconv.ParseInt(status, 10, 64)
		filter.Status = int8(statusInt)
	}

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
		resp, err := h.service.Execute(sessCtx)
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

	resp := result.(*service.ComicSubmissionTotalCoinsAwardedResponseIDO)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

}
