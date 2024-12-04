package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	ccdomain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type StorageTransactionCommitUseCase struct {
	logger                       *slog.Logger
	walletRepo                   domain.WalletRepository
	accountRepo                  domain.AccountRepository
	genesisBlockDataRepo         domain.GenesisBlockDataRepository
	blockchainStateRepo          domain.BlockchainStateRepository
	blockDataRepo                domain.BlockDataRepository
	tokenRepo                    domain.TokenRepository
	pendingSignedTransactionRepo ccdomain.PendingSignedTransactionRepository
}

func NewStorageTransactionCommitUseCase(
	logger *slog.Logger,
	r1 domain.WalletRepository,
	r2 domain.AccountRepository,
	r3 domain.GenesisBlockDataRepository,
	r4 domain.BlockchainStateRepository,
	r5 domain.BlockDataRepository,
	r6 domain.TokenRepository,
	r7 ccdomain.PendingSignedTransactionRepository,
) *StorageTransactionCommitUseCase {
	return &StorageTransactionCommitUseCase{logger, r1, r2, r3, r4, r5, r6, r7}
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
	if err := uc.genesisBlockDataRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for genesis block data",
			slog.Any("error", err))
		return err
	}
	if err := uc.blockchainStateRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for blockchain state",
			slog.Any("error", err))
		return err
	}
	if err := uc.blockDataRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for block data",
			slog.Any("error", err))
		return err
	}
	if err := uc.tokenRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for token",
			slog.Any("error", err))
		return err
	}
	if err := uc.pendingSignedTransactionRepo.CommitTransaction(); err != nil {
		uc.logger.Error("Failed committing transaction for pending signed transaction",
			slog.Any("error", err))
		return err
	}
	return nil
}
