package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	ipfs_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/ipfs"
	pin_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	project_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
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
	Create(ctx context.Context, m *project_s.Project) (*project_s.Project, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*project_s.Project, error)
	UpdateByID(ctx context.Context, m *project_s.Project) (*project_s.Project, error)
	ListByFilter(ctx context.Context, f *project_s.ProjectPaginationListFilter) (*project_s.ProjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *project_s.ProjectPaginationListFilter) ([]*project_s.ProjectAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type ProjectControllerImpl struct {
	Config          *config.Conf
	Logger          *slog.Logger
	UUID            uuid.Provider
	JWT             jwt.Provider
	IPFS            ipfs_storage.IPFSStorager
	Kmutex          kmutex.Provider
	Password        password.Provider
	DbClient        *mongo.Client
	TenantStorer    tenant_s.TenantStorer
	ProjectStorer   project_s.ProjectStorer
	PinObjectStorer pin_s.PinObjectStorer
	UserStorer      user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	jwtp jwt.Provider,
	ipfs ipfs_storage.IPFSStorager,
	kmx kmutex.Provider,
	passwordp password.Provider,
	client *mongo.Client,
	tenant_storer tenant_s.TenantStorer,
	project_storer project_s.ProjectStorer,
	pin_storer pin_s.PinObjectStorer,
	usr_storer user_s.UserStorer,
) ProjectController {
	s := &ProjectControllerImpl{
		Config:          appCfg,
		Logger:          loggerp,
		UUID:            uuidp,
		JWT:             jwtp,
		IPFS:            ipfs,
		Kmutex:          kmx,
		Password:        passwordp,
		DbClient:        client,
		TenantStorer:    tenant_storer,
		ProjectStorer:   project_storer,
		PinObjectStorer: pin_storer,
		UserStorer:      usr_storer,
	}
	s.Logger.Debug("project controller initialization started...")
	s.Logger.Debug("project controller initialized")
	return s
}
