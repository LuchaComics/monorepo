package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	ccdomain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type StorageTransactionOpenUseCase struct {
	logger                       *slog.Logger
	walletRepo                   domain.WalletRepository
	accountRepo                  domain.AccountRepository
	genesisBlockDataRepo         domain.GenesisBlockDataRepository
	blockchainStateRepo          domain.BlockchainStateRepository
	blockDataRepo                domain.BlockDataRepository
	tokenRepo                    domain.TokenRepository
	pendingSignedTransactionRepo ccdomain.PendingSignedTransactionRepository
}

func NewStorageTransactionOpenUseCase(
	logger *slog.Logger,
	r1 domain.WalletRepository,
	r2 domain.AccountRepository,
	r3 domain.GenesisBlockDataRepository,
	r4 domain.BlockchainStateRepository,
	r5 domain.BlockDataRepository,
	r6 domain.TokenRepository,
	r7 ccdomain.PendingSignedTransactionRepository,
) *StorageTransactionOpenUseCase {
	return &StorageTransactionOpenUseCase{logger, r1, r2, r3, r4, r5, r6, r7}
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
	if err := uc.genesisBlockDataRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for genesis block data",
			slog.Any("error", err))
		return err
	}
	if err := uc.blockchainStateRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for blockchain state",
			slog.Any("error", err))
		return err
	}
	if err := uc.blockDataRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for block data",
			slog.Any("error", err))
		return err
	}
	if err := uc.tokenRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for token",
			slog.Any("error", err))
		return err
	}
	if err := uc.pendingSignedTransactionRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for pending signed transaction",
			slog.Any("error", err))
		return err
	}
	return nil
}
