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

type ProofOfAuthorityTokenMintHTTPHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.ProofOfAuthorityTokenMintService
}

func NewProofOfAuthorityTokenMintHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	mintTokenService *service.ProofOfAuthorityTokenMintService,
) *ProofOfAuthorityTokenMintHTTPHandler {
	return &ProofOfAuthorityTokenMintHTTPHandler{cfg, logger, mintTokenService}
}

type ProofOfAuthorityTokenMintRequestIDO struct {
	ProofOfAuthorityAccountAddress string `bson:"poa_address" json:"poa_address"`
	ProofOfAuthorityWalletPassword string `bson:"poa_password" json:"poa_password"`
	To                             string `json:"to"`           // Account receiving the Token.
	MetadataURI                    string `json:"metadata_uri"` // URI pointing to Token metadata file.
}

type BlockchainProofOfAuthorityTokenMintResponseIDO struct {
}

func (h *ProofOfAuthorityTokenMintHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalProofOfAuthorityTokenMintRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	pofAddr := common.HexToAddress(req.ProofOfAuthorityAccountAddress)
	toAddr := common.HexToAddress(req.To)

	h.logger.Debug("Received Token mint request",
		slog.Any("metadata_uri", req.MetadataURI),
	)

	serviceExecErr := h.service.Execute(
		ctx,
		&pofAddr,
		req.ProofOfAuthorityWalletPassword,
		&toAddr,
		req.MetadataURI,
	)
	if serviceExecErr != nil {
		httperror.ResponseError(w, serviceExecErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func unmarshalProofOfAuthorityTokenMintRequest(ctx context.Context, r *http.Request) (*ProofOfAuthorityTokenMintRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *ProofOfAuthorityTokenMintRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
