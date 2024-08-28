package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/adapter/templatedemailer"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/provider/uuid"
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
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	Password         password.Provider
	TemplatedEmailer templatedemailer.TemplatedEmailer
	DbClient         *mongo.Client
	UserStorer       user_s.UserStorer
	TenantStorer      tenant_s.TenantStorer
	PinObjectStorer pinobject_s.PinObjectStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	passwordp password.Provider,
	te templatedemailer.TemplatedEmailer,
	client *mongo.Client,
	org_storer tenant_s.TenantStorer,
	usr_storer user_s.UserStorer,
	attch_storer pinobject_s.PinObjectStorer,
) UserController {
	s := &UserControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		Password:         passwordp,
		TemplatedEmailer: te,
		DbClient:         client,
		UserStorer:       usr_storer,
		TenantStorer:      org_storer,
		PinObjectStorer: attch_storer,
	}
	loggerp.Debug("user controller initialization started...")
	loggerp.Debug("user controller initialized")
	return s
}
