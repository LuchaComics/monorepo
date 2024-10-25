package usecase

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListTokensByOwnerUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewListTokensByOwnerUseCase(config *config.Config, logger *slog.Logger, repo domain.TokenRepository) *ListTokensByOwnerUseCase {
	return &ListTokensByOwnerUseCase{config, logger, repo}
}

func (uc *ListTokensByOwnerUseCase) Execute(owner *common.Address) ([]*domain.Token, error) {
	return uc.repo.ListByOwner(owner)
}
