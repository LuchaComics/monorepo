package datastore

import (
	"context"
	"log/slog"
)

func (impl *accountStorerImpl) Insert(ctx context.Context, b *Account) error {
	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "account-%v", b.Name); err != nil {
		return err
	}
	return nil
}

func (impl *accountStorerImpl) GetByName(ctx context.Context, name string) (*Account, error) {
	bBytes, err := impl.dbClient.Getf("account-%v", name)
	if err != nil {
		return nil, err
	}
	b, err := NewAccountFromDeserialize(bBytes)
	if err != nil {
		impl.logger.Error("failed to deserialize",
			slog.String("name", name),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (impl *accountStorerImpl) List(ctx context.Context) ([]*Account, error) {
	return nil, nil
}
