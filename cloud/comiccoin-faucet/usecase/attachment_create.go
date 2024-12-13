package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type CreateAttachmentUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.AttachmentRepository
}

func NewCreateAttachmentUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.AttachmentRepository,
) *CreateAttachmentUseCase {
	return &CreateAttachmentUseCase{config, logger, repo}
}

func (uc *CreateAttachmentUseCase) Execute(ctx context.Context, attachment *domain.Attachment) error {
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
	// STEP 3: Insert into database.
	//

	return uc.repo.Create(ctx, attachment)
}
