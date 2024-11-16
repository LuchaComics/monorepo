package handler

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type BlockchainStateChangeEventDTOHTTPHandler struct {
	logger  *slog.Logger
	usecase *usecase.BlockchainStateUpdateDetectorUseCase
}

func NewBlockchainStateChangeEventDTOHTTPHandler(
	logger *slog.Logger,
	uc *usecase.BlockchainStateUpdateDetectorUseCase,
) *BlockchainStateChangeEventDTOHTTPHandler {
	return &BlockchainStateChangeEventDTOHTTPHandler{logger, uc}
}

func (h *BlockchainStateChangeEventDTOHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Set CORS headers to allow all origins. You may want to restrict this to specific origins in a production environment.
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Expose-Headers", "Content-Type")

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	chainIDStr := r.URL.Query().Get("chain_id")
	if chainIDStr == "" || (chainIDStr != "1" && chainIDStr != "2") {
		http.Error(w, "Invalid chain_id parameter", http.StatusBadRequest)
		return
	}
	var chainID uint16

	if len(chainIDStr) > 0 {
		var err error
		chainIDInt64, err := strconv.ParseUint(chainIDStr, 10, 16)
		if err != nil {
			log.Println(err)
		}
		chainID = uint16(chainIDInt64)
	}

	h.logger.Debug("Blockchain state change events requested",
		slog.Any("chain_id", chainIDStr))

	defer func() {
		// Simulate closing the connection
		closeNotify := w.(http.CloseNotifier).CloseNotify()
		<-closeNotify
	}()

	for {
		updatedBlockchainState, err := h.usecase.Execute(ctx)
		if err != nil {
			h.logger.Error("Failed detecting blockchain state changes",
				slog.Any("error", err))
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if updatedBlockchainState.ChainID == chainID {
			fmt.Fprintf(w, "data: %v", chainID)
		}
		w.(http.Flusher).Flush()
	}

}
