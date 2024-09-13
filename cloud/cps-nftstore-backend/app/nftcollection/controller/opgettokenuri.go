package controller

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	eth "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type GetTokenURIResponseIDO struct {
	URI string `bson:"uri" json:"uri"`
}

func (impl *NFTCollectionControllerImpl) OperationGetTokenURI(ctx context.Context, collectionID primitive.ObjectID, tokenID uint64) (*GetTokenURIResponseIDO, error) {
	// Retrieve from our database the record for the specific id.
	m, err := impl.NFTCollectionStorer.GetByID(ctx, collectionID)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "collection does not exist")
	}

	eth := eth.NewAdapter(impl.Logger)
	if connectErr := eth.ConnectToNodeAtURL(m.NodeURL); connectErr != nil {
		impl.Logger.Error("failed connecting to node", slog.Any("error", connectErr))
		return nil, httperror.NewForBadRequestWithSingleField("node_url", fmt.Sprintf("connection error: %v", connectErr))
	}
	tokenURI, getErr := eth.GetTokenURI(m.SmartContractAddress, tokenID)
	if getErr != nil {
		impl.Logger.Error("failed getting balance", slog.Any("error", getErr))
		return nil, getErr
	}

	return &GetTokenURIResponseIDO{
		URI: tokenURI,
	}, err
}
