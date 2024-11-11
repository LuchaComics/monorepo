package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type ListBlockDataFilteredInHashesHTTPHandler struct {
	logger  *slog.Logger
	service *service.ListBlockDataFilteredInHashesService
}

func NewListBlockDataFilteredInHashesHTTPHandler(
	logger *slog.Logger,
	s1 *service.ListBlockDataFilteredInHashesService,
) *ListBlockDataFilteredInHashesHTTPHandler {
	return &ListBlockDataFilteredInHashesHTTPHandler{logger, s1}
}

func (h *ListBlockDataFilteredInHashesHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug("List block data in hashes filter requested")

	// Here is where you extract url parameters.
	query := r.URL.Query()

	// chainIDstr := query.Get("chain_id")
	// if chainIDstr == "" {
	// 	chainIDstr = "1"
	// }
	// chainID, err := strconv.ParseInt(chainIDstr, 0, 16)
	// if err != nil {
	// 	http.Error(w, err.Error(), http.StatusInternalServerError)
	// 	return
	// }
	hashesStr := query.Get("hashes")
	if hashesStr == "" {
		http.Error(w, "Missing `hashes` parameter", http.StatusBadRequest)
		return
	}
	hashes := strings.Split(hashesStr, ",")
	blockchainState, err := h.service.Execute(ctx, hashes)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&blockchainState); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
