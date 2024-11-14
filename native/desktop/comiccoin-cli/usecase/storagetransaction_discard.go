package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type StorageTransactionDiscardUseCase struct {
	logger      *slog.Logger
	walletRepo  domain.WalletRepository
	accountRepo domain.AccountRepository
}

func NewStorageTransactionDiscardUseCase(
	logger *slog.Logger,
	r1 domain.WalletRepository,
	r2 domain.AccountRepository,
) *StorageTransactionDiscardUseCase {
	return &StorageTransactionDiscardUseCase{logger, r1, r2}
}

func (uc *StorageTransactionDiscardUseCase) Execute() {
	uc.accountRepo.DiscardTransaction()
	uc.walletRepo.DiscardTransaction()
}
