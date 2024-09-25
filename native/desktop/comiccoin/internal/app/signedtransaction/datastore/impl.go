package datastore

import (
	"context"
	"log/slog"
)

func (impl *signedTansactionStorerImpl) Upsert(ctx context.Context, b *SignedTransaction) error {
	impl.DeleteByNonce(ctx, b.Nonce)

	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "signed--transaction-%v", b.Nonce); err != nil {
		return err
	}
	return nil
}

func (impl *signedTansactionStorerImpl) GetByNonce(ctx context.Context, nonce uint64) (*SignedTransaction, error) {
	bBytes, err := impl.dbClient.Getf("signed--transaction-%v", nonce)
	if err != nil {
		return nil, err
	}
	b, err := NewSignedTransactionFromDeserialize(bBytes)
	if err != nil {
		impl.logger.Error("failed to deserialize",
			slog.Uint64("nonce", nonce),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (impl *signedTansactionStorerImpl) List(ctx context.Context) ([]*SignedTransaction, error) {
	res := make([]*SignedTransaction, 0)
	seekThenIterateKey := ""
	err := impl.dbClient.Iterate("signed--transaction-", seekThenIterateKey, func(key, value []byte) error {
		account, err := NewSignedTransactionFromDeserialize(value)
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

func (impl *signedTansactionStorerImpl) DeleteByNonce(ctx context.Context, nonce uint64) error {
	err := impl.dbClient.Deletef("signed--transaction-%v", nonce)
	if err != nil {
		return err
	}
	return nil
}
