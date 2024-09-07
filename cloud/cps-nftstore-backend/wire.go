//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/cache/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/mongodb"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/templatedemailer"
	collection_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/controller"
	collection_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/datastore"
	collection_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/httptransport"
	gateway_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/controller"
	gateway_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/gateway/httptransport"
	tenant_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/controller"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	tenant_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/httptransport"
	user_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/controller"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	user_http "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
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
		s3_storage.NewStorage,
		user_s.NewDatastore,
		user_c.NewController,
		collection_s.NewDatastore,
		collection_c.NewController,
		tenant_s.NewDatastore,
		tenant_c.NewController,
		gateway_c.NewController,
		gateway_http.NewHandler,
		user_http.NewHandler,
		collection_http.NewHandler,
		tenant_http.NewHandler,
		middleware.NewMiddleware,
		http.NewInputPort,
		NewApplication)
	return Application{}
}