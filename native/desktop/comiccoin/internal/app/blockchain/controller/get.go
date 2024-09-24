package controller

import (
	"context"

	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
)

func (impl *blockchainControllerImpl) GetBlock(ctx context.Context, hash string) (*block_ds.Block, error) {
	return impl.blockStorer.GetByHash(ctx, hash)
}
