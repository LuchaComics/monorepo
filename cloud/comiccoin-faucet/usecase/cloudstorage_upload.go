package usecase

import (
	"context"
	"log/slog"
	"mime/multipart"

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

func (uc *CloudStorageUploadUseCase) Execute(ctx context.Context, objectKey string, file multipart.File) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if objectKey == "" {
		e["object_key"] = "Object key is required"
	}
	if file == nil {
		e["file"] = "File is required"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Upload file to cloud storage (in background).
	//
	go func(multipartfile multipart.File, objkey string) {
		uc.logger.Debug("beginning private s3 image upload...")
		if err := uc.cloudstorage.UploadContentFromMulipart(context.Background(), objkey, multipartfile); err != nil {
			uc.logger.Error("Cloud storage upload failure",
				slog.Any("error", err))
			// Do not return an error, simply continue this function as there might
			// be a case were the file was removed on the s3 bucket by ourselves
			// or some other reason.
		}
		uc.logger.Debug("Finished cloud storage upload with success")
	}(file, objectKey)

	return nil
}
