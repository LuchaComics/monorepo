package controller

import (
	"context"
	"fmt"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	blockdata_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockdata/datastore"
	stx_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/utils/blockchain/signature"
)

func (impl *blockchainControllerImpl) NewGenesisBlock(ctx context.Context, coinbaseKey *keystore.Key) (*blockdata_ds.BlockData, error) {
	if coinbaseKey == nil {
		return nil, fmt.Errorf("missing: %v", "coinbase")
	}

	initialSupply := uint64(5000000000000000000)
	tx := &stx_ds.Transaction{
		ChainID: impl.config.Blockchain.ChainID,
		Nonce:   0, // Will be calculated later.
		From:    coinbaseKey.Address,
		To:      coinbaseKey.Address,
		Value:   initialSupply,
		Tip:     0,
		Data:    make([]byte, 0),
	}
	signedTx, err := tx.Sign(coinbaseKey.PrivateKey)
	if err != nil {
		return nil, fmt.Errorf("Failed to sign transaction: %v", err)
	}

	prevBlockHash := signature.ZeroHash

	gasPrice := uint64(impl.config.Blockchain.GasPrice)
	unitsOfGas := uint64(impl.config.Blockchain.UnitsOfGas)
	blockTx := blockdata_ds.NewBlockTx(signedTx, gasPrice, unitsOfGas)
	trans := make([]blockdata_ds.BlockTx, 1)
	trans = append(trans, blockTx)

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		return nil, fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	// Construct the genesis block.
	block := blockdata_ds.Block{
		Header: blockdata_ds.BlockHeader{
			Number:        0, // Genesis always starts at zero
			PrevBlockHash: prevBlockHash,
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			Beneficiary:   coinbaseKey.Address,
			Difficulty:    impl.config.Blockchain.Difficulty,
			MiningReward:  impl.config.Blockchain.MiningReward,
			// StateRoot:     "",             //args.StateRoot,
			TransRoot: tree.RootHex(), //
			Nonce:     0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	block.PerformPOW(ctx, impl.config.Blockchain.Difficulty)

	// Attach our nonce
	// block.Header.Nonce = nonce

	// Create our data.
	genesisBlockData := blockdata_ds.NewBlockData(block)

	// Save to database.
	if err := impl.blockDataStorer.Insert(ctx, &genesisBlockData); err != nil {
		return nil, fmt.Errorf("Failed to insert into database: %v", err)
	}
	if err := impl.lastHashStorer.Set(ctx, genesisBlockData.Hash); err != nil {
		return nil, fmt.Errorf("Failed to set last hash in database: %v", err)
	}

	// return genesisBlock, nil
	return &genesisBlockData, nil
}
