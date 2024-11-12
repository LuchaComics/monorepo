package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type SignedTransactionSubmissionHTTPHandler struct {
	logger  *slog.Logger
	service *service.SignedTransactionSubmissionService
}

func NewSignedTransactionSubmissionHTTPHandler(
	logger *slog.Logger,
	transferCoinsService *service.SignedTransactionSubmissionService,
) *SignedTransactionSubmissionHTTPHandler {
	return &SignedTransactionSubmissionHTTPHandler{logger, transferCoinsService}
}

type SignedTransactionSubmissionRequestIDO struct {
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

type BlockchainSignedTransactionSubmissionResponseIDO struct {
}

func (h *SignedTransactionSubmissionHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalSignedTransactionSubmissionRequest(ctx, r)
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
	)
	if serviceExecErr != nil {
		httperror.ResponseError(w, serviceExecErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func unmarshalSignedTransactionSubmissionRequest(ctx context.Context, r *http.Request) (*SignedTransactionSubmissionRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *SignedTransactionSubmissionRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
