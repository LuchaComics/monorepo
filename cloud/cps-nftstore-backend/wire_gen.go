// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/cache/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/ipfs"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/mongodb"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/templatedemailer"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/httptransport"
	controller7 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/controller"
	datastore3 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/datastore"
	httptransport7 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/ipfsgateway/httptransport"
	controller6 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/controller"
	datastore4 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	httptransport6 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/httptransport"
	controller4 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/controller"
	datastore6 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/datastore"
	httptransport4 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/httptransport"
	controller5 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/controller"
	datastore5 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/datastore"
	httptransport5 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftmetadata/httptransport"
	controller3 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/controller"
	datastore2 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	httptransport3 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/httptransport"
	controller2 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	httptransport2 "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/httptransport"
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

import (
	_ "go.uber.org/automaxprocs"
	_ "time/tzdata"
)

// Injectors from wire.go:

func InitializeEvent() Application {
	slogLogger := logger.NewProvider()
	conf := config.New()
	provider := uuid.NewProvider()
	timeProvider := time.NewProvider()
	jwtProvider := jwt.NewProvider(conf)
	blacklistProvider := blacklist.NewProvider()
	kmutexProvider := kmutex.NewProvider()
	passwordProvider := password.NewProvider()
	emailer := mailgun.NewEmailer(conf, slogLogger, provider)
	templatedEmailer := templatedemailer.NewTemplatedEmailer(conf, slogLogger, provider, emailer)
	client := mongodb.NewStorage(conf, slogLogger)
	cacher := mongodbcache.NewCache(conf, slogLogger, client)
	userStorer := datastore.NewDatastore(conf, slogLogger, client)
	tenantStorer := datastore2.NewDatastore(conf, slogLogger, client)
	gatewayController := controller.NewController(conf, slogLogger, provider, jwtProvider, kmutexProvider, passwordProvider, templatedEmailer, cacher, client, userStorer, tenantStorer)
	middlewareMiddleware := middleware.NewMiddleware(conf, slogLogger, provider, timeProvider, jwtProvider, blacklistProvider, gatewayController)
	handler := httptransport.NewHandler(slogLogger, gatewayController)
	userController := controller2.NewController(conf, slogLogger, provider, passwordProvider, templatedEmailer, client, tenantStorer, userStorer)
	httptransportHandler := httptransport2.NewHandler(slogLogger, userController)
	tenantController := controller3.NewController(conf, slogLogger, provider, templatedEmailer, client, tenantStorer, userStorer)
	handler2 := httptransport3.NewHandler(slogLogger, tenantController)
	ipfsStorager := ipfs.NewStorage(conf, slogLogger)
	pinObjectStorer := datastore3.NewDatastore(conf, slogLogger, client)
	nftAssetStorer := datastore4.NewDatastore(conf, slogLogger, client)
	nftMetadataStorer := datastore5.NewDatastore(conf, slogLogger, client)
	nftCollectionStorer := datastore6.NewDatastore(conf, slogLogger, client)
	nftCollectionController := controller4.NewController(conf, slogLogger, provider, jwtProvider, kmutexProvider, passwordProvider, ipfsStorager, client, tenantStorer, pinObjectStorer, nftAssetStorer, nftMetadataStorer, nftCollectionStorer, userStorer)
	handler3 := httptransport4.NewHandler(slogLogger, nftCollectionController)
	nftMetadataController := controller5.NewController(conf, slogLogger, provider, jwtProvider, kmutexProvider, passwordProvider, ipfsStorager, client, tenantStorer, pinObjectStorer, nftAssetStorer, nftMetadataStorer, nftCollectionStorer, userStorer)
	handler4 := httptransport5.NewHandler(slogLogger, nftMetadataController)
	nftAssetController := controller6.NewController(conf, slogLogger, provider, passwordProvider, jwtProvider, ipfsStorager, client, pinObjectStorer, nftAssetStorer, nftMetadataStorer, nftCollectionStorer, userStorer)
	handler5 := httptransport6.NewHandler(slogLogger, nftAssetController)
	ipfsGatewayController := controller7.NewController(conf, slogLogger, provider, passwordProvider, jwtProvider, ipfsStorager, client, pinObjectStorer)
	handler6 := httptransport7.NewHandler(slogLogger, ipfsGatewayController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, httptransportHandler, handler2, handler3, handler4, handler5, handler6)
	eventschedulerInputPortServer := eventscheduler.NewInputPort(conf, slogLogger, userController, tenantController, nftCollectionController, nftMetadataController, nftAssetController)
	application := NewApplication(slogLogger, inputPortServer, eventschedulerInputPortServer)
	return application
}
