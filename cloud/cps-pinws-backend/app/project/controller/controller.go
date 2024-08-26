package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
	tenant_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/tenant/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/jwt"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/kmutex"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/password"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/uuid"
)

// ProjectController Interface for tenant business logic controller.
type ProjectController interface {
	Create(ctx context.Context, m *domain.Project) (*domain.Project, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Project, error)
	UpdateByID(ctx context.Context, m *domain.Project) (*domain.Project, error)
	ListByFilter(ctx context.Context, f *domain.ProjectPaginationListFilter) (*domain.ProjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.ProjectPaginationListFilter) ([]*domain.ProjectAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type ProjectControllerImpl struct {
	Config        *config.Conf
	Logger        *slog.Logger
	UUID          uuid.Provider
	JWT           jwt.Provider
	Kmutex        kmutex.Provider
	Password      password.Provider
	DbClient      *mongo.Client
	TenantStorer  tenant_s.TenantStorer
	ProjectStorer domain.ProjectStorer
	UserStorer    user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	jwtp jwt.Provider,
	kmx kmutex.Provider,
	passwordp password.Provider,
	client *mongo.Client,
	org_tenantr tenant_s.TenantStorer,
	sub_tenantr domain.ProjectStorer,
	usr_storer user_s.UserStorer,
) ProjectController {
	s := &ProjectControllerImpl{
		Config:        appCfg,
		Logger:        loggerp,
		UUID:          uuidp,
		JWT:           jwtp,
		Kmutex:        kmx,
		Password:      passwordp,
		DbClient:      client,
		TenantStorer:  org_tenantr,
		ProjectStorer: sub_tenantr,
		UserStorer:    usr_storer,
	}
	s.Logger.Debug("project controller initialization started...")
	s.Logger.Debug("project controller initialized")
	return s
}
