package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type AttachmentUpdateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.AttachmentRepository
}

func NewAttachmentUpdateUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.AttachmentRepository,
) *AttachmentUpdateUseCase {
	return &AttachmentUpdateUseCase{config, logger, repo}
}

func (uc *AttachmentUpdateUseCase) Execute(ctx context.Context, attachment *domain.Attachment) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if attachment == nil {
		e["attachment"] = "Attachment is required"
	} else {

	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Update in our database.
	//

	return uc.repo.UpdateByID(ctx, attachment)
}
