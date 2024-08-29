package controller

import (
	"context"
	"log"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/adapter/storage/ipfs"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/adapter/storage/s3"
	pin_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	project_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/project/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/uuid"
)

// PinObjectController Interface for pinobject business logic controller.
type PinObjectController interface {
	IpfsAdd(ctx context.Context, req *IpfsAddRequestIDO) (string, error)
	Create(ctx context.Context, req *PinObjectCreateRequestIDO) (*pin_s.PinObject, error)
	GetByRequestID(ctx context.Context, requestID primitive.ObjectID) (*pin_s.PinObject, error)
	GetWithContentByRequestID(ctx context.Context, requestID primitive.ObjectID) (*pin_s.PinObject, error)
	UpdateByRequestID(ctx context.Context, ns *PinObjectUpdateRequestIDO) (*pin_s.PinObject, error)
	ListByFilter(ctx context.Context, f *pin_s.PinObjectPaginationListFilter) (*pin_s.PinObjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *pin_s.PinObjectPaginationListFilter) ([]*pin_s.PinObjectAsSelectOption, error)
	DeleteByRequestID(ctx context.Context, requestID primitive.ObjectID) error
	Shutdown()
}

type PinObjectControllerImpl struct {
	Config          *config.Conf
	Logger          *slog.Logger
	UUID            uuid.Provider
	Password        password.Provider
	JWT             jwt.Provider
	IPFS            ipfs_storage.IPFSStorager
	S3              s3_storage.S3Storager
	DbClient        *mongo.Client
	ProjectStorer   project_s.ProjectStorer
	PinObjectStorer pinobject_s.PinObjectStorer
	UserStorer      user_s.UserStorer
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
	proj_storer project_s.ProjectStorer,
	pin_storer pinobject_s.PinObjectStorer,
	usr_storer user_s.UserStorer,
) PinObjectController {
	s := &PinObjectControllerImpl{
		Config:          appCfg,
		Logger:          loggerp,
		UUID:            uuidp,
		Password:        passwordp,
		JWT:             jwtp,
		S3:              s3,
		IPFS:            ipfs,
		DbClient:        client,
		ProjectStorer:   proj_storer,
		PinObjectStorer: pin_storer,
		UserStorer:      usr_storer,
	}
	s.Logger.Debug("pinobject controller initialization started...")
	if err := s.s3SyncWithIpfs(context.Background()); err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}
	s.Logger.Debug("pinobject controller initialized")
	return s
}

func (impl *PinObjectControllerImpl) Shutdown() {
	impl.Logger.Debug("shutting down...")
	impl.IPFS.Shutdown()
}
