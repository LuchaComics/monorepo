package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	collection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

// CollectionController Interface for tenant business logic controller.
type CollectionController interface {
	Create(ctx context.Context, m *collection_s.Collection) (*collection_s.Collection, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*collection_s.Collection, error)
	UpdateByID(ctx context.Context, m *collection_s.Collection) (*collection_s.Collection, error)
	ListByFilter(ctx context.Context, f *collection_s.CollectionPaginationListFilter) (*collection_s.CollectionPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *collection_s.CollectionPaginationListFilter) ([]*collection_s.CollectionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type CollectionControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	JWT              jwt.Provider
	Kmutex           kmutex.Provider
	Password         password.Provider
	IPFS             ipfs_storage.IPFSStorager
	DbClient         *mongo.Client
	TenantStorer     tenant_s.TenantStorer
	CollectionStorer collection_s.CollectionStorer
	UserStorer       user_s.UserStorer
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
	collection_storer collection_s.CollectionStorer,
	usr_storer user_s.UserStorer,
) CollectionController {
	s := &CollectionControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		JWT:              jwtp,
		Kmutex:           kmx,
		Password:         passwordp,
		IPFS:             ipfs,
		DbClient:         client,
		TenantStorer:     tenant_storer,
		CollectionStorer: collection_storer,
		UserStorer:       usr_storer,
	}
	s.Logger.Debug("collection controller initialization started...")
	s.Logger.Debug("collection controller initialized")
	return s
}
