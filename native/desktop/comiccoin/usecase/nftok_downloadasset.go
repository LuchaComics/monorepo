package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type DownloadNonFungibleTokenAssetUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.IPFSRepository
}

func NewDownloadNonFungibleTokenAssetUseCase(config *config.Config, logger *slog.Logger, r domain.IPFSRepository) *DownloadNonFungibleTokenAssetUseCase {
	return &DownloadNonFungibleTokenAssetUseCase{config, logger, r}
}

func (uc *DownloadNonFungibleTokenAssetUseCase) Execute(tokenID uint64, assetURI string) (string, error) {
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
		return uc.executeForIPFS(tokenID, assetURI)
	} else if strings.Contains(assetURI, "https://") {
		uc.logger.Debug("Downloading asset via https...",
			slog.Any("tokenID", tokenID),
			slog.Any("asset_uri", assetURI))

		return uc.executeForHTTP(tokenID, assetURI)
	}

	uc.logger.Error("Token asset uri contains protocol we do not support:",
		slog.Any("tokenID", tokenID),
		slog.Any("asset_uri", assetURI))

	return "", fmt.Errorf("Token asset URI contains protocol we do not support: %v\n", assetURI)
}

func (uc *DownloadNonFungibleTokenAssetUseCase) executeForIPFS(tokenID uint64, assetIpfsPath string) (string, error) {
	assetCID := strings.Replace(assetIpfsPath, "ipfs://", "", -1)
	contentBytes, contentType, err := uc.repo.Get(context.Background(), assetCID)
	if err != nil {
		return "", err
	}

	var filename string
	switch contentType {
	case "image/png": //  Portable Network Graphics (PNG)
		filename = fmt.Sprintf("%v.png", assetCID)
	case "image/jpeg": // Joint Photographic Experts Group (JPEG)
		filename = fmt.Sprintf("%v.jpg", assetCID)
	case "image/jpg": // Joint Photographic Experts Group (JPG)
		filename = fmt.Sprintf("%v.jpg", assetCID)
	case "image/gif": //  Graphics Interchange Format (GIF)
		filename = fmt.Sprintf("%v.gif", assetCID)
	case "image/bmp": // Bitmap (BMP)
		filename = fmt.Sprintf("%v.bmp", assetCID)
	case "image/tiff": // Tagged Image File Format (TIFF)
		filename = fmt.Sprintf("%v.tiff", assetCID)
	case "image/tif": // Tagged Image File Format (TIFF)
		filename = fmt.Sprintf("%v.tif", assetCID)
	case "image/webp": // Web Picture (WEBP)
		filename = fmt.Sprintf("%v.webp", assetCID)
	case "image/svg+xml": // Scalable Vector Graphics (SVG)
		filename = fmt.Sprintf("%v.svg", assetCID)
	case "video/mp4": // MPEG-4
		filename = fmt.Sprintf("%v.mp4", assetCID)
	case "video/x-m4v": // MPEG-4 Video
		filename = fmt.Sprintf("%v.m4v", assetCID)
	case "video/quicktime": // QuickTime
		filename = fmt.Sprintf("%v.mov", assetCID)
	case "video/webm": // WebM
		filename = fmt.Sprintf("%v.webm", assetCID)
	case "video/ogg": // Ogg Theora
		filename = fmt.Sprintf("%v.ogv", assetCID)
	case "video/x-flv": // Flash Video
		filename = fmt.Sprintf("%v.flv", assetCID)
	case "video/x-msvideo": // AVI (Audio Video Interleave)
		filename = fmt.Sprintf("%v.avi", assetCID)
	case "video/x-ms-wmv": //  Windows Media Video
		filename = fmt.Sprintf("%v.wmv", assetCID)
	case "video/3gpp": // 3GPP (3rd Generation Partnership Project)
		filename = fmt.Sprintf("%v.3gp", assetCID)
	case "video/3gpp2": // 3GPP2 (3rd Generation Partnership Project 2)
		filename = fmt.Sprintf("%v.3g2", assetCID)
	default:
		return "", fmt.Errorf("Unsupported file type: %v.", contentType)
	}

	assetFilepath := filepath.Join(uc.config.App.DirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), filename)

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

func (uc *DownloadNonFungibleTokenAssetUseCase) executeForHTTP(tokenID uint64, url string) (string, error) {
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

	assetFilepath := filepath.Join(uc.config.App.DirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), filename)

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
