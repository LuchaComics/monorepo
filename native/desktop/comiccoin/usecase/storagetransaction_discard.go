package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type StorageTransactionDiscardUseCase struct {
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

func NewStorageTransactionDiscardUseCase(
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
) *StorageTransactionDiscardUseCase {
	return &StorageTransactionDiscardUseCase{config, logger, r1, r2, r3, r4, r5, r6, r7, r8, r9}
}

func (uc *StorageTransactionDiscardUseCase) Execute() {
	uc.accountRepo.DiscardTransaction()
	uc.tokenRepo.DiscardTransaction()
	uc.latestHashRepo.DiscardTransaction()
	uc.latestTokenIDRepo.DiscardTransaction()
	uc.blockDataRepo.DiscardTransaction()
	uc.identityKeyRepo.DiscardTransaction()
	uc.mempoolTransactionRepo.DiscardTransaction()
	uc.pendingBlockTransactionRepo.DiscardTransaction()
	uc.walletRepo.DiscardTransaction()
	// uc.signedTransactionRepo.DiscardTransaction()
}
