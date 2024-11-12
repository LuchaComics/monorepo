package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type ListAllWalletUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewListAllWalletUseCase(config *config.Config, logger *slog.Logger, repo domain.WalletRepository) *ListAllWalletUseCase {
	return &ListAllWalletUseCase{config, logger, repo}
}

func (uc *ListAllWalletUseCase) Execute() ([]*domain.Wallet, error) {
	return uc.repo.ListAll()
}
