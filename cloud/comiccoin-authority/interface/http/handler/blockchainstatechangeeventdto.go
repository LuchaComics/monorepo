package handler

import (
	"fmt"
	"log/slog"
	"net/http"
	"time"

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

	h.logger.Debug("Blockchain state requested", slog.Any("chain_id", chainIDStr))

	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if chainIDStr == "1" {
				fmt.Fprintf(w, "data: %v", 1)
			} else if chainIDStr == "1" {
				fmt.Fprintf(w, "data: %v", 2)
			} else {
				fmt.Fprintf(w, "data: %v", 3)
			}
			w.(http.Flusher).Flush()
		}
	}

	// // Simulate closing the connection
	// closeNotify := w.(http.CloseNotifier).CloseNotify()
	// <-closeNotify
}
