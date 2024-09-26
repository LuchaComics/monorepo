package controller

import (
	"context"
	"log/slog"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
)

func (impl *blockchainControllerImpl) RunMinerOperation(ctx context.Context) {
	//TODO: IMPL: If this node is authorized miner then run the following code...

	// Subscribe to the `signed-transactions` topic so we can received
	// all the latest signed transactions to mine!
	sub := impl.pubSubBroker.Subscribe(ctx, "mempool")

	for true {
		signedTransactionBytes := <-sub
		signedTransaction, err := pt_ds.NewSignedTransactionFromDeserialize(signedTransactionBytes)
		if err != nil {
			impl.logger.Error("signed transaction received",
				slog.Uint64("nonce", signedTransaction.Nonce))

			// Do not continue in this loop iteration but skip it and restart it
			// so we are waiting for the next subscription request instead of
			// crashing this function.
			continue
		}
		if miningErr := impl.mine(ctx, signedTransaction); miningErr != nil {
			impl.logger.Error("signed transaction failed mining",
				slog.Any("error", miningErr),
				slog.Uint64("nonce", signedTransaction.Nonce))
			continue
		}
	}
}

func (impl *blockchainControllerImpl) mine(ctx context.Context, signedTransaction *pt_ds.SignedTransaction) error {
	impl.logger.Debug("signed transaction received",
		slog.Uint64("nonce", signedTransaction.Nonce))

	//
	// STEP 1:
	// Fetch all our related data.
	//

	//TODO: IMPL.

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
