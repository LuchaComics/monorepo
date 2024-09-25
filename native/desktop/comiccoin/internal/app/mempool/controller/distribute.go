package controller

import (
	"context"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
)

func (impl *mempoolControllerImpl) ReadyToDistribute(ctx context.Context) (*pt_ds.SignedTransaction, error) {
	// Let us subscribe to the blockchain topic which sends pending signed
	// transactions and block execution until we receive something.
	sub := impl.messageQueueBroker.Subscribe("mempool")

	ptBytes := <-sub
	pt, err := pt_ds.NewSignedTransactionFromDeserialize(ptBytes)
	if err != nil {
		return nil, err
	}
	return pt, nil
}
