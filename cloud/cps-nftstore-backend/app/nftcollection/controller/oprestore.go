package controller

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pinobj_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	a_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type NFTCollectionOperationRestoreRequestIDO struct {
	Name        string
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
}

// NFTCollectionOperationRestoreResponseIDO represents `PinStatus` spec via https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers.
type NFTCollectionOperationRestoreResponseIDO struct {
	RequestID         primitive.ObjectID `bson:"requestid" json:"requestid"`
	Status            string             `bson:"status" json:"status"`
	OperationRestored time.Time          `bson:"created,omitempty" json:"created,omitempty"`
	Delegates         []string           `bson:"delegates" json:"delegates"`
	Info              map[string]string  `bson:"info" json:"info"`
}

func ValidateOperationRestoreRequest(dirtyData *NFTCollectionOperationRestoreRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.Filename == "" {
		e["filename"] = "missing value"
	}
	if dirtyData.ContentType == "" {
		e["content_type"] = "missing value"
	}
	if dirtyData.Data == nil {
		e["data"] = "missing value"
	}
	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTCollectionControllerImpl) OperationRestore(ctx context.Context, reqRaw *NFTCollectionOperationRestoreRequestIDO) (*a_d.NFTCollection, error) {
	if err := ValidateOperationRestoreRequest(reqRaw); err != nil {
		return nil, err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		//
		// STEP 1
		// Attempt to unmarshal the backupfile.
		//

		reqPayload := &NFTCollectionBackupOperationResponseIDO{}

		if err := json.Unmarshal(reqRaw.Data, &reqPayload); err != nil {
			return nil, err
		}
		if reqPayload == nil {
			return nil, httperror.NewForBadRequestWithSingleField("format failure", "failed to unmarshal")
		}
		if reqPayload.NFTCollection == nil {
			return nil, httperror.NewForBadRequestWithSingleField("nft_collection", "missing data")
		}

		//
		// STEP 2
		// Lookup the collection and check to see if it already exists.
		// If the file does not exist then proceed else if it exists then
		// we abort with an error.
		//

		nftCollectionDoesExist, nftColExistsErr := impl.NFTCollectionStorer.CheckIfExistsByID(sessCtx, reqPayload.NFTCollection.ID)
		if nftColExistsErr != nil {
			return nil, nftColExistsErr
		}
		if nftCollectionDoesExist {
			return nil, httperror.NewForBadRequestWithSingleField("message", "collection already exists, restoration cancelled")
		}

		//
		// STEP 3
		// Verify that the NFT Collection has our IPNS key saved in the IPFS
		// node that we are connected to, if not then error.
		//

		// keyNameExists, keyExistsErr := impl.IPFS.CheckIfKeyNameExists(sessCtx, reqPayload.NFTCollection.IPNSKeyName)
		// if keyExistsErr != nil {
		// 	return nil, keyExistsErr
		// }
		// if !keyNameExists {
		// 	return nil, httperror.NewForBadRequestWithSingleField("ipns_key_name", "missing key in ipfs node")
		// }

		//
		// STEP 4
		// Save NFT collection to the database.
		//

		// Save the collection data to the database
		if createCollErr := impl.NFTCollectionStorer.Create(sessCtx, reqPayload.NFTCollection); createCollErr != nil {
			impl.Logger.Error("failed to save collection to database", slog.Any("error", createCollErr))
			return nil, createCollErr
		}

		impl.Logger.Debug("finished creating nft collection",
			slog.String("collection_id", reqPayload.NFTCollection.ID.Hex()))

		//
		// STEP 5
		// Iterate through all the NFT assets and then (1) save to the database
		// and (2) pin to IPFS network.
		//

		for _, asset := range reqPayload.NFTAssets {
			assetsDoesExist, aExErr := impl.NFTAssetStorer.CheckIfExistsByID(sessCtx, asset.ID)
			if aExErr != nil {
				impl.Logger.Error("failed to check by id in nft assets to database", slog.Any("error", aExErr))
				return nil, aExErr
			}
			if assetsDoesExist {
				impl.Logger.Error("nft asset already exists error")
				return nil, httperror.NewForBadRequestWithSingleField("message", "nft asset already exists, restoration cancelled")
			}

			if err := impl.NFTAssetStorer.Create(sessCtx, asset); err != nil {
				impl.Logger.Error("database nft asset create error", slog.Any("error", err))
				return nil, err
			}

			impl.Logger.Debug("successfully created nft asset file",
				slog.String("id", asset.ID.Hex()))

			pinObject := &pinobj_s.PinObject{
				ID:          primitive.NewObjectID(),
				IPNSPath:    "", // Set to empty b/c this pin is not mounted to IPNS.
				CID:         asset.CID,
				Content:     nil,
				Filename:    asset.Filename,
				ContentType: asset.ContentType,
				CreatedAt:   asset.CreatedAt,
				ModifiedAt:  asset.ModifiedAt,
			}
			if createdErr := impl.PinObjectStorer.Create(sessCtx, pinObject); createdErr != nil {
				impl.Logger.Error("database create pin object for nft asset error", slog.Any("error", createdErr))
				return nil, err
			}

			impl.Logger.Debug("nft asset file pinned",
				slog.String("cid", pinObject.CID))
		}

		//
		// STEP 6
		// Iterate through all the NFTs and then (1) save to the database
		// and (2) pin to IPFS network.

		for _, nft := range reqPayload.NFTs {
			nftDoesExist, nftExErr := impl.NFTStorer.CheckIfExistsByID(sessCtx, nft.ID)
			if nftExErr != nil {
				impl.Logger.Error("failed to check by id in nfts to database", slog.Any("error", nftExErr))
				return nil, nftExErr
			}
			if nftDoesExist {
				impl.Logger.Error("nft already exists error")
				return nil, httperror.NewForBadRequestWithSingleField("message", "nft already exists, restoration cancelled")
			}

			if err := impl.NFTStorer.Create(sessCtx, nft); err != nil {
				impl.Logger.Error("database create nft error", slog.Any("error", err))
				return nil, err
			}

			if err := impl.NFTStorer.Create(sessCtx, nft); err != nil {
				impl.Logger.Error("failed to create nft to database",
					slog.Any("error", err))
				return nil, err
			}

			impl.Logger.Debug("nft created",
				slog.Uint64("token_id", nft.TokenID))

			pinObject := &pinobj_s.PinObject{
				ID:          primitive.NewObjectID(),
				IPNSPath:    nft.FileIPNSPath,
				CID:         nft.FileCID,
				Content:     nil,
				Filename:    fmt.Sprintf("%v", nft.TokenID), // We set it to this way b/c it is required by our `Smart Contract` to write the names like this - This is not an error!
				ContentType: "application/json",
				CreatedAt:   nft.CreatedAt,
				ModifiedAt:  nft.ModifiedAt,
			}
			if createdErr := impl.PinObjectStorer.Create(sessCtx, pinObject); createdErr != nil {
				impl.Logger.Error("database create pin object for nft error", slog.Any("error", createdErr))
				return nil, err
			}

			impl.Logger.Debug("nft metedata file pinned",
				slog.String("pinObject", pinObject.CID))
		}

		// return res, nil
		return reqPayload.NFTCollection, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*a_d.NFTCollection), nil
}
