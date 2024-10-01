package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type BroadcastPurposedBlockDataDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.PurposedBlockDataDTORepository
}

func NewBroadcastPurposedBlockDataDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.PurposedBlockDataDTORepository) *BroadcastPurposedBlockDataDTOUseCase {
	return &BroadcastPurposedBlockDataDTOUseCase{config, logger, repo}
}

func (uc *BroadcastPurposedBlockDataDTOUseCase) Execute(ctx context.Context, stx *domain.PurposedBlockData) error {
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

	dto := &domain.PurposedBlockDataDTO{
		Hash:   stx.Hash,
		Header: stx.Header,
		Trans:  stx.Trans,
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.BroadcastToP2PNetwork(ctx, dto)
}
