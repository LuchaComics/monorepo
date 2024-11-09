package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListAllWalletUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewListAllWalletUseCase(config *config.Configuration, logger *slog.Logger, repo domain.WalletRepository) *ListAllWalletUseCase {
	return &ListAllWalletUseCase{config, logger, repo}
}

func (uc *ListAllWalletUseCase) Execute(ctx context.Context) ([]*domain.Wallet, error) {
	return uc.repo.ListAll(ctx)
}
