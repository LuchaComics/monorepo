package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/templatedemailer"
	attachment_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/attachment/datastore"
	submission_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	credit_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	receipt_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	userpurchase_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// UserController Interface for user business logic controller.
type UserController interface {
	Create(ctx context.Context, requestData *UserCreateRequestIDO) (*user_s.User, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	GetUserBySessionUUID(ctx context.Context, sessionUUID string) (*user_s.User, error)
	ArchiveByID(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	ListByFilter(ctx context.Context, f *user_s.UserPaginationListFilter) (*user_s.UserPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *user_s.UserPaginationListFilter) ([]*user_s.UserAsSelectOption, error)
	UpdateByID(ctx context.Context, request *UserUpdateRequestIDO) (*user_s.User, error)
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*user_s.User, error)
	Star(ctx context.Context, id primitive.ObjectID) (*user_s.User, error)
	ChangePassword(ctx context.Context, req *UserOperationChangePasswordRequest) error
	ChangeTwoFactorAuthentication(ctx context.Context, req *UserOperationChangeTwoFactorAuthenticationRequest) error
	//TODO: Add more...
}

type UserControllerImpl struct {
	Config                *config.Conf
	Logger                *slog.Logger
	UUID                  uuid.Provider
	Password              password.Provider
	DbClient              *mongo.Client
	UserStorer            user_s.UserStorer
	ComicSubmissionStorer submission_s.ComicSubmissionStorer
	StoreStorer           store_s.StoreStorer
	CreditStorer          credit_s.CreditStorer
	AttachmentStorer      attachment_s.AttachmentStorer
	ReceiptStorer         receipt_s.ReceiptStorer
	UserPurchaseStorer    userpurchase_s.UserPurchaseStorer
	TemplatedEmailer      templatedemailer.TemplatedEmailer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	client *mongo.Client,
	org_storer store_s.StoreStorer,
	usr_storer user_s.UserStorer,
	st_storer submission_s.ComicSubmissionStorer,
	credit_storer credit_s.CreditStorer,
	attch_storer attachment_s.AttachmentStorer,
	receipt_storer receipt_s.ReceiptStorer,
	up_storer userpurchase_s.UserPurchaseStorer,
	temailer templatedemailer.TemplatedEmailer,
) UserController {
	s := &UserControllerImpl{
		Config:                appCfg,
		Logger:                loggerp,
		UUID:                  uuidp,
		Password:              passwordp,
		DbClient:              client,
		UserStorer:            usr_storer,
		ComicSubmissionStorer: st_storer,
		StoreStorer:           org_storer,
		CreditStorer:          credit_storer,
		AttachmentStorer:      attch_storer,
		ReceiptStorer:         receipt_storer,
		UserPurchaseStorer:    up_storer,
		TemplatedEmailer:      temailer,
	}
	loggerp.Debug("user controller initialization started...")
	loggerp.Debug("user controller initialized")
	return s
}
