package usecase

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type CountTokensByOwnerUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.TokenRepository
}

func NewCountTokensByOwnerUseCase(config *config.Config, logger *slog.Logger, repo domain.TokenRepository) *CountTokensByOwnerUseCase {
	return &CountTokensByOwnerUseCase{config, logger, repo}
}

func (uc *CountTokensByOwnerUseCase) Execute(owner *common.Address) (int, error) {
	toks, err := uc.repo.ListByOwner(owner)
	if err != nil {
		return 0, err
	}
	return len(toks), nil
}
