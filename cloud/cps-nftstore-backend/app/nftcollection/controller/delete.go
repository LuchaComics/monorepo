package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTCollectionControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	d, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if d == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "collection does not exist")
	}

	//
	// STEP 1
	// List through all the NFT metadata and delete them from IPFS network and
	// from the database.
	//

	metadataRes, listErr := impl.NFTStorer.ListByNFTCollectionID(ctx, id)
	if listErr != nil {
		impl.Logger.Error("database list by nft collection error", slog.Any("error", listErr))
		return listErr
	}
	for _, metadata := range metadataRes.Results {
		if err := impl.IPFS.Unpin(ctx, metadata.FileCID); err != nil {
			impl.Logger.Error("ipfs failed unpinning file error", slog.Any("error", err))
			return err
		}
		if deleteErr := impl.NFTStorer.DeleteByID(ctx, metadata.ID); deleteErr != nil {
			impl.Logger.Error("database delete by id error", slog.Any("error", deleteErr))
			return deleteErr
		}
	}

	//
	// STEP 2
	// List through all the NFT assets and delete them from IPFS network and
	// from the database.
	//

	assetsRes, listErr := impl.NFTAssetStorer.ListByNFTCollectionID(ctx, id)
	if listErr != nil {
		impl.Logger.Error("database list by nft asset error", slog.Any("error", listErr))
		return listErr
	}
	for _, asset := range assetsRes.Results {
		if err := impl.IPFS.Unpin(ctx, asset.CID); err != nil {
			impl.Logger.Error("ipfs failed unpinning file error", slog.Any("error", err))
			return err
		}
		if deleteErr := impl.NFTAssetStorer.DeleteByID(ctx, asset.ID); deleteErr != nil {
			impl.Logger.Error("database delete by id error", slog.Any("error", deleteErr))
			return deleteErr
		}
	}

	//
	// STEP 3
	// Delete the collection folder from IPFS network.
	//

	impl.Logger.Debug("deleted all nft metadata files from ipfs for this collection")
	if err := impl.IPFS.Unpin(ctx, d.IPFSDirectoryCID); err != nil {
		impl.Logger.Error("ipfs failed unpinning directory error", slog.Any("error", err))
		return err
	}

	//
	// STEP 4
	// Deleted the database record.
	//

	if err := impl.NFTCollectionStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
