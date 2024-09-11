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
	nftcollection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

// NFTController Interface for tenant business logic controller.
type NFTController interface {
	Create(ctx context.Context, requestData *NFTCreateRequestIDO) (*nft_s.NFT, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*nft_s.NFT, error)
	UpdateByID(ctx context.Context, m *NFTUpdateRequestIDO) (*nft_s.NFT, error)
	ListByFilter(ctx context.Context, f *nft_s.NFTPaginationListFilter) (*nft_s.NFTPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *nft_s.NFTPaginationListFilter) ([]*nft_s.NFTAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ArchiveByID(ctx context.Context, id primitive.ObjectID) error
}

type NFTControllerImpl struct {
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
	NFTStorer   nft_s.NFTStorer
	NFTCollectionStorer nftcollection_s.NFTCollectionStorer
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
	nftcollection_storer nftcollection_s.NFTCollectionStorer,
	usr_storer user_s.UserStorer,
) NFTController {
	s := &NFTControllerImpl{
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
		NFTStorer:   nft_storer,
		NFTCollectionStorer: nftcollection_storer,
		UserStorer:          usr_storer,
	}
	s.Logger.Debug("nft controller initialization started...")
	s.Logger.Debug("nft controller initialized")
	return s
}
