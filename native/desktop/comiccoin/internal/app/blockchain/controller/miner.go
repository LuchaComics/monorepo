package controller

import (
	"context"
	"log/slog"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/pendingtransaction/datastore"
)

func (impl *blockchainControllerImpl) RunMinerOperation(ctx context.Context) {
	//TODO: IMPL: If this node is authorized miner then run the following code...

	// Subscribe to the `pending-transactions` topic so we can received
	// all the latest pending transactions to mine!
	sub := impl.messageQueueBroker.Subscribe("pending-transactions")

	for true {
		pendingTransactionBytes := <-sub
		pendingTransaction, err := pt_ds.NewPendingTransactionFromDeserialize(pendingTransactionBytes)
		if err != nil {
			impl.logger.Error("pending transaction received",
				slog.String("id", pendingTransaction.ID))

			// Do not continue in this loop iteration but skip it and restart it
			// so we are waiting for the next subscription request instead of
			// crashing this function.
			continue
		}
		if miningErr := impl.mine(ctx, pendingTransaction); miningErr != nil {
			impl.logger.Error("pending transaction failed mining",
				slog.Any("error", miningErr),
				slog.String("id", pendingTransaction.ID))
			continue
		}
	}
}

func (impl *blockchainControllerImpl) mine(ctx context.Context, pendingTransaction *pt_ds.PendingTransaction) error {
	impl.logger.Debug("pending transaction received",
		slog.String("id", pendingTransaction.ID))

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
	// (If this record exists locally) Delete the pending transaction record
	// from our database.
	//

	//TODO: IMPL.

	return nil
}
