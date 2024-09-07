package controller

import (
	"context"
	"log/slog"
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
		impl.Logger.Debug("access granted")
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
		// m.CreatedByUserID = uid
		// m.CreatedByUserName = uname
		m.ModifiedAt = time.Now()
		// m.ModifiedByUserID = uid
		// m.ModifiedByUserName = uname
		m.Status = s_d.StatusActive

		// Save to our database.
		if err := impl.CollectionStorer.Create(ctx, m); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		// DEVELOPERS NOTE:
		// Every collection has a unique IPNS record, as a result we will need
		// to do the following:
		// 1. Generate new `key`.
		// 2. Generate a folder for this collection.
		// 3. Generate a file called `0` and populate it with a random string.
		// 4. Public this collection's folder to IPNS
		// 5. Save the IPNS record.

		// Step 1. Generate new `key`.
		if err := impl.IPFS.GenerateKey(sessCtx, m.ID.Hex()); err != nil {
			return nil, err
		}

		// Step 2. Generate a folder for this collection.
		//TODO: Impl.

		// Step 3. Generate a file called `0` and populate it with a random string.
		//TODO: Impl.

		// 4. Public this collection's folder to IPNS
		//TODO: Impl.

		// 5. Save the IPNS record.
		//TODO: Impl.

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
