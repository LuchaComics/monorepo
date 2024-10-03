package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

type GetBlockchainLastestHashUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockchainLastestHashRepository
}

func NewGetBlockchainLastestHashUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockchainLastestHashRepository) *GetBlockchainLastestHashUseCase {
	return &GetBlockchainLastestHashUseCase{config, logger, repo}
}

func (uc *GetBlockchainLastestHashUseCase) Execute() (string, error) {
	return uc.repo.Get()
}
