package controller

import (
	"context"
	"log/slog"
	"time"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config/constants"
)

func (impl *mempoolControllerImpl) RunReceiveFromNetworkOperation(ctx context.Context) error {
	// Wait until we are connected to the network...
	for {
		isConnected := impl.pubSubBroker.IsSubscriberConnectedToNetwork(ctx, constants.PubSubMempoolTopicName)
		if isConnected {
			impl.logger.Debug("Mempool connected to network")
			break
		} else {
			impl.logger.Debug("Waiting for network connection...")
			time.Sleep(10 * time.Second)
		}
	}

	// Subscribe
	subscribeChannel := impl.pubSubBroker.Subscribe(ctx, constants.PubSubMempoolTopicName)

	// Receive data from the channel.
	for {
		select {
		case signedTransactionBytes, ok := <-subscribeChannel:
			//
			// STEP 1
			// Unmarshal the signed transaction which we received from the
			// distributed pub-sub broker.
			//

			if !ok {
				impl.logger.Warn("Subscribe channel closed unexpectedly")
				return nil
			}
			signedTransaction, err := pt_ds.NewSignedTransactionFromDeserialize(signedTransactionBytes)
			if err != nil {
				impl.logger.Error("Failed to deserialize signed transaction", slog.Any("error", err))
				continue
			}
			impl.logger.Debug("Received pending signed transaction from network",
				slog.Any("received", signedTransaction))

			//
			// STEP 2
			// Validate our received transaction and proceed further if validated.
			//

			if validateErr := signedTransaction.Validate(impl.config.Blockchain.ChainID); err != nil {
				impl.logger.Error("Pending signed transaction failed validation",
					slog.Any("chain_id", impl.config.Blockchain.ChainID),
					slog.Any("error", validateErr),
				)
				return validateErr
			}

			impl.logger.Debug("Received pending signed transaction was successufully validated")

			//
			// STEP 3
			// Save to our database.
			//

			insertErr := impl.signedTransactionStorer.Upsert(ctx, signedTransaction)
			if insertErr != nil {
				impl.logger.Debug("Failed to insert (or update) the signed transaction into the database",
					slog.Any("error", insertErr))
				return insertErr
			}

			impl.logger.Debug("Saved pending signed submission to mempool",
				slog.Any("received", signedTransaction))

		case <-ctx.Done():
			impl.logger.Info("Received shutdown signal, stopping ReceiveFromNetwork")
			return nil
		}
	}
}
