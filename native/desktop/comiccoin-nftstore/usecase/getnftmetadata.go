package usecase

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/domain"
)

// GetNFTMetadataUseCase structure represents the use-case for downloading
// the metadata file via a URI and return the bytes content.
type GetNFTMetadataUseCase struct {
	logger *slog.Logger
}

func NewGetNFTMetadataUseCase(logger *slog.Logger) *GetNFTMetadataUseCase {
	return &GetNFTMetadataUseCase{logger}
}

func (uc *GetNFTMetadataUseCase) Execute(metadataURI string) (*domain.NFTMetadata, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if metadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed downloading metadata uri.",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Check protocol and download using the particular protocol.
	//

	if strings.Contains(metadataURI, "https://") {
		return uc.downloadViaHTTPS(metadataURI)
	} else if strings.Contains(metadataURI, "ipfs://") {
		return uc.downloadViaIPFS(metadataURI)
	}

	return nil, fmt.Errorf("Unsupported protocol in URI, only accepted options are: %v", "https and ipfs")
}

func (uc *GetNFTMetadataUseCase) downloadViaHTTPS(metadataURI string) (*domain.NFTMetadata, error) {
	uc.logger.Debug("fetching from uri via HTTPS protocol.",
		slog.String("metadataURI", metadataURI))

	resp, err := http.Get(metadataURI)
	if err != nil {
		uc.logger.Error("Failed fetching metadata uri via http.",
			slog.Any("error", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code error: %d", resp.StatusCode)
		uc.logger.Error("Status code error",
			slog.Any("error", err))
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		uc.logger.Error("Failed read all.",
			slog.Any("error", err))
		return nil, err
	}

	var metadata domain.NFTMetadata
	err = json.Unmarshal(body, &metadata)
	if err != nil {
		uc.logger.Error("Failed unmarshalling body",
			slog.Any("error", err))
	}

	return &metadata, nil
}

func (uc *GetNFTMetadataUseCase) downloadViaIPFS(metadataURI string) (*domain.NFTMetadata, error) {
	cid := strings.Replace(metadataURI, "ipfs://", "", -1)
	url := strings.Replace(constants.IPFSPublicHTTPSGatewayBaseURL, "{CID}", cid, -1)
	uc.logger.Debug("fetching from uri via IPFS protocol.",
		slog.String("ipfs_cid", cid),
		slog.String("url", url))

	return uc.downloadViaHTTPS(url)
}
