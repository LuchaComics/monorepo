package controller

import (
	"context"
	"log"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	a_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

// LedgerController represents a ledger to record all transactions of ComicCoin.
// The purpose of this controller is to provide operations on the ledger like
// add block, get block, etc. Functionality like minting, verification, etc
// are done elsewere.
type LedgerController interface {
	NewGenesisBlock(ctx context.Context, coinbaseKey *keystore.Key) (*block_ds.Block, error)
	GetBlock(ctx context.Context, hash string) (*block_ds.Block, error)
	GetBalanceByAddress(ctx context.Context, address common.Address) (*big.Int, error)
}

type ledgerControllerImpl struct {
	logger         *slog.Logger
	accountStorer  a_ds.AccountStorer
	lastHashStorer lasthash_ds.LastHashStorer
	blockStorer    block_ds.BlockStorer
}

func NewController(cfg *config.Config, logger *slog.Logger, as a_ds.AccountStorer, lhDS lasthash_ds.LastHashStorer, blockDS block_ds.BlockStorer) LedgerController {
	// Defensive code to protect the programmer from any errors.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	return &ledgerControllerImpl{
		logger:         logger,
		accountStorer:  as,
		lastHashStorer: lhDS,
		blockStorer:    blockDS,
	}
}
