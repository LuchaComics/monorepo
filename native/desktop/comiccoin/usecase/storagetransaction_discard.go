package usecase

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	ccdomain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type StorageTransactionDiscardUseCase struct {
	logger                       *slog.Logger
	walletRepo                   domain.WalletRepository
	accountRepo                  domain.AccountRepository
	genesisBlockDataRepo         domain.GenesisBlockDataRepository
	blockchainStateRepo          domain.BlockchainStateRepository
	blockDataRepo                domain.BlockDataRepository
	tokenRepo                    domain.TokenRepository
	pendingSignedTransactionRepo ccdomain.PendingSignedTransactionRepository
}

func NewStorageTransactionDiscardUseCase(
	logger *slog.Logger,
	r1 domain.WalletRepository,
	r2 domain.AccountRepository,
	r3 domain.GenesisBlockDataRepository,
	r4 domain.BlockchainStateRepository,
	r5 domain.BlockDataRepository,
	r6 domain.TokenRepository,
	r7 ccdomain.PendingSignedTransactionRepository,
) *StorageTransactionDiscardUseCase {
	return &StorageTransactionDiscardUseCase{logger, r1, r2, r3, r4, r5, r6, r7}
}

func (uc *StorageTransactionDiscardUseCase) Execute() {
	uc.accountRepo.DiscardTransaction()
	uc.walletRepo.DiscardTransaction()
	uc.genesisBlockDataRepo.DiscardTransaction()
	uc.blockchainStateRepo.DiscardTransaction()
	uc.blockDataRepo.DiscardTransaction()
	uc.tokenRepo.DiscardTransaction()
	uc.pendingSignedTransactionRepo.DiscardTransaction()
}
