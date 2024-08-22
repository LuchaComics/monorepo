//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/cache/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/emailer/mailgun"
	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/ipfs"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/mongodb"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/templatedemailer"
	attachment_c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/controller"
	attachment_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	attachment_http "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/httptransport"
	gateway_c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/gateway/controller"
	gateway_http "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/gateway/httptransport"
	tenant_c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/controller"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	tenant_http "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/httptransport"
	user_c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/controller"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	user_http "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/inputport/http"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/inputport/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/blacklist"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/logger"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/time"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/uuid"
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
		ipfs_storage.NewStorage,
		user_s.NewDatastore,
		user_c.NewController,
		tenant_s.NewDatastore,
		tenant_c.NewController,
		gateway_c.NewController,
		attachment_s.NewDatastore,
		attachment_c.NewController,
		gateway_http.NewHandler,
		user_http.NewHandler,
		tenant_http.NewHandler,
		attachment_http.NewHandler,
		middleware.NewMiddleware,
		http.NewInputPort,
		NewApplication)
	return Application{}
}
