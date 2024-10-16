package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type TransferTokenHTTPHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.TransferTokenService
}

func NewTransferTokenHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	mintTokenService *service.TransferTokenService,
) *TransferTokenHTTPHandler {
	return &TransferTokenHTTPHandler{cfg, logger, mintTokenService}
}

type TransferTokenRequestIDO struct {
	TokenOwnerAddress  string `bson:"token_owner_address" json:"token_owner_address"`
	TokenOwnerPassword string `bson:"token_owner_password" json:"token_owner_password"`
	RecipientAddress   string `bson:"recipient_address" json:"recipient_address"`
	TokenID            uint64 `bson:"token_id" json:"token_id"`
}

type BlockchainTransferTokenResponseIDO struct {
}

func (h *TransferTokenHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalTransferTokenRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	tokenOwnerAddr := common.HexToAddress(req.TokenOwnerAddress)
	recipientAddr := common.HexToAddress(req.RecipientAddress)

	h.logger.Debug("Received Token transfer request",
		slog.Any("token_id", req.TokenID),
	)

	serviceExecErr := h.service.Execute(
		ctx,
		&tokenOwnerAddr,
		req.TokenOwnerPassword,
		&recipientAddr,
		req.TokenID,
	)
	if serviceExecErr != nil {
		httperror.ResponseError(w, serviceExecErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func unmarshalTransferTokenRequest(ctx context.Context, r *http.Request) (*TransferTokenRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *TransferTokenRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
