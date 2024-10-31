package usecase

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
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

func (uc *DownloadNonFungibleTokenAssetUseCase) Execute(tokenID uint64, assetCIDString string) (string, error) {
	assetCID := strings.Replace(assetCIDString, "ipfs://", "", -1)
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
