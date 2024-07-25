package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pm "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/paymentprocessor/stripe"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/pdfbuilder"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/templatedemailer"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// CustomerController Interface for customer business logic controller.
type CustomerController interface {
	Create(ctx context.Context, m *CustomerCreateRequestIDO) (*user_s.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	UpdateByID(ctx context.Context, m *user_s.User) (*user_s.User, error)
	ListByFilter(ctx context.Context, f *user_s.UserPaginationListFilter) (*user_s.UserPaginationListResult, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error)
	Star(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
}

type CustomerControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Password              password.Provider
	PaymentProcessor      pm.PaymentProcessor
	CBFFBuilder           pdfbuilder.CBFFBuilder
	DbClient              *mongo.Client
	UserStorer            user_s.UserStorer
	ComicSubmissionStorer submission_s.ComicSubmissionStorer
	TemplatedEmailer      templatedemailer.TemplatedEmailer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	passwordp password.Provider,
	paymentProcessor pm.PaymentProcessor,
	cbffb pdfbuilder.CBFFBuilder,
	temailer templatedemailer.TemplatedEmailer,
	client *mongo.Client,
	u_storer user_s.UserStorer,
	sub_storer submission_s.ComicSubmissionStorer,
) CustomerController {
	s := &CustomerControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		S3:                    s3,
		Password:              passwordp,
		PaymentProcessor:      paymentProcessor,
		CBFFBuilder:           cbffb,
		TemplatedEmailer:      temailer,
		DbClient:              client,
		UserStorer:            u_storer,
		ComicSubmissionStorer: sub_storer,
	}
	s.Logger.Debug("customer controller initialization started...")
	s.Logger.Debug("customer controller initialized")
	return s
}
