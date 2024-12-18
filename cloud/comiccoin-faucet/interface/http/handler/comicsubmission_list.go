package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"
	"time"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
)

type ComicSubmissionListByFilterHTTPHandler struct {
	logger   *slog.Logger
	dbClient *mongo.Client
	service  *service.ComicSubmissionListByFilterService
}

func NewComicSubmissionListByFilterHTTPHandler(
	logger *slog.Logger,
	dbClient *mongo.Client,
	service *service.ComicSubmissionListByFilterService,
) *ComicSubmissionListByFilterHTTPHandler {
	return &ComicSubmissionListByFilterHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *ComicSubmissionListByFilterHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	tenantID, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)

	// Get query parameters
	query := r.URL.Query()

	// Initialize the filter with required fields
	filter := &service.ComicSubmissionFilterRequestID{
		TenantID: tenantID,
	}

	// Handle name filter
	if name := query.Get("name"); name != "" {
		filter.Name = &name
	}

	// Handle status filter
	if statusStr := query.Get("status"); statusStr != "" {
		if status, err := strconv.ParseInt(statusStr, 10, 8); err == nil {
			statusInt8 := int8(status)
			filter.Status = statusInt8
		} else {
			h.logger.Error("status parse error",
				slog.Any("error", err),
				slog.String("status", statusStr))
		}
	}

	// Handle created_at range filters
	if startDateStr := query.Get("created_at_start"); startDateStr != "" {
		if startDate, err := time.Parse(time.RFC3339, startDateStr); err == nil {
			filter.CreatedAtStart = &startDate
		} else {
			h.logger.Error("created_at_start parse error",
				slog.Any("error", err),
				slog.String("created_at_start", startDateStr))
		}
	}

	if endDateStr := query.Get("created_at_end"); endDateStr != "" {
		if endDate, err := time.Parse(time.RFC3339, endDateStr); err == nil {
			filter.CreatedAtEnd = &endDate
		} else {
			h.logger.Error("created_at_end parse error",
				slog.Any("error", err),
				slog.String("created_at_end", endDateStr))
		}
	}

	// Handle cursor-based pagination
	if lastIDStr := query.Get("last_id"); lastIDStr != "" {
		if lastID, err := primitive.ObjectIDFromHex(lastIDStr); err == nil {
			filter.LastID = &lastID
		} else {
			h.logger.Error("last_id parse error",
				slog.Any("error", err),
				slog.String("last_id", lastIDStr))
		}
	}

	if lastCreatedAtStr := query.Get("last_created_at"); lastCreatedAtStr != "" {
		if lastCreatedAt, err := time.Parse(time.RFC3339, lastCreatedAtStr); err == nil {
			filter.LastCreatedAt = &lastCreatedAt
		} else {
			h.logger.Error("last_created_at parse error",
				slog.Any("error", err),
				slog.String("last_created_at", lastCreatedAtStr))
		}
	}

	// Handle limit
	if limitStr := query.Get("limit"); limitStr != "" {
		if limit, err := strconv.ParseInt(limitStr, 10, 64); err == nil {
			filter.Limit = limit
		} else {
			h.logger.Error("limit parse error",
				slog.Any("error", err),
				slog.String("limit", limitStr))
		}
	}

	// Filter by user account.
	loggedInUserRole, _ := ctx.Value(constants.SessionUserRole).(int8)
	if userIDStr := query.Get("user_id"); userIDStr != "" {
		if userID, err := primitive.ObjectIDFromHex(userIDStr); err == nil {
			filter.UserID = userID
		} else {
			h.logger.Error("user_id parse error",
				slog.Any("error", err),
				slog.String("user_id", userID.Hex()))
			httperror.ResponseError(w, err)
			return
		}
	}

	// Permission handling: Non-administrators cannot filter other userIDs.
	if loggedInUserRole != domain.UserRoleRoot {
		loggedInUserID, _ := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
		filter.UserID = loggedInUserID
		h.logger.Debug("Forcing filtering",
			slog.Any("loggedInUserID", loggedInUserID),
			slog.Any("loggedInUserRole", loggedInUserRole))
	}

	// Start the transaction
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
		resp, err := h.service.Execute(sessCtx, filter)
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

	resp := result.(*service.ComicSubmissionFilterResultResponseIDO)

	// Set response headers
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		h.logger.Error("encoding response error",
			slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
