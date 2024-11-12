package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type GetBlockchainStateFromCentralAuthorityByChainIDUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainStateDTORepository
}

func NewGetBlockchainStateFromCentralAuthorityByChainIDUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainStateDTORepository) *GetBlockchainStateFromCentralAuthorityByChainIDUseCase {
	return &GetBlockchainStateFromCentralAuthorityByChainIDUseCase{config, logger, repo}
}

func (uc *GetBlockchainStateFromCentralAuthorityByChainIDUseCase) Execute(ctx context.Context, chainID uint16) (*domain.BlockchainStateDTO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if chainID == 0 {
		e["chain_id"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validation",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Insert into database.
	//

	return uc.repo.GetFromCentralAuthorityByChainID(ctx, chainID)
}
