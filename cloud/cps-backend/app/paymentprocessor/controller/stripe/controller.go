package stripe

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/emailer/mailgun"
	pm "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/paymentprocessor/stripe"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/templatedemailer"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	eventlog_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
	offer_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	r_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	org_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	up_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

type StripePaymentProcessorController interface {
	Webhook(ctx context.Context, header string, b []byte) error
	CreateStripeCheckoutSessionURLForComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (string, error)
}

type StripePaymentProcessorControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Password              password.Provider
	Emailer               mg.Emailer
	TemplatedEmailer      templatedemailer.TemplatedEmailer
	PaymentProcessor      pm.PaymentProcessor
	Kmutex                kmutex.Provider
	DbClient              *mongo.Client
	StoreStorer           org_s.StoreStorer
	UserStorer            user_s.UserStorer
	ReceiptStorer         r_s.ReceiptStorer
	OfferStorer           offer_s.OfferStorer
	EventLogStorer        eventlog_s.EventLogStorer
	ComicSubmissionStorer submission_s.ComicSubmissionStorer
	UserPurchaseStorer    up_s.UserPurchaseStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	emailer mg.Emailer,
	te templatedemailer.TemplatedEmailer,
	paymentProcessor pm.PaymentProcessor,
	kmux kmutex.Provider,
	client *mongo.Client,
	org_storer org_s.StoreStorer,
	sub_storer user_s.UserStorer,
	is r_s.ReceiptStorer,
	offs offer_s.OfferStorer,
	evel eventlog_s.EventLogStorer,
	sub_s submission_s.ComicSubmissionStorer,
	up up_s.UserPurchaseStorer,
) StripePaymentProcessorController {
	loggerp.Debug("payment processor controller initialization started...")
	s := &StripePaymentProcessorControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		S3:                    s3,
		Password:              passwordp,
		Kmutex:                kmux,
		Emailer:               emailer,
		TemplatedEmailer:      te,
		PaymentProcessor:      paymentProcessor,
		DbClient:              client,
		StoreStorer:           org_storer,
		UserStorer:            sub_storer,
		ReceiptStorer:         is,
		OfferStorer:           offs,
		EventLogStorer:        evel,
		ComicSubmissionStorer: sub_s,
		UserPurchaseStorer:    up,
	}
	s.Logger.Debug("payment processor controller initialized")
	return s
}
