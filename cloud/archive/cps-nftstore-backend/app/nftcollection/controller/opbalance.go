package controller

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"go.mongodb.org/mongo-driver/bson/primitive"

	eth "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/blockchain/eth"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type GetWalletBalanceOperationResponseIDO struct {
	Value *big.Float `bson:"value" json:"value"`
}

func (impl *NFTCollectionControllerImpl) OperationGetWalletBalanceByID(ctx context.Context, id primitive.ObjectID) (*GetWalletBalanceOperationResponseIDO, error) {
	// Retrieve from our database the record for the specific id.
	m, err := impl.NFTCollectionStorer.GetByID(ctx, id)
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
	balance, balanceErr := eth.GetBalance(m.WalletAccountAddress)
	if balanceErr != nil {
		impl.Logger.Error("failed getting balance", slog.Any("error", balanceErr))
		return nil, httperror.NewForBadRequestWithSingleField("node_url", fmt.Sprintf("balance error: %v", balanceErr))
	}

	return &GetWalletBalanceOperationResponseIDO{
		Value: balance,
	}, err
}
