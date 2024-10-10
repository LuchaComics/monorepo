package handler

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
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
	To                             string `json:"to"` // Account receiving the NFT.

	// This is the URL to the image of the item. Can be just about any type of image (including SVGs, which will be cached into PNGs by OpenSea), IPFS or Arweave URLs or paths. We recommend using a minimum 3000 x 3000 image.
	Image string `bson:"image" json:"image"`

	// This is the URL that will appear below the asset's image on OpenSea and will allow users to leave OpenSea and view the item on your site.
	ExternalURL string `bson:"external_url" json:"external_url"`

	// A human-readable description of the item. Markdown is supported.
	Description string `bson:"description" json:"description"`

	// Name of the item.
	Name string `bson:"name" json:"name"`

	// These are the attributes for the item, which will show up on the OpenSea page for the item. (see below
	Attributes []*NFTMetadataAttributeRequestIDO `bson:"attributes" json:"attributes"`

	// Background color of the item on OpenSea. Must be a six-character hexadecimal without a pre-pended #.
	BackgroundColor string `bson:"background_color" json:"background_color"`

	// A URL to a multi-media attachment for the item. The file extensions GLTF, GLB, WEBM, MP4, M4V, OGV, and OGG are supported, along with the audio-only extensions MP3, WAV, and OGA.
	//
	// Animation_url also supports HTML pages, allowing you to build rich experiences and interactive NFTs using JavaScript canvas, WebGL, and more. Scripts and relative paths within the HTML page are now supported. However, access to browser extensions is not supported.
	AnimationURL string `bson:"animation_url" json:"animation_url"`

	// A URL to a YouTube video (only used if animation_url is not provided).
	YoutubeURL string `bson:"youtube_url" json:"youtube_url"`

	MetadataURI string `json:"metadata_uri"` // URI pointing to NFT metadata file.
}

type NFTMetadataAttributeRequestIDO struct {
	DisplayType string `bson:"display_type" json:"display_type"`
	TraitType   string `bson:"trait_type" json:"trait_type"`
	Value       string `bson:"value" json:"value"`
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

	attributes := make([]*domain.NFTMetadataAttribute, len(req.Attributes))
	for _, attr := range req.Attributes {
		newAttr := &domain.NFTMetadataAttribute{
			DisplayType: attr.DisplayType,
			TraitType:   attr.TraitType,
			Value:       attr.Value,
		}
		attributes = append(attributes, newAttr)
	}

	nftMetadata := &domain.NFTMetadata{
		Image:           req.Image,
		ExternalURL:     req.ExternalURL,
		Description:     req.Description,
		Attributes:      attributes,
		Name:            req.Name,
		BackgroundColor: req.BackgroundColor,
		AnimationURL:    req.AnimationURL,
		YoutubeURL:      req.YoutubeURL,
	}

	h.logger.Debug("Received NFT mint request",
		slog.Any("image", req.Image),
		slog.Any("external_url", req.ExternalURL),
		slog.Any("description", req.Description),
		slog.Any("attributes", attributes),
		slog.Any("name", req.Name),
		slog.Any("background_color", req.BackgroundColor),
		slog.Any("animation_url", req.AnimationURL),
		slog.Any("youtube_url", req.YoutubeURL),
		slog.Any("metadata_uri", req.MetadataURI),
	)

	serviceExecErr := h.service.Execute(
		ctx,
		&pofAddr,
		req.ProofOfAuthorityWalletPassword,
		&toAddr,
		nftMetadata,
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
