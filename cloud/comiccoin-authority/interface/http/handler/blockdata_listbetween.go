package handler

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type ListBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler struct {
	logger  *slog.Logger
	service *service.ListBlockDataFilteredBetweenBlockNumbersInChainIDService
}

func NewListBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler(
	logger *slog.Logger,
	s1 *service.ListBlockDataFilteredBetweenBlockNumbersInChainIDService,
) *ListBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler {
	return &ListBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler{logger, s1}
}

func (h *ListBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	h.logger.Debug("List block data in hashes filter requested")

	// Here is where you extract url parameters.
	query := r.URL.Query()

	chainIDstr := query.Get("chain_id")
	if chainIDstr == "" {
		http.Error(w, "Missing `chain_id` parameter.", http.StatusInternalServerError)
		return
	}
	chainID, err := strconv.ParseInt(chainIDstr, 0, 16)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	startStr := query.Get("start")
	if startStr == "" {
		http.Error(w, "Missing `start` parameter.", http.StatusInternalServerError)
		return
	}
	start, err := strconv.ParseInt(startStr, 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	endStr := query.Get("end")
	if endStr == "" {
		http.Error(w, "Missing `end` parameter.", http.StatusInternalServerError)
		return
	}
	end, err := strconv.ParseInt(endStr, 0, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	blockchainState, err := h.service.Execute(ctx, uint64(start), uint64(end), uint16(chainID))
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
