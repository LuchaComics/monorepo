package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	nftcollection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	nftmetadata_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

// NFTMetadataController Interface for tenant business logic controller.
type NFTMetadataController interface {
	Create(ctx context.Context, requestData *NFTMetadataCreateRequestIDO) (*nftmetadata_s.NFTMetadata, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*nftmetadata_s.NFTMetadata, error)
	UpdateByID(ctx context.Context, m *nftmetadata_s.NFTMetadata) (*nftmetadata_s.NFTMetadata, error)
	ListByFilter(ctx context.Context, f *nftmetadata_s.NFTMetadataPaginationListFilter) (*nftmetadata_s.NFTMetadataPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *nftmetadata_s.NFTMetadataPaginationListFilter) ([]*nftmetadata_s.NFTMetadataAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type NFTMetadataControllerImpl struct {
	Config              *config.Conf
	Logger              *slog.Logger
	UUID                uuid.Provider
	JWT                 jwt.Provider
	Kmutex              kmutex.Provider
	Password            password.Provider
	IPFS                ipfs_storage.IPFSStorager
	DbClient            *mongo.Client
	TenantStorer        tenant_s.TenantStorer
	NFTAssetStorer      nftasset_s.NFTAssetStorer
	NFTMetadataStorer   nftmetadata_s.NFTMetadataStorer
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
	nftasset_storer nftasset_s.NFTAssetStorer,
	nftmetadata_storer nftmetadata_s.NFTMetadataStorer,
	nftcollection_storer nftcollection_s.NFTCollectionStorer,
	usr_storer user_s.UserStorer,
) NFTMetadataController {
	s := &NFTMetadataControllerImpl{
		Config:              appCfg,
		Logger:              loggerp,
		UUID:                uuidp,
		JWT:                 jwtp,
		Kmutex:              kmx,
		Password:            passwordp,
		IPFS:                ipfs,
		DbClient:            client,
		TenantStorer:        tenant_storer,
		NFTAssetStorer:      nftasset_storer,
		NFTMetadataStorer:   nftmetadata_storer,
		NFTCollectionStorer: nftcollection_storer,
		UserStorer:          usr_storer,
	}
	s.Logger.Debug("nftmetadata controller initialization started...")
	s.Logger.Debug("nftmetadata controller initialized")
	return s
}
