package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/credit/datastore"
	offer_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// CreditController Interface for store business logic controller.
type CreditController interface {
	Create(ctx context.Context, req *CreditCreateRequest) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Credit, error)
	UpdateByID(ctx context.Context, m *domain.Credit) (*domain.Credit, error)
	ListByFilter(ctx context.Context, f *domain.CreditPaginationListFilter) (*domain.CreditPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.CreditPaginationListFilter) ([]*domain.CreditAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type CreditControllerImpl struct {
	Config       *config.Conf
	Logger       *slog.Logger
	UUID         uuid.Provider
	DbClient     *mongo.Client
	StoreStorer  store_s.StoreStorer
	CreditStorer domain.CreditStorer
	UserStorer   user_s.UserStorer
	OfferStorer  offer_s.OfferStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	org_storer store_s.StoreStorer,
	credit_storer domain.CreditStorer,
	usr_storer user_s.UserStorer,
	offer_storer offer_s.OfferStorer,
) CreditController {
	loggerp.Debug("store controller initialization started...")
	s := &CreditControllerImpl{
		Config:       appCfg,
		Logger:       loggerp,
		UUID:         uuidp,
		DbClient:     client,
		StoreStorer:  org_storer,
		CreditStorer: credit_storer,
		UserStorer:   usr_storer,
		OfferStorer:  offer_storer,
	}
	s.Logger.Debug("store controller initialized")
	return s
}
