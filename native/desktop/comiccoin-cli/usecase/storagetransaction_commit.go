package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type StorageTransactionCommitUseCase struct {
	logger      *slog.Logger
	walletRepo  domain.ExtendedWalletRepository
	accountRepo domain.ExtendedAccountRepository
}

func NewStorageTransactionCommitUseCase(
	logger *slog.Logger,
	r1 domain.ExtendedWalletRepository,
	r2 domain.ExtendedAccountRepository,
) *StorageTransactionCommitUseCase {
	return &StorageTransactionCommitUseCase{logger, r1, r2}
}

func (uc *StorageTransactionCommitUseCase) Execute() error {
	if err := uc.accountRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for accounts",
			slog.Any("error", err))
		return err
	}
	if err := uc.walletRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for wallet",
			slog.Any("error", err))
		return err
	}
	return nil
}
