package datastore

import (
	"context"
	"log/slog"
)

func (impl *pendingTransactionStorerImpl) Insert(ctx context.Context, b *SignedPendingTransaction) error {
	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "signed-pending-transaction-%v", b.Nonce); err != nil {
		return err
	}
	return nil
}

func (impl *pendingTransactionStorerImpl) GetByNonce(ctx context.Context, nonce uint64) (*SignedPendingTransaction, error) {
	bBytes, err := impl.dbClient.Getf("signed-pending-transaction-%v", nonce)
	if err != nil {
		return nil, err
	}
	b, err := NewSignedPendingTransactionFromDeserialize(bBytes)
	if err != nil {
		impl.logger.Error("failed to deserialize",
			slog.Uint64("nonce", nonce),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (impl *pendingTransactionStorerImpl) List(ctx context.Context) ([]*SignedPendingTransaction, error) {
	res := make([]*SignedPendingTransaction, 0)
	seekThenIterateKey := ""
	err := impl.dbClient.Iterate("signed-pending-transaction-", seekThenIterateKey, func(key, value []byte) error {
		account, err := NewSignedPendingTransactionFromDeserialize(value)
		if err != nil {
			impl.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, account)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (impl *pendingTransactionStorerImpl) DeleteByNonce(ctx context.Context, nonce uint64) error {
	err := impl.dbClient.Deletef("signed-pending-transaction-%v", nonce)
	if err != nil {
		return err
	}
	return nil
}
