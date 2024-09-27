package worker

import (
	"context"
)

func (impl *blockchainWorkerImpl) RunMinerOperation(ctx context.Context) {
	impl.logger.Info("miner started...")

	// // Execute the miner tick on startup of this function.
	// if err := impl.handleMineTimerTicker(ctx); err != nil {
	// 	return
	// }
	//
	// // Create a timer that ticks every minute
	// ticker := time.NewTicker(time.Minute)
	//
	// // Start the timer in a separate goroutine
	// go func() {
	// 	for {
	// 		select {
	// 		case <-ticker.C:
	// 			if err := impl.handleMineTimerTicker(ctx); err != nil {
	// 				return
	// 			}
	// 		case <-ctx.Done():
	// 			// Clean up and exit
	// 			ticker.Stop()
	// 			fmt.Println("Timer stopped")
	// 			return
	// 		}
	// 	}
	// }()
}

// func (impl *blockchainWorkerImpl) handleMineTimerTicker(ctx context.Context) error {
// 	impl.logger.Debug("miner tick")
// 	// slog.Uint64("nonce", signedTransaction.Nonce))
//
// 	//
// 	// STEP 1:
// 	// Fetch all our related data.
// 	//
//
// 	txs, err := impl.signedTransactionStorer.List(ctx)
// 	if err != nil {
// 		impl.logger.Error("failed getting list of pending signed transactions",
// 			slog.Any("error", err))
// 	}
//
// 	impl.logger.Debug("miner fetched the following txs", slog.Any("txs", txs))
//
// 	// Apply the transactions per block limit.
// 	if len(txs) > 0 {
// 		txs = txs[:impl.config.Blockchain.TransPerBlock]
// 	}
//
// 	prevBlockHash, err := impl.lastHashStorer.Get(ctx)
// 	if err != nil {
// 		return fmt.Errorf("Failed to get last hash in database: %v", err)
// 	}
// 	prevBlock, err := impl.blockDataStorer.GetByHash(ctx, prevBlockHash)
// 	if err != nil {
// 		return fmt.Errorf("Failed to get lastest block from database: %v", err)
// 	}
//
// 	//
// 	// STEP 2:
// 	// Setup our new block
// 	//
//
// 	gasPrice := uint64(impl.config.Blockchain.GasPrice)
// 	unitsOfGas := uint64(impl.config.Blockchain.UnitsOfGas)
// 	trans := make([]blockdata_ds.BlockTx, 1)
// 	for _, signedTx := range txs {
// 		blockTx := blockdata_ds.NewBlockTx(*signedTx, gasPrice, unitsOfGas)
// 		trans = append(trans, blockTx)
// 	}
//
// 	// Construct a merkle tree from the transaction for this block. The root
// 	// of this tree will be part of the block to be mined.
// 	tree, err := merkle.NewTree(trans)
// 	if err != nil {
// 		return fmt.Errorf("Failed to create merkle tree: %v", err)
// 	}
//
// 	// Construct the genesis block.
// 	block := blockdata_ds.Block{
// 		Header: blockdata_ds.BlockHeader{
// 			Number:        prevBlock.Header.Number + 1,
// 			PrevBlockHash: prevBlockHash,
// 			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
// 			Beneficiary:   prevBlock.Header.Beneficiary,
// 			Difficulty:    impl.config.Blockchain.Difficulty,
// 			MiningReward:  impl.config.Blockchain.MiningReward,
// 			// StateRoot:     "",             //args.StateRoot, // SKIP!
// 			TransRoot: tree.RootHex(), //
// 			Nonce:     0,              // Will be identified by the POW algorithm.
// 		},
// 		MerkleTree: tree,
// 	}
//
// 	//
// 	// STEP 3:
// 	// Execute the proof of work to find our nounce to meet the hash difficulty.
// 	//
//
// 	if mineErr := block.PerformPOW(ctx, impl.config.Blockchain.Difficulty); mineErr != nil {
// 		return fmt.Errorf("Failed to mine block: %v", err)
// 	}
//
// 	// Attach our nonce
// 	// block.Header.Nonce = nonce
//
// 	// Create our data.
// 	blockData := blockdata_ds.NewBlockData(block)
//
// 	impl.logger.Debug("new block data created",
// 		slog.Any("mined_nonce", block.Header.Nonce),
// 		slog.Any("mined_hash", block.Hash()),
// 		slog.Any("blockdata", blockData))
//
// 	//
// 	// STEP 4:
// 	// Submit to the blockchain network for verification.
// 	//
//
// 	//TODO: IMPL.
//
// 	//
// 	// STEP 5:
// 	// (If this record exists locally) Delete the signed transaction record
// 	// from our database.
// 	//
//
// 	//TODO: IMPL.
//
// 	return nil
// }
