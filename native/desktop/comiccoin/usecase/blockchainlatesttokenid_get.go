package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type GetBlockchainLastestTokenIDUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestTokenIDRepository
}

func NewGetBlockchainLastestTokenIDUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestTokenIDRepository) *GetBlockchainLastestTokenIDUseCase {
	return &GetBlockchainLastestTokenIDUseCase{config, logger, repo}
}

func (uc *GetBlockchainLastestTokenIDUseCase) Execute() (uint64, error) {
	latestTokenID, err := uc.repo.Get()
	if err != nil {
		uc.logger.Warn("Failed getting latest token id, automatically setting returned value to zero",
			slog.Any("error", err))
		latestTokenID = 0
	}
	return latestTokenID, nil
}
