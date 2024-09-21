package controller

import (
	"context"
	"log"
	"log/slog"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type BlockchainController interface {
	NewGenesisBlock(ctx context.Context, coinbaseKey *keystore.Key) (*block_ds.Block, error)
	GetBlock(ctx context.Context, hash string) (*block_ds.Block, error)
	GetBalanceByAddress(ctx context.Context, address common.Address) (*big.Int, error)
}

type blockchainControllerImpl struct {
	logger         *slog.Logger
	lastHashStorer lasthash_ds.LastHashStorer
	blockStorer    block_ds.BlockStorer
}

func NewController(cfg *config.Config, logger *slog.Logger, lhDS lasthash_ds.LastHashStorer, blockDS block_ds.BlockStorer) BlockchainController {
	// Defensive code to protect the programmer from any errors.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	return &blockchainControllerImpl{
		logger:         logger,
		lastHashStorer: lhDS,
		blockStorer:    blockDS,
	}
}