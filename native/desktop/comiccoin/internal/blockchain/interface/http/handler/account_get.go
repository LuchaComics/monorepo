package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetAccountHTTPHandler struct {
	config            *config.Config
	logger            *slog.Logger
	getAccountService *service.GetAccountService
}

func NewGetAccountHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.GetAccountService,
) *GetAccountHTTPHandler {
	return &GetAccountHTTPHandler{cfg, logger, s}
}

type AccountGetResponseIDO struct {
	ID            string `json:"id"`
	WalletAddress string `json:"wallet_address"`
}

func (h *GetAccountHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, accountID string) {
	// ctx := r.Context()

	account, serviceErr := h.getAccountService.Execute(accountID)
	if serviceErr != nil {
		httperror.ResponseError(w, serviceErr)
		return
	}

	// Conver to our HTTP response and send back to the user.
	responsePayload := &AccountGetResponseIDO{
		ID:            account.ID,
		WalletAddress: account.WalletAddress.String(),
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&responsePayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
