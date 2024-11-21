package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/ethereum/go-ethereum/common"
)

type TokenListByOwnerHTTPHandler struct {
	logger  *slog.Logger
	service *service.TokenListByOwnerService
}

func NewTokenListByOwnerHTTPHandler(
	logger *slog.Logger,
	s1 *service.TokenListByOwnerService,
) *TokenListByOwnerHTTPHandler {
	return &TokenListByOwnerHTTPHandler{logger, s1}
}

func (h *TokenListByOwnerHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Here is where you extract url parameters.
	query := r.URL.Query()
	ownerAddressStr := query.Get("owner_address")
	if ownerAddressStr == "" {
		err := httperror.NewForNotFoundWithSingleField("owner_address", "missing value")
		httperror.ResponseError(w, err)
		return
	}

	h.logger.Debug("Token list by owner requested",
		slog.Any("owner_address", ownerAddressStr))

	address := common.HexToAddress(strings.ToLower(ownerAddressStr))
	resp, err := h.service.Execute(ctx, &address)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
