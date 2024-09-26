package controller

import (
	"context"
	"log"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	pubsub "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/pubsub"
	a_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/uuid"
)

// BlockchainController provides all the functionality that can be performed
// on the `ComicCoin` cryptocurrency.

type BlockchainController interface {
	NewGenesisBlock(ctx context.Context, coinbaseKey *keystore.Key) (*block_ds.Block, error)
	GetBlock(ctx context.Context, hash string) (*block_ds.Block, error)
	GetBalanceByAccountName(ctx context.Context, accountName string) (*BlockchainBalanceResponseIDO, error)
	Submit(ctx context.Context, req *BlockchainSubmitRequestIDO) (*BlockchainSubmitResponseIDO, error)
	GetSignedTransactions(ctx context.Context) ([]*pt_ds.SignedTransaction, error)
}

type blockchainControllerImpl struct {
	config                  *config.Config
	logger                  *slog.Logger
	uuid                    uuid.Provider
	localPubSubBroker       pubsub.PubSubBroker
	p2pPubSubBroker         pubsub.PubSubBroker
	accountStorer           a_ds.AccountStorer
	signedTransactionStorer pt_ds.SignedTransactionStorer
	lastHashStorer          lasthash_ds.LastHashStorer
	blockStorer             block_ds.BlockStorer
}

func NewController(
	cfg *config.Config,
	logger *slog.Logger,
	uuid uuid.Provider,
	locpsbroker pubsub.PubSubBroker,
	p2psbroker pubsub.PubSubBroker,
	as a_ds.AccountStorer,
	pt pt_ds.SignedTransactionStorer,
	lhDS lasthash_ds.LastHashStorer,
	blockDS block_ds.BlockStorer,
) BlockchainController {
	// Defensive code to protect the programmer from any errors.
	if cfg.Blockchain.Difficulty <= 0 {
		log.Fatal("cannot have blochain difficulty less then or equal to zero")
	}
	if cfg.Blockchain.TransPerBlock <= 0 {
		log.Fatal("cannot have blochain transactions per block less then or equal to zero")
	}

	impl := &blockchainControllerImpl{
		config:                  cfg,
		logger:                  logger,
		uuid:                    uuid,
		localPubSubBroker:       locpsbroker,
		p2pPubSubBroker:         p2psbroker,
		accountStorer:           as,
		signedTransactionStorer: pt,
		lastHashStorer:          lhDS,
		blockStorer:             blockDS,
	}

	//TODO: Only if miner activated
	impl.runMinerOperationInBackground(context.Background())

	return impl
}
