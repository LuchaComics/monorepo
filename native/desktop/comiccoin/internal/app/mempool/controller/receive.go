package controller

import (
	"context"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
)

func (impl *mempoolControllerImpl) Receive(ctx context.Context, pendingTx *pt_ds.SignedTransaction) error {
	return nil
}
