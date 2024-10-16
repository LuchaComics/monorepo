package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type BroadcastProposedBlockDataDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ProposedBlockDataDTORepository
}

func NewBroadcastProposedBlockDataDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.ProposedBlockDataDTORepository) *BroadcastProposedBlockDataDTOUseCase {
	return &BroadcastProposedBlockDataDTOUseCase{config, logger, repo}
}

func (uc *BroadcastProposedBlockDataDTOUseCase) Execute(ctx context.Context, stx *domain.ProposedBlockData) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if stx == nil {
		e["signed_transaction"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our strucutre.
	//

	dto := &domain.ProposedBlockDataDTO{
		Hash:            stx.Hash,
		Header:          stx.Header,
		HeaderSignatureBytes: stx.HeaderSignatureBytes,
		Trans:           stx.Trans,
		Validator:       stx.Validator,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.BroadcastToP2PNetwork(ctx, dto)
}
