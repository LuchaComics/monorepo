package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	cloudinterface "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/cloud"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
)

type CloudStorageDeleteUseCase struct {
	config       *config.Configuration
	logger       *slog.Logger
	cloudstorage cloudinterface.CloudStorage
}

func NewCloudStorageDeleteUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	cloudstorage cloudinterface.CloudStorage,
) *CloudStorageDeleteUseCase {
	return &CloudStorageDeleteUseCase{config, logger, cloudstorage}
}

func (uc *CloudStorageDeleteUseCase) Execute(ctx context.Context, objectKeys []string) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if len(objectKeys) == 0 {
		e["object_keys"] = "Object keys are required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get presigned url
	//

	return uc.cloudstorage.DeleteByKeys(ctx, objectKeys)
}
