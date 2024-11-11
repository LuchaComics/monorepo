package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

type IpfsGatewayController interface {
	GetByCID(ctx context.Context, cid string) (*pinobject_s.PinObject, error)
	GetByIPNSPath(ctx context.Context, ipnsPath string) (*pinobject_s.PinObject, error)
	Shutdown()
}

type IpfsGatewayControllerImpl struct {
	Config          *config.Conf
	Logger          *slog.Logger
	UUID            uuid.Provider
	Password        password.Provider
	JWT             jwt.Provider
	IPFS            ipfs_storage.IPFSStorager
	DbClient        *mongo.Client
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
	client *mongo.Client,
	pinobject_storer pinobject_s.PinObjectStorer,
) IpfsGatewayController {
	s := &IpfsGatewayControllerImpl{
		Config:          appCfg,
		Logger:          loggerp,
		UUID:            uuidp,
		Password:        passwordp,
		JWT:             jwtp,
		IPFS:            ipfs,
		DbClient:        client,
		PinObjectStorer: pinobject_storer,
	}
	s.Logger.Debug("pinobject controller initialization started...")
	s.Logger.Debug("pinobject controller initialized")
	return s
}

func (impl *IpfsGatewayControllerImpl) Shutdown() {
	impl.Logger.Debug("shutting down...")
	impl.IPFS.Shutdown()
}
