package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	mg "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/emailer/mailgun"
	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/templatedemailer"
	attachment_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/datastore"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	receipt_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	userpurchase_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// StoreController Interface for store business logic controller.
type StoreController interface {
	Create(ctx context.Context, m *store_s.Store) (*store_s.Store, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*store_s.Store, error)
	UpdateByID(ctx context.Context, m *store_s.Store) (*store_s.Store, error)
	ListByFilter(ctx context.Context, f *store_s.StorePaginationListFilter) (*store_s.StorePaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *store_s.StorePaginationListFilter) ([]*store_s.StoreAsSelectOption, error)
	PublicListAsSelectOptionByFilter(ctx context.Context, f *store_s.StorePaginationListFilter) ([]*store_s.StoreAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*store_s.Store, error)
}

type StoreControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	S3                    s3_storage.S3Storager
	Emailer               mg.Emailer
	TemplatedEmailer      templatedemailer.TemplatedEmailer
	DbClient              *mongo.Client
	UserStorer            user_s.UserStorer
	StoreStorer           store_s.StoreStorer
	CreditStorer          credit_s.CreditStorer
	AttachmentStorer      attachment_s.AttachmentStorer
	ReceiptStorer         receipt_s.ReceiptStorer
	UserPurchaseStorer    userpurchase_s.UserPurchaseStorer
	ComicSubmissionStorer submission_s.ComicSubmissionStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	emailer mg.Emailer,
	te templatedemailer.TemplatedEmailer,
	client *mongo.Client,
	org_storer store_s.StoreStorer,
	usr_storer user_s.UserStorer,
	st_storer submission_s.ComicSubmissionStorer,
	credit_storer credit_s.CreditStorer,
	attch_storer attachment_s.AttachmentStorer,
	receipt_storer receipt_s.ReceiptStorer,
	up_storer userpurchase_s.UserPurchaseStorer,
) StoreController {
	loggerp.Debug("store controller initialization started...")
	s := &StoreControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		S3:                    s3,
		Emailer:               emailer,
		DbClient:              client,
		TemplatedEmailer:      te,
		UserStorer:            usr_storer,
		ComicSubmissionStorer: st_storer,
		StoreStorer:           org_storer,
		CreditStorer:          credit_storer,
		AttachmentStorer:      attch_storer,
		ReceiptStorer:         receipt_storer,
		UserPurchaseStorer:    up_storer,
	}
	s.Logger.Debug("store controller initialized")
	return s
}
