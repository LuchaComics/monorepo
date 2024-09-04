package controller

import (
	"context"
	"log/slog"

	"github.com/bartmika/arraydiff"

	"go.mongodb.org/mongo-driver/mongo"
)

func (impl *PinObjectControllerImpl) s3SyncWithIpfs(ctx context.Context) error {
	impl.Logger.Debug("synching s3 with ipfs...")
	defer impl.Logger.Debug("synched s3 with ipfs")

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		dbCids, err := impl.PinObjectStorer.GetAllCIDs(ctx)
		if err != nil {
			return nil, err
		}
		ipfsCids, err := impl.IPFS.ListPins(ctx)
		if err != nil {
			return nil, err
		}

		// See what are the differences between the two arrays of type `uint64` data-types.
		addCids, comCids, rmCids := arraydiff.Strings(ipfsCids, dbCids)

		// For debugging purposes only.
		impl.Logger.Debug("fetched cids",
			slog.Any("insert_cids", addCids),
			slog.Any("common_cids", comCids),
			slog.Any("remove_cids", rmCids),
			slog.Any("db_cids", dbCids),
			slog.Any("ipfs_cids", ipfsCids))

		for _, addCid := range addCids {
			impl.Logger.Debug("save to local ipfs from s3",
				slog.Any("add_cid", addCid))

			pin, err := impl.PinObjectStorer.GetByCID(sessCtx, addCid)
			if err != nil {
				impl.Logger.Warn("error looking up pin by cid",
					slog.Any("cid", addCid),
					slog.Any("error", err))
				continue
			}

			if pin != nil {
				fileContent, err := impl.S3.GetContentByKey(sessCtx, pin.ObjectKey)
				if err != nil {
					impl.Logger.Warn("error getting content from s3",
						slog.Any("cid", addCid),
						slog.Any("object_key", pin.ObjectKey),
						slog.Any("error", err))
					continue
				}

				// // Upload to IPFS network and pin as well.
				if _, err := impl.IPFS.AddFileContentAndPin(ctx, fileContent); err != nil {
					impl.Logger.Error("failed uploading and pinning to IPFS", slog.Any("error", err))
					continue
				}

				impl.Logger.Warn("content saved locally from remote via s3",
					slog.Any("cid", addCid))
			}
		}

		for _, rmCid := range rmCids {
			impl.Logger.Debug("removing from local ipfs",
				slog.Any("remove_cid", rmCid))

			if err := impl.IPFS.DeleteContent(ctx, rmCid); err != nil {
				impl.Logger.Error("failed deleting content from IPFS",
					slog.Any("cid", rmCid),
					slog.Any("error", err))
				continue
			}

			impl.Logger.Debug("removed from local ipfs",
				slog.Any("remove_cid", rmCid))
		}

		// for _, cid := range cidsInDB {
		// 	// Check to see if our ipfs has the `cid` locally and if not
		// 	// then download from our s3 and save it to our local ipfs node.
		// 	// Else if we already have it then do nothing as we are synched.
		// 	pinnedContent, getContentErr := impl.IPFS.GetContent(sessCtx, cid)
		// 	if getContentErr != nil {
		// 		impl.Logger.Error("get content from ipfs by cid error",
		// 			slog.String("cid", cid),
		// 			slog.Any("error", err))
		// 		return nil, err
		// 	}
		//
		// 	if pinnedContent != nil {
		// 		impl.Logger.Debug("content found locally, skipping sync",
		// 			slog.Any("cid", cid))
		// 	} else {
		// 		impl.Logger.Warn("content not found locally, fetching from remote via s3",
		// 			slog.Any("cid", cid))
		//
		// 		pin, err := impl.PinObjectStorer.GetByCID(sessCtx, cid)
		// 		if err != nil {
		// 			impl.Logger.Warn("error looking up pin by cid",
		// 				slog.Any("cid", cid),
		// 				slog.Any("error", err))
		// 			continue
		// 		}
		// 		if pin != nil {
		// 			fileContent, err := impl.S3.GetContentByKey(sessCtx, pin.ObjectKey)
		// 			if err != nil {
		// 				impl.Logger.Warn("error getting content from s3",
		// 					slog.Any("cid", cid),
		// 					slog.Any("object_key", pin.ObjectKey),
		// 					slog.Any("error", err))
		// 				continue
		// 			}
		//
		// 			// // Upload to IPFS network.
		// 			if _, err := impl.IPFS.UploadContentAndPin(ctx, pin.Meta["filename"], fileContent); err != nil {
		// 				impl.Logger.Error("failed uploading and pinning to IPFS", slog.Any("error", err))
		// 				continue
		// 			}
		//
		// 			impl.Logger.Warn("content saved locally from remote via s3",
		// 				slog.Any("cid", cid))
		// 		}
		// 	}
		// }

		return nil, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}
