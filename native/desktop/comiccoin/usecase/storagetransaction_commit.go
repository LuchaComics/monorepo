package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type StorageTransactionCommitUseCase struct {
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

func NewStorageTransactionCommitUseCase(
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
	// r9 domain.SignedTransactionRepository,
) *StorageTransactionCommitUseCase {
	return &StorageTransactionCommitUseCase{config, logger, r1, r2, r3, r4, r5, r6, r7, r8, r9}
}

func (uc *StorageTransactionCommitUseCase) Execute() error {
	if err := uc.accountRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.tokenRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.latestHashRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.latestTokenIDRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.blockDataRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.identityKeyRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.mempoolTransactionRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.pendingBlockTransactionRepo.CommitTransaction(); err != nil {
		return err
	}
	if err := uc.walletRepo.CommitTransaction(); err != nil {
		return err
	}
	return nil
}
