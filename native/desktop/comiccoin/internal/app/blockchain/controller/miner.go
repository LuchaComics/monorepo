package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"
)

func (impl *blockchainControllerImpl) runMinerOperationInBackground(ctx context.Context) {
	impl.logger.Info("miner started...")

	// Execute the miner tick on startup of this function.
	if err := impl.handleMineTimerTicker(ctx); err != nil {
		return
	}

	// Create a timer that ticks every minute
	ticker := time.NewTicker(time.Minute)

	// Start the timer in a separate goroutine
	go func() {
		for {
			select {
			case <-ticker.C:
				if err := impl.handleMineTimerTicker(ctx); err != nil {
					return
				}
			case <-ctx.Done():
				// Clean up and exit
				ticker.Stop()
				fmt.Println("Timer stopped")
				return
			}
		}
	}()
}

func (impl *blockchainControllerImpl) handleMineTimerTicker(ctx context.Context) error {
	impl.logger.Debug("miner tick")
	// slog.Uint64("nonce", signedTransaction.Nonce))

	//
	// STEP 1:
	// Fetch all our related data.
	//

	txs, err := impl.signedTransactionStorer.List(ctx)
	if err != nil {
		impl.logger.Error("failed getting list of pending signed transactions",
			slog.Any("error", err))
	}

	// Apply the transactions per block limit.
	if len(txs) > 0 {
		txs = txs[:impl.config.Blockchain.TransPerBlock]
	}

	impl.logger.Debug("miner fetched the following txs", slog.Any("txs", txs))

	//
	// STEP 2:
	// Setup our new block
	//

	//TODO: IMPL.

	//
	// STEP 3:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	//TODO: IMPL.

	//
	// STEP 4:
	// Submit to the blockchain network for verification.
	//

	//TODO: IMPL.

	//
	// STEP 5:
	// (If this record exists locally) Delete the signed transaction record
	// from our database.
	//

	//TODO: IMPL.

	return nil
}
