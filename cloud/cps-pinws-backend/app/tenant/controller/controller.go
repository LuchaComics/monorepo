package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/s3"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/templatedemailer"
	pinobject_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/uuid"
)

// TenantController Interface for tenant business logic controller.
type TenantController interface {
	Create(ctx context.Context, m *tenant_s.Tenant) (*tenant_s.Tenant, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*tenant_s.Tenant, error)
	UpdateByID(ctx context.Context, m *tenant_s.Tenant) (*tenant_s.Tenant, error)
	ListByFilter(ctx context.Context, f *tenant_s.TenantPaginationListFilter) (*tenant_s.TenantPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *tenant_s.TenantPaginationListFilter) ([]*tenant_s.TenantAsSelectOption, error)
	PublicListAsSelectOptionByFilter(ctx context.Context, f *tenant_s.TenantPaginationListFilter) ([]*tenant_s.TenantAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CreateComment(ctx context.Context, customerID primitive.ObjectID, content string) (*tenant_s.Tenant, error)
}

type TenantControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	S3               s3_storage.S3Storager
	TemplatedEmailer templatedemailer.TemplatedEmailer
	DbClient         *mongo.Client
	UserStorer       user_s.UserStorer
	TenantStorer     tenant_s.TenantStorer
	PinObjectStorer pinobject_s.PinObjectStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	te templatedemailer.TemplatedEmailer,
	client *mongo.Client,
	org_tenantr tenant_s.TenantStorer,
	usr_tenantr user_s.UserStorer,
	attch_tenantr pinobject_s.PinObjectStorer,
) TenantController {
	loggerp.Debug("tenant controller initialization started...")
	s := &TenantControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		S3:               s3,
		TemplatedEmailer: te,
		DbClient:         client,
		UserStorer:       usr_tenantr,
		TenantStorer:     org_tenantr,
		PinObjectStorer: attch_tenantr,
	}
	s.Logger.Debug("tenant controller initialized")
	return s
}
