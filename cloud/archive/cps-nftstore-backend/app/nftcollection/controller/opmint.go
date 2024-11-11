package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	eth "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/cryptowrapper"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type MintOperationRequestIDO struct {
	CollectionID   primitive.ObjectID `bson:"collection_id" json:"collection_id"`
	TokenID        uint64             `bson:"token_id" json:"token_id"`
	ToAddress      string             `bson:"to_address" json:"to_address"`
	WalletPassword string             `bson:"wallet_password" json:"wallet_password"`
}

func validateMintOperationRequest(dirtyData *MintOperationRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.CollectionID.IsZero() {
		e["collection_id"] = "missing value"
	}
	if dirtyData.ToAddress == "" {
		e["to_address"] = "missing value"
	}
	if dirtyData.WalletPassword == "" {
		e["wallet_password"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *NFTCollectionControllerImpl) OperationMint(ctx context.Context, req *MintOperationRequestIDO) error {
	if valErr := validateMintOperationRequest(req); valErr != nil {
		return valErr
	}

	ipAddress, _ := ctx.Value(constants.SessionIPAddress).(string)

	//
	// STEP 1
	// Fetch all the related records from the database.
	//

	// Retrieve from our database the record for the specific id.
	collection, err := impl.NFTCollectionStorer.GetByID(ctx, req.CollectionID)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if collection == nil {
		return httperror.NewForBadRequestWithSingleField("id", "collection does not exist")
	}
	nft, err := impl.NFTStorer.GetByTokenID(ctx, req.TokenID)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if nft == nil {
		return httperror.NewForBadRequestWithSingleField("id", "nft does not exist")
	}

	// Validate the collection with token.
	if nft.TokenID != collection.TokensCount {
		return httperror.NewForBadRequestWithSingleField("token_id", "this nft token cannot be minted because it is not in the queque to be minted yet")
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
		return httperror.NewForBadRequestWithSingleField("wallet_password", "incorrect password used")
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
		return httperror.NewForBadRequestWithSingleField("node_url", fmt.Sprintf("connection error: %v", connectErr))
	}

	//
	// STEP 4
	// Execute the `mint` funciton to our smart contract in the ethereum
	// blockchain network.
	//

	if mintErr := eth.Mint(collection.SmartContract, plaintextPrivateKey, collection.SmartContractAddress, req.ToAddress); mintErr != nil {
		impl.Logger.Error("failed minting", slog.Any("error", mintErr))
		return mintErr
	}

	impl.Logger.Debug("successfully minted",
		slog.String("collection_id", collection.ID.Hex()),
		slog.String("smart_contract_address", collection.SmartContractAddress))

	//
	// STEP 5
	// Update our database record.
	//

	nft.MintedToAddress = req.ToAddress
	nft.ModifiedAt = time.Now()
	nft.ModifiedFromIPAddress = ipAddress
	if err := impl.NFTStorer.UpdateByID(ctx, nft); err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}

	impl.Logger.Debug("updated nft in database",
		slog.String("collection_id", collection.ID.Hex()),
		slog.Uint64("token_id", nft.TokenID))

	return nil
}
