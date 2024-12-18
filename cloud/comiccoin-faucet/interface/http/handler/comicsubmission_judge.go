package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
)

type ComicSubmissionJudgeHTTPHandler struct {
	logger   *slog.Logger
	dbClient *mongo.Client
	service  *service.ComicSubmissionJudgeService
}

func NewComicSubmissionJudgeHTTPHandler(
	logger *slog.Logger,
	dbClient *mongo.Client,
	service *service.ComicSubmissionJudgeService,
) *ComicSubmissionJudgeHTTPHandler {
	return &ComicSubmissionJudgeHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *ComicSubmissionJudgeHTTPHandler) unmarshalLoginRequest(
	ctx context.Context,
	r *http.Request,
) (*service.ComicSubmissionJudgeVerdictRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData service.ComicSubmissionJudgeVerdictRequestIDO

	defer r.Body.Close()

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(r.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(teeReader).Decode(&requestData) // [1]
	if err != nil {
		h.logger.Error("decoding error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *ComicSubmissionJudgeHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.unmarshalLoginRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// For security, we will attach via user logged in context the authenticated
	// user in our system.
	userID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)
	data.AdminUserID = userID
	data.AdminUserIPAddress = ipAddress

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
		resp, err := h.service.Execute(sessCtx, data)
		if err != nil {
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

	resp := result.(*domain.ComicSubmission)

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
