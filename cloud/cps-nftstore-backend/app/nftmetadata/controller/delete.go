package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTMetadataControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	d, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if d == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "nftmetadata does not exist")
	}

	// for _, cid := range d.MetadataFileCIDs {
	// 	if err := impl.IPFS.Unpin(ctx, cid); err != nil {
	// 		impl.Logger.Error("ipfs failed unpinning file error", slog.Any("error", err))
	// 		return err
	// 	}
	// }
	// impl.Logger.Debug("deleted all nft metadata files from ipfs for this nftmetadata")
	// if err := impl.IPFS.Unpin(ctx, d.IPFSDirectoryCID); err != nil {
	// 	impl.Logger.Error("ipfs failed unpinning directory error", slog.Any("error", err))
	// 	return err
	// }

	if err := impl.NFTMetadataStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
