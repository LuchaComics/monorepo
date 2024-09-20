package datastore

import "context"

func (impl *blockStorerImpl) Insert(ctx context.Context, b *Block) error {
	bBytes, err := b.Serialize()
	if err != nil {
		return err
	}
	if err := impl.dbClient.Setf(bBytes, "BLOCK_%v", b.Hash); err != nil {
		return err
	}
	return nil
}

func (impl *blockStorerImpl) GetByHash(ctx context.Context, hash string) (*Block, error) {
	bBytes, err := impl.dbClient.Getf("BLOCK_%v", hash)
	if err != nil {
		return nil, err
	}
	b, err := NewBlockFromDeserialize(bBytes)
	if err != nil {
		return nil, err
	}
	return b, nil
}
