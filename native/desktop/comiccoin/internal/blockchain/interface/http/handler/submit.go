package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type SubmitHTTPHandler struct {
	config *config.Config
	logger *slog.Logger
	// createAccountService *service.SubmitService
}

func NewSubmitHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	// s *service.SubmitService,
) *SubmitHTTPHandler {
	return &SubmitHTTPHandler{cfg, logger}
}

type BlockchainSubmitRequestIDO struct {
	// Name of the account
	FromAccountName string `json:"from_account_name"`

	AccountWalletPassword string `json:"account_wallet_password"`

	// Recipient’s public key
	To string `json:"to"`

	// Value is amount of coins being transferred
	Value uint64 `json:"value"`

	// Data is any NFT related data attached
	Data []byte `json:"data"`
}

type BlockchainSubmitResponseIDO struct {
}

func (h *SubmitHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestPayload, err := unmarshalSubmitRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	_ = requestPayload

	// account, serviceErr := h.createAccountService.Execute(h.config.App.DirPath, requestPayload.ID, requestPayload.WalletPassword)
	// if serviceErr != nil {
	// 	httperror.ResponseError(w, serviceErr)
	// 	return
	// }
	//
	// // Conver to our HTTP response and send back to the user.
	// responsePayload := &AccountCreateResponseIDO{
	// 	ID:            account.ID,
	// 	WalletAddress: account.WalletAddress.String(),
	// }
	//
	// w.WriteHeader(http.StatusCreated)
	// if err := json.NewEncoder(w).Encode(&responsePayload); err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
}

func unmarshalSubmitRequest(ctx context.Context, r *http.Request) (*BlockchainSubmitRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *BlockchainSubmitRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
