package controller

import (
	"context"

	blockdata_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockdata/datastore"
)

func (impl *blockchainControllerImpl) GetBlockData(ctx context.Context, hash string) (*blockdata_ds.BlockData, error) {
	// return impl.blockStorer.GetByHash(ctx, hash)
	return nil, nil
}
