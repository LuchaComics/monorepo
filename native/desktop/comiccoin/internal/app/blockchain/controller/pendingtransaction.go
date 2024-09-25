package controller

import (
	"context"
	"log/slog"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/pendingtransaction/datastore"
)

func (impl *blockchainControllerImpl) GetPendingTransactions(ctx context.Context) ([]*pt_ds.PendingTransaction, error) {
	res, err := impl.pendingTransactionStorer.List(ctx)
	if err != nil {
		impl.logger.Error("failed to get the pending transactions list",
			slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
