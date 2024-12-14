package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	cloudinterface "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/cloud"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
)

type CloudStorageUploadUseCase struct {
	config       *config.Configuration
	logger       *slog.Logger
	cloudstorage cloudinterface.CloudStorage
}

func NewCloudStorageUploadUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	cloudstorage cloudinterface.CloudStorage,
) *CloudStorageUploadUseCase {
	return &CloudStorageUploadUseCase{config, logger, cloudstorage}
}

func (uc *CloudStorageUploadUseCase) Execute(ctx context.Context, objectKey string, dataBytes []byte, contentType string) error {
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
	// STEP 2: Upload file to cloud storage (in background).
	//
	go func(dataBin []byte, objkey string, binType string) {
		uc.logger.Debug("beginning private s3 image upload...")
		if err := uc.cloudstorage.UploadContentFromBytes(context.Background(), objkey, dataBin, binType); err != nil {
			uc.logger.Error("Cloud storage upload failure",
				slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		uc.logger.Debug("Finished cloud storage upload with success")
	}(dataBytes, objectKey, contentType)

	return nil
}
