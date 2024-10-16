package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
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
	Address string `json:"address"`
	Balance uint64 `json:"balance"`
	Nonce   uint64 `json:"nonce"`
}

func (h *GetAccountHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, accountAddressStr string) {
	// ctx := r.Context()

	accountAddress := common.HexToAddress(accountAddressStr)

	account, serviceErr := h.getAccountService.Execute(&accountAddress)
	if serviceErr != nil {
		httperror.ResponseError(w, serviceErr)
		return
	}

	// Conver to our HTTP response and send back to the user.
	responsePayload := &AccountGetResponseIDO{
		Address: account.Address.String(),
		Balance: account.Balance,
		Nonce:   account.Nonce,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&responsePayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
