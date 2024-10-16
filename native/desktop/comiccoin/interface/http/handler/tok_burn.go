package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type BurnTokenHTTPHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.BurnTokenService
}

func NewBurnTokenHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	mintTokenService *service.BurnTokenService,
) *BurnTokenHTTPHandler {
	return &BurnTokenHTTPHandler{cfg, logger, mintTokenService}
}

type BurnTokenRequestIDO struct {
	TokenOwnerAddress  string `bson:"token_owner_address" json:"token_owner_address"`
	TokenOwnerPassword string `bson:"token_owner_password" json:"token_owner_password"`
	TokenID            uint64 `bson:"token_id" json:"token_id"`
}

type BlockchainBurnTokenResponseIDO struct {
}

func (h *BurnTokenHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalBurnTokenRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	tokenOwnerAddr := common.HexToAddress(req.TokenOwnerAddress)

	h.logger.Debug("Received Token burn request",
		slog.Any("token_id", req.TokenID),
	)

	serviceExecErr := h.service.Execute(
		ctx,
		&tokenOwnerAddr,
		req.TokenOwnerPassword,
		req.TokenID,
	)
	if serviceExecErr != nil {
		httperror.ResponseError(w, serviceExecErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func unmarshalBurnTokenRequest(ctx context.Context, r *http.Request) (*BurnTokenRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *BurnTokenRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
