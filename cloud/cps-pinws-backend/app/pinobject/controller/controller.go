package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/ipfs"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/s3"
	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/uuid"
)

// PinObjectController Interface for pinobject business logic controller.
type PinObjectController interface {
	Create(ctx context.Context, req *PinObjectCreateRequestIDO) (*domain.PinObject, error)
	GetByRequestID(ctx context.Context, requestID primitive.ObjectID) (*domain.PinObject, error)
	UpdateByRequestID(ctx context.Context, ns *PinObjectUpdateRequestIDO) (*domain.PinObject, error)
	ListByFilter(ctx context.Context, f *domain.PinObjectPaginationListFilter) (*domain.PinObjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.PinObjectPaginationListFilter) ([]*domain.PinObjectAsSelectOption, error)
	DeleteByRequestID(ctx context.Context, requestID primitive.ObjectID) error
	PermanentlyDeleteByRequestID(ctx context.Context, requestID primitive.ObjectID) error
	Shutdown()
}

type PinObjectControllerImpl struct {
	Config          *config.Conf
	Logger          *slog.Logger
	UUID            uuid.Provider
	IPFS            ipfs_storage.IPFSStorager
	S3              s3_storage.S3Storager
	DbClient        *mongo.Client
	PinObjectStorer pinobject_s.PinObjectStorer
	UserStorer      user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	ipfs ipfs_storage.IPFSStorager,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	org_storer pinobject_s.PinObjectStorer,
	usr_storer user_s.UserStorer,
) PinObjectController {
	s := &PinObjectControllerImpl{
		Config:          appCfg,
		Logger:          loggerp,
		UUID:            uuidp,
		S3:              s3,
		IPFS:            ipfs,
		DbClient:        client,
		PinObjectStorer: org_storer,
		UserStorer:      usr_storer,
	}
	s.Logger.Debug("pinobject controller initialization started...")
	s.Logger.Debug("pinobject controller initialized")
	return s
}

func (impl *PinObjectControllerImpl) Shutdown() {
	impl.Logger.Debug("shutting down...")
	impl.IPFS.Shutdown()
}
