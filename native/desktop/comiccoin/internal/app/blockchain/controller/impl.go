package controller

import (
	"context"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
)

func (impl *blockchainControllerImpl) NewGenesisBlock(ctx context.Context, coinbaseKey *keystore.Key) (*block_ds.Block, error) {
	if coinbaseKey == nil {
		return nil, fmt.Errorf("missing: %v", "coinbase")
	}

	initialSupply, _ := new(big.Int).SetString("50000000000000000000", 10) // 50 coins with 18 decimal places
	genesisTransaction := block_ds.NewTransaction(
		common.Address{},    // From: zero address for genesis block
		coinbaseKey.Address, // To: coinbase address (usually the miner's address)
		initialSupply,       // 50 coins (with 18 decimal places)
		[]byte("Genesis Block"),
		0, // Nonce
	)

	if err := genesisTransaction.Sign(coinbaseKey.PrivateKey); err != nil {
		return nil, fmt.Errorf("Failed to sign transaction: %v", err)
	}

	genesisBlock := &block_ds.Block{
		Hash:         "",
		PreviousHash: "",
		Timestamp:    time.Now(),
		Nonce:        0,
		Difficulty:   1,
		Transactions: []*block_ds.Transaction{genesisTransaction},
	}

	miningDifficulty := 1

	nonce, hash := genesisBlock.Mine(miningDifficulty)

	// Update our genesis block with our later mining information.
	genesisBlock.Nonce = nonce
	genesisBlock.Hash = hash

	// Save to database.
	if err := impl.blockStorer.Insert(ctx, genesisBlock); err != nil {
		return nil, fmt.Errorf("Failed to insert into database: %v", err)
	}
	if err := impl.lastHashStorer.Set(ctx, genesisBlock.Hash); err != nil {
		return nil, fmt.Errorf("Failed to set last hash in database: %v", err)
	}

	return genesisBlock, nil
}

func (impl *blockchainControllerImpl) GetBlock(ctx context.Context, hash string) (*block_ds.Block, error) {
	return impl.blockStorer.GetByHash(ctx, hash)
}
