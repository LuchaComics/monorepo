package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type StorageTransactionDiscardUseCase struct {
	logger      *slog.Logger
	walletRepo  domain.ExtendedWalletRepository
	accountRepo domain.ExtendedAccountRepository
}

func NewStorageTransactionDiscardUseCase(
	logger *slog.Logger,
	r1 domain.ExtendedWalletRepository,
	r2 domain.ExtendedAccountRepository,
) *StorageTransactionDiscardUseCase {
	return &StorageTransactionDiscardUseCase{logger, r1, r2}
}

func (uc *StorageTransactionDiscardUseCase) Execute() {
	uc.accountRepo.DiscardTransaction()
	uc.walletRepo.DiscardTransaction()
}
