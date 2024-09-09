package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *CollectionControllerImpl) Create(ctx context.Context, m *s_d.Collection) (*s_d.Collection, error) {
	// Extract from our session the following data.
	urole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	oname, _ := ctx.Value(constants.SessionUserTenantName).(string)
	otz, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		// impl.Logger.Debug("access granted")
		// Do nothing
	default:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
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
		// Add defaults.
		m.TenantID = oid
		m.TenantName = oname
		m.TenantTimezone = otz
		m.ID = primitive.NewObjectID()
		m.CreatedAt = time.Now()
		m.ModifiedAt = time.Now()
		m.Status = s_d.StatusActive

		// DEVELOPERS NOTE:
		// Every collection has a unique IPNS record, as a result we will need
		// to generate a new `key` for it. In addition the IPFS RPC returns
		// the IPNS name.

		keyName := fmt.Sprintf("ipns_key_%s", m.ID.Hex())
		ipnsName, err := impl.IPFS.GenerateKey(sessCtx, keyName)
		if err != nil {
			impl.Logger.Error("failed to generate key error",
				slog.Any("error", err))
			return nil, err
		}

		// Give our collection's folder a custom name.
		m.IpfsDirectoryName = fmt.Sprintf("%v_metadata", m.ID.Hex())

		// Save the IPNS record related data.
		m.IpnsKeyName = keyName
		m.IpnsName = ipnsName

		// Create our NFT collections folder and create a sample file named `0`
		// because our `token_id` increments by one.
		collectionDirCid, firstTokenFileCid, ipfsApiAddErr := impl.IPFS.UploadContentFromStringWithFolder(
			context.Background(),
			"Hello world via `Collectibles Protective Services`!", "0", // Create a sample file...
			m.IpfsDirectoryName)
		if ipfsApiAddErr != nil {
			return nil, fmt.Errorf("ipfs failed adding to api: %v\n", ipfsApiAddErr)
		}
		impl.Logger.Debug("ipfs storage successfully uploaded",
			slog.String("collection_cid", collectionDirCid),
			slog.String("first_token_cid", firstTokenFileCid))

		// Publish our NFT collection folder to IPNS.
		resIpnsName, ipfsPublishErr := impl.IPFS.PublishToIPNS(sessCtx, keyName, collectionDirCid)
		if ipfsPublishErr != nil {
			return nil, fmt.Errorf("ipns failed publishing to api: %v\n", ipfsApiAddErr)
		}

		// For defensive code purposes only.
		if !strings.Contains(ipnsName, resIpnsName) {
			return nil, fmt.Errorf("ipns error: does not match: %s and %s", ipnsName, resIpnsName)
		}

		// Save to our database.
		if dbCreateErr := impl.CollectionStorer.Create(sessCtx, m); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", dbCreateErr))
			return nil, err
		}

		return m, nil
	}

	// Start a transaction

	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*s_d.Collection), nil
}
