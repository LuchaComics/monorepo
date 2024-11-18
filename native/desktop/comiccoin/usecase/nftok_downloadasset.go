package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/slog"
	"math/big"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type DownloadNonFungibleTokenAssetUseCase struct {
	logger *slog.Logger
	repo   domain.NFTAssetRepository
}

func NewDownloadNonFungibleTokenAssetUseCase(logger *slog.Logger, r domain.NFTAssetRepository) *DownloadNonFungibleTokenAssetUseCase {
	return &DownloadNonFungibleTokenAssetUseCase{logger, r}
}

func (uc *DownloadNonFungibleTokenAssetUseCase) Execute(tokenID *big.Int, assetURI string, dirPath string) (string, error) {
	if assetURI == "" {
		uc.logger.Warn("No asset to download, skipping function...",
			slog.Any("tokenID", tokenID),
			slog.Any("asset_uri", assetURI))
		return "", nil
	}

	// Confirm URI is using protocol our app supports.
	if strings.Contains(assetURI, "ipfs://") {
		uc.logger.Debug("Downloading asset via ipfs...",
			slog.Any("tokenID", tokenID),
			slog.Any("asset_uri", assetURI))
		return uc.executeForIPFS(tokenID, assetURI, dirPath)
	} else if strings.Contains(assetURI, "https://") {
		uc.logger.Debug("Downloading asset via https...",
			slog.Any("tokenID", tokenID),
			slog.Any("asset_uri", assetURI))

		return uc.executeForHTTP(tokenID, assetURI, dirPath)
	}

	uc.logger.Error("Token asset uri contains protocol we do not support:",
		slog.Any("tokenID", tokenID),
		slog.Any("asset_uri", assetURI))

	return "", fmt.Errorf("Token asset URI contains protocol we do not support: %v\n", assetURI)
}

func (uc *DownloadNonFungibleTokenAssetUseCase) executeForIPFS(tokenID *big.Int, assetIpfsPath string, dirPath string) (string, error) {
	assetCID := strings.Replace(assetIpfsPath, "ipfs://", "", -1)
	nftAsset, err := uc.repo.Get(context.Background(), assetCID)
	if err != nil {
		return "", err
	}

	var filename string = nftAsset.Filename
	assetFilepath := filepath.Join(dirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), filename)

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(assetFilepath), 0755); err != nil {
		return "", fmt.Errorf("Failed create directories: %v\n", err)
	}

	// Convert the response bytes into reader.
	contentBytesReader := bytes.NewReader(nftAsset.Content)

	// Save the data to file.
	f, err := os.Create(assetFilepath)
	if err != nil {
		return "", fmt.Errorf("Failed create file: %v\n", err)
	}
	defer f.Close()

	// Save to local directory.
	_, err = io.Copy(f, contentBytesReader)
	if err != nil {
		return "", fmt.Errorf("Failed to copy into file contents %v\n", err)
	}

	return assetFilepath, nil
}

func (uc *DownloadNonFungibleTokenAssetUseCase) executeForHTTP(tokenID *big.Int, url string, dirPath string) (string, error) {
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

	// Read the response body from the `res` variable and store it in the `contentBytes` variable.
	contentBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to_read_response_body: %v\n", err)
	}

	// Get the filename at the end of the URL path (special thanks: https://stackoverflow.com/a/44570361).
	filename := path.Base(r.URL.Path)

	assetFilepath := filepath.Join(dirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), filename)

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(assetFilepath), 0755); err != nil {
		return "", fmt.Errorf("Failed create directories: %v\n", err)
	}

	// Convert the response bytes into reader.
	contentBytesReader := bytes.NewReader(contentBytes)

	// Save the data to file.
	f, err := os.Create(assetFilepath)
	if err != nil {
		return "", fmt.Errorf("Failed create file: %v\n", err)
	}
	defer f.Close()

	// Save to local directory.
	_, err = io.Copy(f, contentBytesReader)
	if err != nil {
		return "", fmt.Errorf("Failed to copy into file contents %v\n", err)
	}

	return assetFilepath, nil
}
