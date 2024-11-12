package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type UpsertBlockchainStateUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainStateRepository
}

func NewUpsertBlockchainStateUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainStateRepository) *UpsertBlockchainStateUseCase {
	return &UpsertBlockchainStateUseCase{config, logger, repo}
}

func (uc *UpsertBlockchainStateUseCase) Execute(ctx context.Context, blockchainState *domain.BlockchainState) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if blockchainState == nil {
		e["genesis_block_data"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.UpsertByChainID(ctx, blockchainState)
}
