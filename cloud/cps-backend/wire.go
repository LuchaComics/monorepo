//go:build wireinject
// +build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/cache/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/emailer/mailgun"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/paymentprocessor/stripe"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/pdfbuilder"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/storage/mongodb"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/templatedemailer"
	attachment_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/controller"
	attachment_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/datastore"
	attachment_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/httptransport"
	comicsub_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/controller"
	comicsub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	comicsub_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/httptransport"
	credit_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/controller"
	credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	credit_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/httptransport"
	customer_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/customer/controller"
	customer_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/customer/httptransport"
	eventlog_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
	gateway_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/controller"
	gateway_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/httptransport"
	off_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/controller"
	off_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	off_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/httptransport"
	strpayproc_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/paymentprocessor/controller/stripe"
	strpayproc_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/paymentprocessor/httptransport/stripe"
	receipt_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/controller"
	receipt_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	receipt_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/httptransport"
	store_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/controller"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	store_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/httptransport"
	user_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/controller"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	user_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/httptransport"
	userpurchase_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/controller"
	userpurchase_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	userpurchase_http "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/httptransport"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/inputport/http"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/inputport/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/cpsrn"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/logger"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/time"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/blacklist"
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
		mailgun.NewEmailer,
		templatedemailer.NewTemplatedEmailer,
		password.NewProvider,
		cpsrn.NewProvider,
		mongodb.NewStorage,
		blacklist.NewProvider,
		mongodbcache.NewCache,
		s3_storage.NewStorage,
		pdfbuilder.NewCBFFBuilder,
		pdfbuilder.NewPCBuilder,
		pdfbuilder.NewCCIMGBuilder,
		pdfbuilder.NewCCSCBuilder,
		pdfbuilder.NewCCBuilder,
		pdfbuilder.NewCCUGBuilder,
		stripe.NewPaymentProcessor,
		eventlog_s.NewDatastore,
		user_s.NewDatastore,
		user_c.NewController,
		customer_c.NewController,
		store_s.NewDatastore,
		store_c.NewController,
		off_s.NewDatastore,
		off_c.NewController,
		receipt_s.NewDatastore,
		receipt_c.NewController,
		userpurchase_s.NewDatastore,
		userpurchase_c.NewController,
		comicsub_s.NewDatastore,
		comicsub_c.NewController,
		strpayproc_c.NewController,
		gateway_c.NewController,
		attachment_s.NewDatastore,
		attachment_c.NewController,
		credit_s.NewDatastore,
		credit_c.NewController,
		strpayproc_http.NewHandler,
		gateway_http.NewHandler,
		user_http.NewHandler,
		receipt_http.NewHandler,
		userpurchase_http.NewHandler,
		customer_http.NewHandler,
		store_http.NewHandler,
		off_http.NewHandler,
		comicsub_http.NewHandler,
		attachment_http.NewHandler,
		credit_http.NewHandler,
		middleware.NewMiddleware,
		http.NewInputPort,
		NewApplication)
	return Application{}
}
