package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/ethereum/go-ethereum/common"
)

type MintNFTHTTPHandler struct {
	config  *config.Config
	logger  *slog.Logger
	service *service.MintNFTService
}

func NewMintNFTHTTPHandler(
	cfg *config.Config,
	logger *slog.Logger,
	mintNFTService *service.MintNFTService,
) *MintNFTHTTPHandler {
	return &MintNFTHTTPHandler{cfg, logger, mintNFTService}
}

type MintNFTRequestIDO struct {
	ProofOfAuthorityAccountAddress string `bson:"poa_address" json:"poa_address"`
	ProofOfAuthorityWalletPassword string `bson:"poa_password" json:"poa_password"`
	To                             string `json:"to"`           // Account receiving the NFT.
	MetadataURI                    string `json:"metadata_uri"` // URI pointing to NFT metadata file.
}

type BlockchainMintNFTResponseIDO struct {
}

func (h *MintNFTHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	req, err := unmarshalMintNFTRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	_ = req

	pofAddr := common.HexToAddress(req.ProofOfAuthorityAccountAddress)
	toAddr := common.HexToAddress(req.To)

	h.logger.Debug("Received NFT mint request",
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

func unmarshalMintNFTRequest(ctx context.Context, r *http.Request) (*MintNFTRequestIDO, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData *MintNFTRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return requestData, nil
}
