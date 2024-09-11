package eventscheduler

import (
	"log/slog"

	"github.com/mileusna/crontab"

	nftasset "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/controller"
	collection "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/controller"
	nftmetadata "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/controller"
	tenant "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/controller"
	user "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
)

type InputPortServer interface {
	Run()
	Shutdown()
}

type crontabInputPort struct {
	Config        *config.Conf
	Logger        *slog.Logger
	Crontab       *crontab.Crontab
	User          user.UserController
	Tenant        tenant.TenantController
	NFTCollection collection.NFTCollectionController
	NFTMetadata   nftmetadata.NFTMetadataController
	NFTAsset      nftasset.NFTAssetController
}

func NewInputPort(
	configp *config.Conf,
	loggerp *slog.Logger,
	cu user.UserController,
	org tenant.TenantController,
	co collection.NFTCollectionController,
	nftmetadata nftmetadata.NFTMetadataController,
	nftasset nftasset.NFTAssetController,
) InputPortServer {

	ctab := crontab.New() // create cron table

	// Create our HTTP server controller.
	p := &crontabInputPort{
		Config:        configp,
		Logger:        loggerp,
		Crontab:       ctab,
		User:          cu,
		Tenant:        org,
		NFTCollection: co,
		NFTMetadata:   nftmetadata,
		NFTAsset:      nftasset,
	}

	return p
}

func (port *crontabInputPort) Run() {
	port.Logger.Info("event scheduler running")
	// port.Crontab.MustAddJob("* * * * *", port.runGarbageCollection) // every minute
	port.Crontab.MustAddJob("0 0 * * *", port.runReprovideIPNS)       // every 24 hours
	port.Crontab.MustAddJob("00 22 * * *", port.runGarbageCollection) // every day at 10 pm

}

func (port *crontabInputPort) Shutdown() {
	port.Logger.Info("event scheduler shutting down now...")
	port.Logger.Info("event scheduler shutdown")
}
