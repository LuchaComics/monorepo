package handler

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type BlockchainStateChangeEventDTOHTTPHandler struct {
	logger  *slog.Logger
	service *service.BlockchainStateChangeSubscriptionService
}

func NewBlockchainStateChangeEventDTOHTTPHandler(
	logger *slog.Logger,
	s *service.BlockchainStateChangeSubscriptionService,
) *BlockchainStateChangeEventDTOHTTPHandler {
	return &BlockchainStateChangeEventDTOHTTPHandler{logger, s}
}

func (h *BlockchainStateChangeEventDTOHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

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
		slog.Any("chain_id", chainIDStr),
		slog.Any("ip_address", ipAddress))

	defer func() {
		// Simulate closing the connection
		closeNotify := w.(http.CloseNotifier).CloseNotify()
		<-closeNotify
	}()

	for {
		select {
		case <-ctx.Done():
			h.logger.Debug("Context canceled. Stopping event stream.",
				slog.Any("ip_address", ipAddress))
			return
		default:
			updatedBlockchainState, err := h.service.Execute(ctx)
			if err != nil {
				h.logger.Error("Failed detecting blockchain state changes",
					slog.Any("error", err),
					slog.Any("ip_address", ipAddress))
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if updatedBlockchainState.ChainID == chainID {
				fmt.Fprintf(w, "data: %v", chainID)
			}
			w.(http.Flusher).Flush()
		}
	}
}
