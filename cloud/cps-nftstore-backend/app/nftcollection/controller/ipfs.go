package controller

import (
	"context"
	"log/slog"
	"time"

	domain "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
)

func (impl *NFTCollectionControllerImpl) ReprovidehCollectionsInIPNS(ctx context.Context) error {
	f := &domain.NFTCollectionPaginationListFilter{
		PageSize:  1_000_000_000,
		SortField: "id",
		SortOrder: 1, // 1=ascending | -1=descending
	}

	res, listErr := impl.NFTCollectionStorer.ListByFilter(ctx, f)
	if listErr != nil {
		impl.Logger.Error("database list by filter error", slog.Any("error", listErr))
		return listErr
	}
	if res != nil {
		for _, collection := range res.Results {
			_, publishErr := impl.IPFS.PublishToIPNS(ctx, collection.IPNSKeyName, collection.IPFSDirectoryCID)
			if publishErr != nil {
				impl.Logger.Error("failed publishing to ipns",
					slog.String("id", collection.ID.Hex()),
					slog.Any("error", publishErr))
				return publishErr
			}
			collection.ModifiedAt = time.Now()
			if updateErr := impl.NFTCollectionStorer.UpdateByID(ctx, collection); updateErr != nil {
				impl.Logger.Error("database update error", slog.Any("error", updateErr))
				return updateErr
			}
		}
	}
	return nil
}
