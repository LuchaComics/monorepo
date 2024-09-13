package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	eth "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	collection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/cryptowrapper"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type DeploySmartContractOperationRequestIDO struct {
	CollectionID   primitive.ObjectID `bson:"collection_id" json:"collection_id"`
	WalletPassword string             `bson:"wallet_password" json:"wallet_password"`
}

type DeploySmartContractOperationResponseIDO struct {
	// Value *big.Float `bson:"value" json:"value"`
}

func validateDeploySmartContractOperationRequest(dirtyData *DeploySmartContractOperationRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.CollectionID.IsZero() {
		e["collection_id"] = "missing value"
	}
	if dirtyData.WalletPassword == "" {
		e["wallet_password"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTCollectionControllerImpl) OperationDeploySmartContract(ctx context.Context, req *DeploySmartContractOperationRequestIDO) (*collection_s.NFTCollection, error) {
	if valErr := validateDeploySmartContractOperationRequest(req); valErr != nil {
		return nil, valErr
	}

	// Extract user and tenant information from the session context
	userRole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// userName := ctx.Value(constants.SessionUserName).(string)
	// tenantID, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	// tenantName, _ := ctx.Value(constants.SessionUserTenantName).(string)
	// tenantTimezone, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)
	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	// Check if the user has the necessary permissions
	switch userRole {
	case u_d.UserRoleRoot:
		// Access is granted; proceed with the operation
	default:
		// Deny access if the user does not have the 'Root' role
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Start a MongoDB session for transaction management
	session, startSessErr := impl.DbClient.StartSession()
	if startSessErr != nil {
		impl.Logger.Error("failed to start database session", slog.Any("error", startSessErr))
		return nil, startSessErr
	}
	defer session.EndSession(ctx)

	// Define the transaction function to perform a series of operations atomically
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		//
		// STEP 1
		// Fetch all the related records from the database.
		//

		// Retrieve from our database the record for the specific id.
		collection, getErr := impl.NFTCollectionStorer.GetByID(ctx, req.CollectionID)
		if getErr != nil {
			impl.Logger.Error("database get by id error", slog.Any("error", getErr))
			return nil, getErr
		}
		if collection == nil {
			return nil, httperror.NewForBadRequestWithSingleField("collection_id", "collection does not exist")
		}

		//
		// STEP 2
		// Decrypt the wallet private key (which is saved in our database in
		// encrypted form) so we can use the plaintext private key for our
		// ethereum deploy smart contract operation.
		//

		plaintextPrivateKey, cryptoErr := cryptowrapper.SymmetricKeyDecryption(collection.WalletEncryptedPrivateKey, req.WalletPassword)
		if cryptoErr != nil {
			impl.Logger.Error("failed to decrypt wallet private key", slog.Any("error", cryptoErr))
			return nil, httperror.NewForBadRequestWithSingleField("wallet_password", "incorrect password used")
		}

		impl.Logger.Debug("decrypted ethereum wallet private key",
			slog.String("collection_id", collection.ID.Hex()))

		//
		// STEP 3
		// Connect to ethereum blockchain network via our node.
		//

		eth := eth.NewAdapter(impl.Logger)
		if connectErr := eth.ConnectToNodeAtURL(collection.NodeURL); connectErr != nil {
			impl.Logger.Error("failed connecting to node", slog.Any("error", connectErr))
			return nil, httperror.NewForBadRequestWithSingleField("node_url", fmt.Sprintf("connection error: %v", connectErr))
		}

		//
		// STEP 4
		// Deploy our smart contract to the ethereum blockchain.
		//

		// DEVELOPERS NOTE:
		// (1) We specify the private key of the ethereum wallet whom has money
		//     in their account, because this function will use that money to
		//     deploy the smart contract; hence, this account will lose money
		//     as result of this operation.
		// (2) We specify the collections wallet address because according to
		//     our contract when we deploy we need to specify an ethereum wallet
		//     which will be the `owner` that is able to mint NFTs. You as the
		//     programmer will have to review the smart contract to verify
		//     yourself.
		smartContractAddress, deployErr := eth.DeploySmartContract(collection.SmartContract, plaintextPrivateKey, collection.IPNSName)
		if deployErr != nil {
			impl.Logger.Error("failed deploying to ethereum blockchain",
				slog.Any("error", deployErr))
			return nil, httperror.NewForBadRequestWithSingleField("deployment_error", fmt.Sprintf("failed deploying: %v", deployErr))
		}

		impl.Logger.Debug("successfully deploy smart contract to blockchain",
			slog.String("collection_id", collection.ID.Hex()),
			slog.String("smart_contract_address", smartContractAddress))

		//
		// STEP 5
		// Update our database record.
		//

		collection.SmartContractAddress = smartContractAddress
		collection.SmartContractStatus = collection_s.SmartContractStatusDeployed
		collection.ModifiedAt = time.Now()
		collection.ModifiedFromIPAddress = ipAddress

		if updateErr := impl.NFTCollectionStorer.UpdateByID(ctx, collection); updateErr != nil {
			impl.Logger.Error("database update by id error", slog.Any("error", updateErr))
			return nil, updateErr
		}

		impl.Logger.Debug("collection updated in database",
			slog.String("collection_id", collection.ID.Hex()))

		return collection, nil
	}

	// Execute the transaction function within a MongoDB session
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("transaction failed", slog.Any("error", err))
		return nil, err
	}

	return result.(*s_d.NFTCollection), nil
}
