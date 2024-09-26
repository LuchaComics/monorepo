package datastore

import (
	"context"
	"log/slog"
)

func (impl *blockdataStorerImpl) Insert(ctx context.Context, b *BlockData) error {
	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "block-data-%v", b.Hash); err != nil {
		return err
	}
	return nil
}

func (impl *blockdataStorerImpl) GetByHash(ctx context.Context, hash string) (*BlockData, error) {
	bBytes, err := impl.dbClient.Getf("block-data-%v", hash)
	if err != nil {
		return nil, err
	}
	b, err := NewBlockDataFromDeserialize(bBytes)
	if err != nil {
		impl.logger.Error("failed to deserialize",
			slog.String("hash", hash),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}
