package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type CreateAccountHTTPHandler struct {
	config               *config.Config
	logger               *slog.Logger
	createAccountService *service.CreateAccountService
}

func NewCreateAccountHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.CreateAccountService,
) *CreateAccountHTTPHandler {
	return &CreateAccountHTTPHandler{cfg, logger, s}
}

type AccountCreateRequestIDO struct {
	ID             string `json:"id"`
	WalletPassword string `json:"wallet_password"`
}

type AccountCreateResponseIDO struct {
	ID            string `json:"id"`
	WalletAddress string `json:"wallet_address"`
}

func (h *CreateAccountHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	requestPayload, err := unmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	account, serviceErr := h.createAccountService.Execute(h.config.App.DirPath, requestPayload.ID, requestPayload.WalletPassword)
	if serviceErr != nil {
		httperror.ResponseError(w, serviceErr)
		return
	}

	// Conver to our HTTP response and send back to the user.
	responsePayload := &AccountCreateResponseIDO{
		ID:            account.ID,
		WalletAddress: account.WalletAddress.String(),
	}

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&responsePayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func unmarshalCreateRequest(ctx context.Context, r *http.Request) (*AccountCreateRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *AccountCreateRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}