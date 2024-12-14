package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AttachmentGetUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.AttachmentRepository
}

func NewAttachmentGetUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.AttachmentRepository,
) *AttachmentGetUseCase {
	return &AttachmentGetUseCase{config, logger, repo}
}

func (uc *AttachmentGetUseCase) Execute(ctx context.Context, id primitive.ObjectID) (*domain.Attachment, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if id.IsZero() {
		e["id"] = "Attachment is required"
	} else {

	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByID(ctx, id)
}
