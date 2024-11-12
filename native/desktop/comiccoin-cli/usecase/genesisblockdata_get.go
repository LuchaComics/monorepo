package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type GetGenesisBlockDataUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.GenesisBlockDataRepository
}

func NewGetGenesisBlockDataUseCase(config *config.Config, logger *slog.Logger, repo domain.GenesisBlockDataRepository) *GetGenesisBlockDataUseCase {
	return &GetGenesisBlockDataUseCase{config, logger, repo}
}

func (uc *GetGenesisBlockDataUseCase) Execute(ctx context.Context, chainID uint16) (*domain.GenesisBlockData, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chain_id"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting genesis block data",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByChainID(ctx, chainID)
}
