package handler

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type GetTokenHTTPHandler struct {
	config          *config.Config
	logger          *slog.Logger
	getTokenService *service.GetTokenService
}

func NewGetTokenHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	s *service.GetTokenService,
) *GetTokenHTTPHandler {
	return &GetTokenHTTPHandler{cfg, logger, s}
}

type TokenGetResponseIDO struct {
	ID          uint64 `json:"id"`
	Owner       string `json:"owner"`
	MetadataURI string `json:"metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
	Nonce       uint64 `json:"nonce"`        // ComicCoin: Newly minted tokens always start at zero and for every transaction action afterwords (transfer, burn, etc) this value is increment by 1.
}

func (h *GetTokenHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, tokenIDStr string) {
	tokenID, err := strconv.ParseUint(tokenIDStr, 10, 64)
	if err != nil {
		h.logger.Error("failed parsing argument",
			slog.String("token_id_str", tokenIDStr))
		badRequestErr := httperror.NewForBadRequestWithSingleField("token_id", fmt.Sprintf("error: %v", err))
		httperror.ResponseError(w, badRequestErr)
		return
	}

	token, serviceErr := h.getTokenService.Execute(tokenID)
	if serviceErr != nil {
		httperror.ResponseError(w, serviceErr)
		return
	}

	// Conver to our HTTP response and send back to the user.
	responsePayload := &TokenGetResponseIDO{
		ID:          token.ID,
		Owner:       token.Owner.String(),
		MetadataURI: token.MetadataURI,
		Nonce:       token.Nonce,
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(&responsePayload); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
