package controller

import (
	"context"
	"log/slog"

	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
)

func (impl *blockchainControllerImpl) GetSignedTransactions(ctx context.Context) ([]*pt_ds.SignedTransaction, error) {
	res, err := impl.signedTransactionStorer.List(ctx)
	if err != nil {
		impl.logger.Error("failed to get the signed transactions list",
			slog.Any("error", err))
		return nil, err
	}
	return res, nil
}
