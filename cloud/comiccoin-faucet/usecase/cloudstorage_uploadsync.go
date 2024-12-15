package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	cloudinterface "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/cloud"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
)

type CloudStorageSyncUploadUseCase struct {
	config       *config.Configuration
	logger       *slog.Logger
	cloudstorage cloudinterface.CloudStorage
}

func NewCloudStorageSyncUploadUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	cloudstorage cloudinterface.CloudStorage,
) *CloudStorageSyncUploadUseCase {
	return &CloudStorageSyncUploadUseCase{config, logger, cloudstorage}
}

func (uc *CloudStorageSyncUploadUseCase) Execute(ctx context.Context, objectKey string, dataBytes []byte, contentType string) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if objectKey == "" {
		e["object_key"] = "Object key is required"
	}
	if dataBytes == nil {
		e["dataBytes"] = "Data is required"
	}
	if contentType == "" {
		e["content_type"] = "Content type is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Upload file to cloud storage (synchronously).
	//
	uc.logger.Debug("beginning private s3 image upload...")
	if err := uc.cloudstorage.UploadContentFromBytes(ctx, objectKey, dataBytes, contentType); err != nil {
		uc.logger.Error("Cloud storage upload failure",
			slog.Any("error", err))
		return err
	}
	uc.logger.Debug("Finished cloud storage upload with success")

	return nil
}
