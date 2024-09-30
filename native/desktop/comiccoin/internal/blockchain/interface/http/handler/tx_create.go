package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateTransactionHTTPHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.CreateTransactionService
}

func NewCreateTransactionHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	createTransactionService *service.CreateTransactionService,
) *CreateTransactionHTTPHandler {
	return &CreateTransactionHTTPHandler{cfg, logger, createTransactionService}
}

type BlockchainCreateTransactionRequestIDO struct {
	// Name of the account
	FromAccountID string `json:"from_account_id"`

	AccountWalletPassword string `json:"account_wallet_password"`

	// Recipient’s public key
	To string `json:"to"`

	// Value is amount of coins being transferred
	Value uint64 `json:"value"`

	// Data is any NFT related data attached
	Data []byte `json:"data"`
}

type BlockchainCreateTransactionResponseIDO struct {
}

func (h *CreateTransactionHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalCreateTransactionRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	toAddr := common.HexToAddress(req.To)

	serviceExecErr := h.service.Execute(
		ctx,
		req.FromAccountID,
		req.AccountWalletPassword,
		&toAddr,
		req.Value,
		req.Data,
	)
	if serviceExecErr != nil {
		httperror.ResponseError(w, serviceExecErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func unmarshalCreateTransactionRequest(ctx context.Context, r *http.Request) (*BlockchainCreateTransactionRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *BlockchainCreateTransactionRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
