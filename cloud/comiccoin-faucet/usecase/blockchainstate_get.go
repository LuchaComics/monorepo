package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

//
// Copied from `github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase`
//

type GetBlockchainStateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.BlockchainStateRepository
}

func NewGetBlockchainStateUseCase(
	config *config.Configuration,
	logger *slog.Logger,
	repo domain.BlockchainStateRepository,
) *GetBlockchainStateUseCase {
	return &GetBlockchainStateUseCase{config, logger, repo}
}

func (uc *GetBlockchainStateUseCase) Execute(ctx context.Context, chainID uint16) (*domain.BlockchainState, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chain_id"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed getting blockchain state",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetByChainID(ctx, chainID)
}
