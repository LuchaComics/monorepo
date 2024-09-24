package datastore

import (
	"context"
	"log/slog"
)

func (impl *pendingTransactionStorerImpl) Insert(ctx context.Context, b *PendingTransaction) error {
	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "pending-transaction-%v", b.ID); err != nil {
		return err
	}
	return nil
}

func (impl *pendingTransactionStorerImpl) GetByID(ctx context.Context, id string) (*PendingTransaction, error) {
	bBytes, err := impl.dbClient.Getf("pending-transaction-%v", id)
	if err != nil {
		return nil, err
	}
	b, err := NewPendingTransactionFromDeserialize(bBytes)
	if err != nil {
		impl.logger.Error("failed to deserialize",
			slog.String("id", id),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (impl *pendingTransactionStorerImpl) List(ctx context.Context) ([]*PendingTransaction, error) {
	res := make([]*PendingTransaction, 0)
	seekThenIterateKey := ""
	err := impl.dbClient.Iterate("pending-transaction-", seekThenIterateKey, func(key, value []byte) error {
		account, err := NewPendingTransactionFromDeserialize(value)
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

func (impl *pendingTransactionStorerImpl) DeleteByID(ctx context.Context, id string) error {
	err := impl.dbClient.Deletef("pending-transaction-%v", id)
	if err != nil {
		return err
	}
	return nil
}
