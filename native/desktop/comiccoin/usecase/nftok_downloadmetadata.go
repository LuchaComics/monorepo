package usecase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type DownloadMetadataNonFungibleTokenUseCase struct {
	logger *slog.Logger
	repo   domain.NFTAssetRepository
}

func NewDownloadMetadataNonFungibleTokenUseCase(logger *slog.Logger, r domain.NFTAssetRepository) *DownloadMetadataNonFungibleTokenUseCase {
	return &DownloadMetadataNonFungibleTokenUseCase{logger, r}
}

func (uc *DownloadMetadataNonFungibleTokenUseCase) Execute(tokenID *big.Int, metadataURI string, dirPath string) (*domain.NonFungibleTokenMetadata, string, error) {
	// Confirm URI is using protocol our app supports.
	if strings.Contains(metadataURI, "ipfs://") {
		uc.logger.Debug("Downloading metadata via ipfs...",
			slog.Any("token_id", tokenID),
			slog.Any("metadata_uri", metadataURI))
		return uc.executeForIPFS(tokenID, metadataURI, dirPath)
	} else if strings.Contains(metadataURI, "https://") {
		uc.logger.Debug("Downloading metadata via https...",
			slog.Any("token_id", tokenID),
			slog.Any("metadata_uri", metadataURI))
		return uc.executeForHTTPS(tokenID, metadataURI, dirPath)
	}

	uc.logger.Error("Token metadata URI contains protocol we do not support:",
		slog.Any("tokenID", tokenID),
		slog.Any("metadataURI", metadataURI))

	return nil, "", fmt.Errorf("Token metadata URI contains protocol we do not support: %v\n", metadataURI)
}

func (uc *DownloadMetadataNonFungibleTokenUseCase) executeForIPFS(tokenID *big.Int, ipfsPath string, dirPath string) (*domain.NonFungibleTokenMetadata, string, error) {
	cid := strings.Replace(ipfsPath, "ipfs://", "", -1)
	nftAsset, err := uc.repo.Get(context.Background(), cid)
	if err != nil {
		return nil, "", err
	}

	var metadata *domain.NonFungibleTokenMetadata
	if err := json.Unmarshal(nftAsset.Content, &metadata); err != nil {
		return nil, "", err
	}

	metadataFilepath := filepath.Join(dirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), "metadata.json")

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(metadataFilepath), 0755); err != nil {
		log.Fatalf("Failed create directories: %v\n", err)
	}

	if err := ioutil.WriteFile(metadataFilepath, nftAsset.Content, 0644); err != nil {
		log.Fatalf("Failed write metadata file: %v\n", err)
	}

	return metadata, metadataFilepath, nil
}

func (uc *DownloadMetadataNonFungibleTokenUseCase) executeForHTTPS(tokenID *big.Int, url string, dirPath string) (*domain.NonFungibleTokenMetadata, string, error) {
	r, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Fatalf("failed to setup get request: %v", err)
	}

	r.Header.Add("Content-Type", "application/json")

	uc.logger.Debug("Getting via HTTPS",
		slog.String("url", url),
		slog.String("method", "GET"))

	client := &http.Client{}
	res, err := client.Do(r)
	if err != nil {
		log.Fatalf("failed to do post request: %v", err)
	}

	defer res.Body.Close()

	if res.StatusCode == http.StatusNotFound {
		log.Fatalf("URL does not exist for: %v", url)
	}

	// if res.StatusCode == http.StatusBadRequest {
	//
	// }

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(res.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	var metadata *domain.NonFungibleTokenMetadata
	if err := json.NewDecoder(teeReader).Decode(&metadata); err != nil {
		uc.logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		log.Fatalf("failed to decode: %v", err)
	}

	uc.logger.Debug("Metadata retrieved")

	metadataFilepath := filepath.Join(dirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), "metadata.json")

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(metadataFilepath), 0755); err != nil {
		log.Fatalf("Failed create directories: %v\n", err)
	}

	metadataBytes, err := json.Marshal(metadata)
	if err != nil {
		log.Fatalf("failed to marshal: %v", err)
	}

	if err := ioutil.WriteFile(metadataFilepath, metadataBytes, 0644); err != nil {
		log.Fatalf("Failed write metadata file: %v\n", err)
	}

	return metadata, metadataFilepath, nil
}
