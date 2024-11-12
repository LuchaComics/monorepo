package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type GetBlockDataFromCentralAuthorityByBlockNumberUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataDTORepository
}

func NewGetBlockDataFromCentralAuthorityByBlockNumberUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataDTORepository) *GetBlockDataFromCentralAuthorityByBlockNumberUseCase {
	return &GetBlockDataFromCentralAuthorityByBlockNumberUseCase{config, logger, repo}
}

func (uc *GetBlockDataFromCentralAuthorityByBlockNumberUseCase) Execute(ctx context.Context, blockNumber uint64) (*domain.BlockDataDTO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if blockNumber == 0 {
		e["block_number"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting block data",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get from authority.
	//

	return uc.repo.GetFromCentralAuthorityByBlockNumber(ctx, blockNumber)
}
