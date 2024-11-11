package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/adapter/templatedemailer"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/provider/uuid"
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
	TemplatedEmailer templatedemailer.TemplatedEmailer
	DbClient         *mongo.Client
	UserStorer       user_s.UserStorer
	TenantStorer     tenant_s.TenantStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	te templatedemailer.TemplatedEmailer,
	client *mongo.Client,
	org_tenantr tenant_s.TenantStorer,
	usr_tenantr user_s.UserStorer,
) TenantController {
	loggerp.Debug("tenant controller initialization started...")
	s := &TenantControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		TemplatedEmailer: te,
		DbClient:         client,
		UserStorer:       usr_tenantr,
		TenantStorer:     org_tenantr,
	}
	s.Logger.Debug("tenant controller initialized")
	return s
}
