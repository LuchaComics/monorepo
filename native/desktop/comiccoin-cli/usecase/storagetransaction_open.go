package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type StorageTransactionOpenUseCase struct {
	logger      *slog.Logger
	walletRepo  domain.WalletRepository
	accountRepo domain.AccountRepository
}

func NewStorageTransactionOpenUseCase(
	logger *slog.Logger,
	r1 domain.WalletRepository,
	r2 domain.AccountRepository,
) *StorageTransactionOpenUseCase {
	return &StorageTransactionOpenUseCase{logger, r1, r2}
}

func (uc *StorageTransactionOpenUseCase) Execute() error {
	if err := uc.accountRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for accounts",
			slog.Any("error", err))
		return err
	}
	if err := uc.walletRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for wallet",
			slog.Any("error", err))
		return err
	}
	return nil
}
