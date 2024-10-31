package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type DownloadMetadataNonFungibleTokenUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.IPFSRepository
}

func NewDownloadMetadataNonFungibleTokenUseCase(config *config.Config, logger *slog.Logger, r domain.IPFSRepository) *DownloadMetadataNonFungibleTokenUseCase {
	return &DownloadMetadataNonFungibleTokenUseCase{config, logger, r}
}

func (uc *DownloadMetadataNonFungibleTokenUseCase) Execute(tokenID uint64, cidString string) (*domain.NonFungibleTokenMetadata, string, error) {
	cid := strings.Replace(cidString, "ipfs://", "", -1)
	metadataBytes, _, err := uc.repo.Get(context.Background(), cid)
	if err != nil {
		return nil, "", err
	}

	var metadata *domain.NonFungibleTokenMetadata
	if err := json.Unmarshal(metadataBytes, &metadata); err != nil {
		return nil, "", err
	}

	metadataFilepath := filepath.Join(uc.config.App.DirPath, "non_fungible_token_assets", fmt.Sprintf("%v", tokenID), "metadata.json")

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(metadataFilepath), 0755); err != nil {
		log.Fatalf("Failed create directories: %v\n", err)
	}

	if err := ioutil.WriteFile(metadataFilepath, metadataBytes, 0644); err != nil {
		log.Fatalf("Failed write metadata file: %v\n", err)
	}

	return metadata, metadataFilepath, nil
}
