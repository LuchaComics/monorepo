package service

import (
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type AttachmentGarbageCollectorService struct {
	logger                        *slog.Logger
	attachmentListByFilterUseCase *usecase.AttachmentListByFilterUseCase
	attachmentDeleteUseCase       *usecase.AttachmentDeleteUseCase
	cloudStorageDeleteUseCase     *usecase.CloudStorageDeleteUseCase
}

func NewAttachmentGarbageCollectorService(
	logger *slog.Logger,
	uc1 *usecase.AttachmentListByFilterUseCase,
	uc2 *usecase.AttachmentDeleteUseCase,
	uc3 *usecase.CloudStorageDeleteUseCase,
) *AttachmentGarbageCollectorService {
	return &AttachmentGarbageCollectorService{logger, uc1, uc2, uc3}
}

func (s *AttachmentGarbageCollectorService) Execute(sessCtx mongo.SessionContext, tenantID primitive.ObjectID) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if tenantID.IsZero() {
		e["tenant_id"] = "Tenant ID is required"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get data from database.
	//

	// Look back 24 hours from now.
	nowMinus24Hours := time.Now().Add(-24 * time.Hour)

	filter := &domain.AttachmentFilter{
		TenantID:      tenantID,
		Status:        domain.AttachmentStatusActive,
		BelongsToType: domain.AttachmentBelongsToTypeUnassigned,
		CreatedAtEnd:  &nowMinus24Hours,
	}

	// Lookup the user in our database, else return a `400 Bad Request` error.
	resp, err := s.attachmentListByFilterUseCase.Execute(sessCtx, filter)
	if err != nil {
		s.logger.Error("failed listing attachments",
			slog.Any("err", err))
		return err
	}

	//
	// STEP 3: Iterate over attachments and delete.
	//

	for _, attch := range resp.Attachments {
		s.logger.Debug("Fetched old attachment",
			slog.Any("id", attch.ID),
			slog.Any("created", attch.CreatedAt),
			slog.Any("status", attch.Status))

		//
		// STEP 4: Delete uploaded data in cloud storage.
		//

		if err := s.cloudStorageDeleteUseCase.Execute(sessCtx, []string{attch.ObjectKey}); err != nil {
			s.logger.Error("Failed deleting attachment from cloud storage",
				slog.Any("err", err))
			return err
		}

		//
		// STEP 5: Delete attachment record from our database.
		//

		if err := s.attachmentDeleteUseCase.Execute(sessCtx, attch.ID); err != nil {
			s.logger.Error("failed deleting attachment",
				slog.Any("err", err))
			return err
		}

		s.logger.Debug("Deleted old attachment",
			slog.Any("id", attch.ID))
	}

	return nil
}
