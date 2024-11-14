package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type StorageTransactionOpenUseCase struct {
	logger      *slog.Logger
	walletRepo  domain.ExtendedWalletRepository
	accountRepo domain.ExtendedAccountRepository
}

func NewStorageTransactionOpenUseCase(
	logger *slog.Logger,
	r1 domain.ExtendedWalletRepository,
	r2 domain.ExtendedAccountRepository,
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
