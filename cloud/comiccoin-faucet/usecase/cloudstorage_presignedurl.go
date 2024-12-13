package usecase

import (
	"context"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	cloudinterface "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/cloud"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
)

type CloudStoragePresignedURLUseCase struct {
	config       *config.Configuration
	logger       *slog.Logger
	cloudstorage cloudinterface.CloudStorage
}

func NewCloudStoragePresignedURLUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	cloudstorage cloudinterface.CloudStorage,
) *CloudStoragePresignedURLUseCase {
	return &CloudStoragePresignedURLUseCase{config, logger, cloudstorage}
}

func (uc *CloudStoragePresignedURLUseCase) Execute(ctx context.Context, objectKey string) (string, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if objectKey == "" {
		e["object_key"] = "Object key is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return "", httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get presigned url
	//

	return uc.cloudstorage.GetPresignedURL(ctx, objectKey, 5*time.Minute)
}
