package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	nft_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	collection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

// NFTCollectionController Interface for tenant business logic controller.
type NFTCollectionController interface {
	Create(ctx context.Context, req *NFTCollectionCreateRequestIDO) (*collection_s.NFTCollection, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*collection_s.NFTCollection, error)
	UpdateByID(ctx context.Context, m *collection_s.NFTCollection) (*collection_s.NFTCollection, error)
	ListByFilter(ctx context.Context, f *collection_s.NFTCollectionPaginationListFilter) (*collection_s.NFTCollectionPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *collection_s.NFTCollectionPaginationListFilter) ([]*collection_s.NFTCollectionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ReprovidehCollectionsInIPNS(ctx context.Context) error
	OperationGetWalletBalanceByID(ctx context.Context, id primitive.ObjectID) (*GetWalletBalanceOperationResponseIDO, error)
	OperationDeploySmartContract(ctx context.Context, req *DeploySmartContractOperationRequestIDO) (*collection_s.NFTCollection, error)
}

type NFTCollectionControllerImpl struct {
	Config              *config.Conf
	Logger              *slog.Logger
	UUID                uuid.Provider
	JWT                 jwt.Provider
	Kmutex              kmutex.Provider
	Password            password.Provider
	IPFS                ipfs_storage.IPFSStorager
	DbClient            *mongo.Client
	TenantStorer        tenant_s.TenantStorer
	PinObjectStorer     pinobject_s.PinObjectStorer
	NFTAssetStorer      nftasset_s.NFTAssetStorer
	NFTStorer           nft_s.NFTStorer
	NFTCollectionStorer collection_s.NFTCollectionStorer
	UserStorer          user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	jwtp jwt.Provider,
	kmx kmutex.Provider,
	passwordp password.Provider,
	ipfs ipfs_storage.IPFSStorager,
	client *mongo.Client,
	tenant_storer tenant_s.TenantStorer,
	pinobject_storer pinobject_s.PinObjectStorer,
	nftasset_storer nftasset_s.NFTAssetStorer,
	nft_storer nft_s.NFTStorer,
	collection_storer collection_s.NFTCollectionStorer,
	usr_storer user_s.UserStorer,
) NFTCollectionController {
	s := &NFTCollectionControllerImpl{
		Config:              appCfg,
		Logger:              loggerp,
		UUID:                uuidp,
		JWT:                 jwtp,
		Kmutex:              kmx,
		Password:            passwordp,
		IPFS:                ipfs,
		DbClient:            client,
		TenantStorer:        tenant_storer,
		PinObjectStorer:     pinobject_storer,
		NFTAssetStorer:      nftasset_storer,
		NFTStorer:           nft_storer,
		NFTCollectionStorer: collection_storer,
		UserStorer:          usr_storer,
	}
	s.Logger.Debug("collection controller initialization started...")
	s.Logger.Debug("collection controller initialized")
	return s
}
