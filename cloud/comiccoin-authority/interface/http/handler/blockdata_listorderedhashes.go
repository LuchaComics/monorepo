package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type ListAllBlockDataOrderedHashesHTTPHandler struct {
	logger  *slog.Logger
	service *service.BlockDataListAllOrderedHashesService
}

func NewListAllBlockDataOrderedHashesHTTPHandler(
	logger *slog.Logger,
	s1 *service.BlockDataListAllOrderedHashesService,
) *ListAllBlockDataOrderedHashesHTTPHandler {
	return &ListAllBlockDataOrderedHashesHTTPHandler{logger, s1}
}

func (h *ListAllBlockDataOrderedHashesHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug("Blockdata ordered hashes requested")

	// Here is where you extract url parameters.
	query := r.URL.Query()

	chainIDstr := query.Get("chain_id")
	if chainIDstr == "" {
		chainIDstr = "1"
	}
	chainID, err := strconv.ParseInt(chainIDstr, 0, 16)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp, err := h.service.Execute(ctx, uint16(chainID))
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
