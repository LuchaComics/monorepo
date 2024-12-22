package handler

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
)

type BlockchainStateServerSentEventsHTTPHandler struct {
	logger  *slog.Logger
	service *service.GetBlockchainStateService
}

func NewBlockchainStateServerSentEventsHTTPHandler(
	logger *slog.Logger,
	s *service.GetBlockchainStateService,
) *BlockchainStateServerSentEventsHTTPHandler {
	return &BlockchainStateServerSentEventsHTTPHandler{logger, s}
}

func (h *BlockchainStateServerSentEventsHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
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

	h.logger.Debug("Blockchain state server sent events connected client",
		slog.Any("chain_id", chainIDStr),
		slog.Any("ip_address", ipAddress))

	// Create a channel for client disconnection
	clientGone := r.Context().Done()

	// Create ticker that will send a SSE to the client every ten seconds.
	t := time.NewTicker(10 * time.Second)
	defer t.Stop()
	for {
		select {
		case <-clientGone:
			h.logger.Debug("Client disconnected",
				slog.Any("chain_id", chainIDStr),
				slog.Any("ip_address", ipAddress))
			return
		case <-t.C:

			blockchainState, err := h.service.Execute(ctx, uint16(chainID))
			if err != nil {
				httperror.ResponseError(w, err)
				return
			}

			if blockchainState.ChainID == chainID {
				h.logger.Debug("Sending sse to client...",
					slog.Any("chain_id", chainIDStr),
					slog.Any("latest_hash", blockchainState.LatestHash),
					slog.Any("ip_address", ipAddress))

				// Send an event to the client.
				fmt.Fprintf(w, "data: %v\n\n", blockchainState.LatestHash)
			}

			w.(http.Flusher).Flush()
			time.Sleep(1 * time.Second)
		}
	}
}
