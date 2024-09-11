//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/cache/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/emailer/mailgun"
	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/mongodb"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/templatedemailer"
	gateway_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/controller"
	gateway_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/httptransport"
	ipfsgate_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/controller"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	ipfsgate_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/httptransport"
	nftasset_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/controller"
	nftasset_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	nftasset_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/httptransport"
	nftcollection_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/controller"
	nftcollection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	nftcollection_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/httptransport"
	nftmetadata_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/controller"
	nftmetadata_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	nftmetadata_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/httptransport"
	tenant_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/controller"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	tenant_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/httptransport"
	user_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/controller"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	user_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/inputport/eventscheduler"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/inputport/http"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/inputport/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/blacklist"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/logger"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/time"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
)

func InitializeEvent() Application {
	// Our application is dependent on the following Golang packages. We need to
	// provide them to Google wire so it can sort out the dependency injection
	// at compile time.
	wire.Build(
		config.New,
		uuid.NewProvider,
		time.NewProvider,
		logger.NewProvider,
		jwt.NewProvider,
		kmutex.NewProvider,
		password.NewProvider,
		mailgun.NewEmailer,
		templatedemailer.NewTemplatedEmailer,
		mongodb.NewStorage,
		blacklist.NewProvider,
		mongodbcache.NewCache,
		ipfs_storage.NewStorage,
		pinobject_s.NewDatastore,
		ipfsgate_c.NewController,
		ipfsgate_http.NewHandler,
		user_s.NewDatastore,
		user_c.NewController,
		nftmetadata_s.NewDatastore,
		nftmetadata_c.NewController,
		nftasset_s.NewDatastore,
		nftasset_c.NewController,
		nftcollection_s.NewDatastore,
		nftcollection_c.NewController,
		tenant_s.NewDatastore,
		tenant_c.NewController,
		gateway_c.NewController,
		gateway_http.NewHandler,
		user_http.NewHandler,
		nftasset_http.NewHandler,
		nftmetadata_http.NewHandler,
		nftcollection_http.NewHandler,
		tenant_http.NewHandler,
		middleware.NewMiddleware,
		http.NewInputPort,
		eventscheduler.NewInputPort,
		NewApplication)
	return Application{}
}
