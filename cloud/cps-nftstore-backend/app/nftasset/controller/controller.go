package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/s3"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	collection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	nftmetadata_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

// NFTAssetController Interface for nftasset business logic controller.
type NFTAssetController interface {
	// IpfsAdd(ctx context.Context, req *IpfsAddRequestIDO) (string, error)
	Create(ctx context.Context, req *NFTAssetCreateRequestIDO) (*nftasset_s.NFTAsset, error)
	// GetByRequestID(ctx context.Context, requestID primitive.ObjectID) (*pin_s.NFTAsset, error)
	// GetWithContentByRequestID(ctx context.Context, requestID primitive.ObjectID) (*pin_s.NFTAsset, error)
	// UpdateByRequestID(ctx context.Context, ns *NFTAssetUpdateRequestIDO) (*pin_s.NFTAsset, error)
	// ListByFilter(ctx context.Context, f *pin_s.NFTAssetPaginationListFilter) (*pin_s.NFTAssetPaginationListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *pin_s.NFTAssetPaginationListFilter) ([]*pin_s.NFTAssetAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	Shutdown()
}

type NFTAssetControllerImpl struct {
	Config              *config.Conf
	Logger              *slog.Logger
	UUID                uuid.Provider
	Password            password.Provider
	JWT                 jwt.Provider
	IPFS                ipfs_storage.IPFSStorager
	S3                  s3_storage.S3Storager
	DbClient            *mongo.Client
	NFTAssetStorer      nftasset_s.NFTAssetStorer
	NFTMetadataStorer   nftmetadata_s.NFTMetadataStorer
	NFTCollectionStorer collection_s.NFTCollectionStorer
	UserStorer          user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	jwtp jwt.Provider,
	ipfs ipfs_storage.IPFSStorager,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	nftasset_storer nftasset_s.NFTAssetStorer,
	nftmetadata_storer nftmetadata_s.NFTMetadataStorer,
	collection_storer collection_s.NFTCollectionStorer,
	usr_storer user_s.UserStorer,
) NFTAssetController {
	s := &NFTAssetControllerImpl{
		Config:              appCfg,
		Logger:              loggerp,
		UUID:                uuidp,
		Password:            passwordp,
		JWT:                 jwtp,
		S3:                  s3,
		IPFS:                ipfs,
		DbClient:            client,
		NFTAssetStorer:      nftasset_storer,
		NFTMetadataStorer:   nftmetadata_storer,
		NFTCollectionStorer: collection_storer,
		UserStorer:          usr_storer,
	}
	s.Logger.Debug("nftasset controller initialization started...")
	s.Logger.Debug("nftasset controller initialized")
	return s
}

func (impl *NFTAssetControllerImpl) Shutdown() {
	impl.Logger.Debug("shutting down...")
	impl.IPFS.Shutdown()
}
