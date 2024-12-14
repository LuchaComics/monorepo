package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type AttachmentCreateService struct {
	config                          *config.Configuration
	logger                          *slog.Logger
	cloudStorageUploadUseCase       *usecase.CloudStorageUploadUseCase
	createAttachmentUseCase         *usecase.CreateAttachmentUseCase
	cloudStoragePresignedURLUseCase *usecase.CloudStoragePresignedURLUseCase
}

func NewAttachmentCreateService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc1 *usecase.CloudStorageUploadUseCase,
	uc2 *usecase.CreateAttachmentUseCase,
	uc3 *usecase.CloudStoragePresignedURLUseCase,
) *AttachmentCreateService {
	return &AttachmentCreateService{cfg, logger, uc1, uc2, uc3}
}

type AttachmentCreateRequestIDO struct {
	Name        string
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
}

type AttachmentCreateResponseIDO domain.Attachment

func (s *AttachmentCreateService) Execute(sessCtx mongo.SessionContext, req *AttachmentCreateRequestIDO) (*AttachmentCreateResponseIDO, error) {
	//
	// STEP 1: Validation
	//

	e := make(map[string]string)

	if req == nil {
		err := errors.New("No request data inputted")
		s.logger.Error("validation error", slog.Any("err", err))
		return nil, err
	}
	if req.Filename == "" {
		e["filename"] = "missing value"
	}
	if req.ContentType == "" {
		e["content_type"] = "missing value"
	}
	if req.Data == nil {
		e["data"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("validation failure",
			slog.Any("e", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Upload to cloud storage.
	//

	timestamp := uint64(time.Now().UTC().UnixMilli())
	objectKey := fmt.Sprintf("attachments/%v_%v", timestamp, req.Filename)

	s.logger.Debug("Uploading to cloud storage...",
		slog.String("object_id", objectKey),
		slog.String("filename", req.Filename),
		slog.String("content_type", req.ContentType),
		slog.Int("content_length", len(req.Data)),
	)

	if err := s.cloudStorageUploadUseCase.Execute(sessCtx, objectKey, req.Data, req.ContentType); err != nil {
		s.logger.Error("Failed uploading to cloud storage", slog.Any("err", err))
		return nil, err
	}

	//
	// STEP 3: Create the record in our database.
	//

	// Extract from our session the following data.
	userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName, _ := sessCtx.Value(constants.SessionUserName).(string)

	attach := &domain.Attachment{
		ID:                        primitive.NewObjectID(),
		CreatedAt:                 time.Now(),
		CreatedByUserName:         userName,
		CreatedByUserID:           userID,
		ModifiedAt:                time.Now(),
		ModifiedByUserName:        userName,
		ModifiedByUserID:          userID,
		Name:                      req.Name,
		Description:               "{}",
		Filename:                  req.Filename,
		ObjectKey:                 objectKey,
		ObjectURL:                 "",
		Status:                    domain.AttachmentStatusActive,
		ContentType:               req.ContentType,
		UserID:                    userID,
		TenantID:                  s.config.App.TenantID,
		BelongsToUniqueIdentifier: primitive.NilObjectID,
		BelongsToType:             domain.AttachmentBelongsToTypeUnassigned,
	}

	if err := s.createAttachmentUseCase.Execute(sessCtx, attach); err != nil {
		s.logger.Error("Failed uploading to cloud storage", slog.Any("err", err))
		return nil, err
	}

	//
	// STEP 4: Get presigned URL.
	//

	presignedURL, err := s.cloudStoragePresignedURLUseCase.Execute(sessCtx, objectKey)
	if err != nil {
		s.logger.Error("Failed generating presigned url via cloud storage", slog.Any("err", err))
		return nil, err
	}

	attach.ObjectURL = presignedURL

	return (*AttachmentCreateResponseIDO)(attach), nil
}
