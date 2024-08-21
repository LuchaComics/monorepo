// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package main

import (
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/cache/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/mongodb"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/templatedemailer"
	controller4 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/controller"
	datastore3 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	httptransport4 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/gateway/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/gateway/httptransport"
	controller3 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/controller"
	datastore2 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	httptransport3 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/httptransport"
	controller2 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	httptransport2 "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/httptransport"
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
	attachmentStorer := datastore3.NewDatastore(conf, slogLogger, client)
	userController := controller2.NewController(conf, slogLogger, provider, passwordProvider, templatedEmailer, client, tenantStorer, userStorer, attachmentStorer)
	httptransportHandler := httptransport2.NewHandler(slogLogger, userController)
	s3Storager := s3.NewStorage(conf, slogLogger, provider)
	tenantController := controller3.NewController(conf, slogLogger, provider, s3Storager, templatedEmailer, client, tenantStorer, userStorer, attachmentStorer)
	handler2 := httptransport3.NewHandler(slogLogger, tenantController)
	attachmentController := controller4.NewController(conf, slogLogger, provider, s3Storager, client, attachmentStorer, userStorer)
	handler3 := httptransport4.NewHandler(slogLogger, attachmentController)
	inputPortServer := http.NewInputPort(conf, slogLogger, middlewareMiddleware, handler, httptransportHandler, handler2, handler3)
	application := NewApplication(slogLogger, inputPortServer)
	return application
}
