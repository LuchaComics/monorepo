package service

import (
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type UploadUnassignedAttachmentService struct {
	logger                          *slog.Logger
	cloudStorageUploadUseCase       *usecase.CloudStorageUploadUseCase
	createAttachmentUseCase         *usecase.CreateAttachmentUseCase
	cloudStoragePresignedURLUseCase *usecase.CloudStoragePresignedURLUseCase
}

func NewUploadUnassignedAttachmentService(
	logger *slog.Logger,
	uc1 *usecase.CloudStorageUploadUseCase,
	uc2 *usecase.CreateAttachmentUseCase,
	uc3 *usecase.CloudStoragePresignedURLUseCase,
) *UploadUnassignedAttachmentService {
	return &UploadUnassignedAttachmentService{logger, uc1, uc2, uc3}
}

type UploadUnassignedAttachmentRequestIDO struct {
	Name        string
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
}

type UploadUnassignedAttachmentResponseIDO struct {
	Name        string
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
}

func (s *UploadUnassignedAttachmentService) Execute(sessCtx mongo.SessionContext, req *UploadUnassignedAttachmentRequestIDO) (*UploadUnassignedAttachmentResponseIDO, error) {
	// // Extract from our session the following data.
	// userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

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

	// timestamp :=
	// objectKey := fmt.Sprintf("/attachments/%v_%v", )

	// // Lookup the user in our database, else return a `400 Bad Request` error.
	// u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	// if err != nil {
	// 	s.logger.Error("database error", slog.Any("err", err))
	// 	return err
	// }
	// if u == nil {
	// 	s.logger.Warn("user does not exist validation error")
	// 	return httperror.NewForBadRequestWithSingleField("id", "does not exist")
	// }
	//
	// // Generate unique token and save it to the user record.
	// u.EmailVerificationCode = primitive.NewObjectID().Hex()
	// if err := s.userUpdateUseCase.Execute(sessCtx, u); err != nil {
	// 	s.logger.Warn("user update by id failed", slog.Any("error", err))
	// 	return err
	// }
	//
	// // Send password reset email.
	// return s.templatedEmailer.SendForgotPasswordEmail(req.Email, u.EmailVerificationCode, u.FirstName)
	return nil, nil
}
