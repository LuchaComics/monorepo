package controller

import (
	"context"
	"log"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	a_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

// BlockchainController provides all the functionality that can be performed
// on the `ComicCoin` cryptocurrency.

type BlockchainController interface {
	NewGenesisBlock(ctx context.Context, coinbaseKey *keystore.Key) (*block_ds.Block, error)
	GetBlock(ctx context.Context, hash string) (*block_ds.Block, error)
	GetBalanceByAccountName(ctx context.Context, accountName string) (*BlockchainBalanceResponseIDO, error)
}

type blockchainControllerImpl struct {
	logger         *slog.Logger
	accountStorer  a_ds.AccountStorer
	lastHashStorer lasthash_ds.LastHashStorer
	blockStorer    block_ds.BlockStorer
}

func NewController(
	cfg *config.Config,
	logger *slog.Logger,
	as a_ds.AccountStorer,
	lhDS lasthash_ds.LastHashStorer,
	blockDS block_ds.BlockStorer,
) BlockchainController {
	// Defensive code to protect the programmer from any errors.
	if cfg.BlockchainDifficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}

	return &blockchainControllerImpl{
		logger:         logger,
		accountStorer:  as,
		lastHashStorer: lhDS,
		blockStorer:    blockDS,
	}
}
