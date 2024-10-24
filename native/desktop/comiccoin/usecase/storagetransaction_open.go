package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type StorageTransactionOpenUseCase struct {
	config                      *config.Config
	logger                      *slog.Logger
	accountRepo                 domain.AccountRepository
	tokenRepo                   domain.TokenRepository
	latestHashRepo              domain.BlockchainLastestHashRepository
	latestTokenIDRepo           domain.BlockchainLastestTokenIDRepository
	blockDataRepo               domain.BlockDataRepository
	identityKeyRepo             domain.IdentityKeyRepository
	mempoolTransactionRepo      domain.MempoolTransactionRepository
	pendingBlockTransactionRepo domain.PendingBlockTransactionRepository
	walletRepo                  domain.WalletRepository
	// signedTransactionRepo       domain.SignedTransactionRepository
}

func NewStorageTransactionOpenUseCase(
	config *config.Config,
	logger *slog.Logger,
	r1 domain.AccountRepository,
	r2 domain.TokenRepository,
	r3 domain.BlockchainLastestHashRepository,
	r4 domain.BlockchainLastestTokenIDRepository,
	r5 domain.BlockDataRepository,
	r6 domain.IdentityKeyRepository,
	r7 domain.MempoolTransactionRepository,
	r8 domain.PendingBlockTransactionRepository,
	r9 domain.WalletRepository,
	// r10 domain.SignedTransactionRepository,
) *StorageTransactionOpenUseCase {
	return &StorageTransactionOpenUseCase{config, logger, r1, r2, r3, r4, r5, r6, r7, r8, r9}
}

func (uc *StorageTransactionOpenUseCase) Execute() error {
	if err := uc.accountRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for accounts",
			slog.Any("error", err))
		return err
	}
	if err := uc.tokenRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for tokens",
			slog.Any("error", err))
		return err
	}
	if err := uc.latestHashRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for latest hash",
			slog.Any("error", err))
		return err
	}
	if err := uc.latestTokenIDRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for token id",
			slog.Any("error", err))
		return err
	}
	if err := uc.blockDataRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for block data",
			slog.Any("error", err))
		return err
	}
	if err := uc.identityKeyRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for identity key",
			slog.Any("error", err))
		return err
	}
	if err := uc.mempoolTransactionRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for mempool transaction",
			slog.Any("error", err))
		return err
	}
	if err := uc.pendingBlockTransactionRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for pending block transaction",
			slog.Any("error", err))
		return err
	}
	if err := uc.walletRepo.OpenTransaction(); err != nil {
		uc.logger.Error("Failed opening transaction for wallet",
			slog.Any("error", err))
		return err
	}
	// if err := uc.signedTransactionRepo.OpenTransaction(); err != nil {
	// 	return err
	// }
	return nil
}
