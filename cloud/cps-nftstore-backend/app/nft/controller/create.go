package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTControllerImpl) Create(ctx context.Context, nft *s_d.NFT) (*s_d.NFT, error) {
	// Extract user and tenant information from the session context
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName := ctx.Value(constants.SessionUserName).(string)
	tenantID, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	tenantName, _ := ctx.Value(constants.SessionUserTenantName).(string)
	tenantTimezone, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)

	// Check if the user has the necessary permissions
	switch userRole {
	case u_d.UserRoleRoot:
		// Access is granted; proceed with the operation
	default:
		// Deny access if the user does not have the 'Root' role
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Start a MongoDB session for transaction management
	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("failed to start database session", slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define the transaction function to perform a series of operations atomically
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Populate nft with default and tenant-specific information
		nft.TenantID = tenantID
		nft.TenantName = tenantName
		nft.TenantTimezone = tenantTimezone
		nft.ID = primitive.NewObjectID()
		nft.CreatedAt = time.Now()
		nft.ModifiedAt = time.Now()
		nft.Status = s_d.StatusActive

		// Save the nft data to the database
		if err := impl.NFTStorer.Create(sessCtx, nft); err != nil {
			impl.Logger.Error("failed to save nft to database", slog.Any("error", err))
			return nil, err
		}

		return nft, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*s_d.NFT), nil
}
