package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type ListAllWalletUseCase struct {
	logger *slog.Logger
	repo   domain.WalletRepository
}

func NewListAllWalletUseCase(logger *slog.Logger, repo domain.WalletRepository) *ListAllWalletUseCase {
	return &ListAllWalletUseCase{logger, repo}
}

func (uc *ListAllWalletUseCase) Execute(ctx context.Context) ([]*domain.Wallet, error) {
	return uc.repo.ListAll(ctx)
}
