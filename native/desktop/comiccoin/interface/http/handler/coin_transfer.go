package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
)

type TransferCoinHTTPHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.TransferCoinService
}

func NewTransferCoinHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	transferCoinService *service.TransferCoinService,
) *TransferCoinHTTPHandler {
	return &TransferCoinHTTPHandler{cfg, logger, transferCoinService}
}

type TransferCoinRequestIDO struct {
	// Name of the account
	SenderAccountAddress string `json:"sender_account_address"`

	SenderAccountPassword string `json:"sender_account_password"`

	// Value is amount of coins being transferred
	Value uint64 `json:"value"`

	// Recipient’s public key
	RecipientAddress string `json:"recipient_address"`

	// Data is any Token related data attached
	Data string `json:"data"`
}

type BlockchainTransferCoinResponseIDO struct {
}

func (h *TransferCoinHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalTransferCoinRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	toAddr := common.HexToAddress(strings.ToLower(req.RecipientAddress))
	senderAddr := common.HexToAddress(strings.ToLower(req.SenderAccountAddress))

	h.logger.Debug("tx submit received",
		slog.Any("sender", senderAddr),
		slog.Any("receipient", toAddr),
		slog.Any("value", req.Value),
		slog.Any("data", req.Data))

	serviceExecErr := h.service.Execute(
		ctx,
		&senderAddr,
		req.SenderAccountPassword,
		&toAddr,
		req.Value,
		[]byte(req.Data),
	)
	if serviceExecErr != nil {
		httperror.ResponseError(w, serviceExecErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func unmarshalTransferCoinRequest(ctx context.Context, r *http.Request) (*TransferCoinRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *TransferCoinRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
