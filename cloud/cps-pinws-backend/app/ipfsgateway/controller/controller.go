package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/ipfs"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/s3"
	pin_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	project_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/uuid"
)

type IpfsGatewayController interface {
	GetByContentID(ctx context.Context, cid string) (*pin_s.PinObject, error)
	Shutdown()
}

type IpfsGatewayControllerImpl struct {
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
) IpfsGatewayController {
	s := &IpfsGatewayControllerImpl{
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
	s.Logger.Debug("pinobject controller initialized")
	return s
}

func (impl *IpfsGatewayControllerImpl) Shutdown() {
	impl.Logger.Debug("shutting down...")
	impl.IPFS.Shutdown()
}
