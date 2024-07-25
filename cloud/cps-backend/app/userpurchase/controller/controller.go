package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// UserPurchaseController Interface for store business logic controller.
type UserPurchaseController interface {
	Create(ctx context.Context, m *domain.UserPurchase) (*domain.UserPurchase, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.UserPurchase, error)
	UpdateByID(ctx context.Context, m *domain.UserPurchase) (*domain.UserPurchase, error)
	ListByFilter(ctx context.Context, f *domain.UserPurchasePaginationListFilter) (*domain.UserPurchasePaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.UserPurchasePaginationListFilter) ([]*domain.UserPurchaseAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type UserPurchaseControllerImpl struct {
	Config             *config.Conf
	Logger             *slog.Logger
	UUID               uuid.Provider
	DbClient           *mongo.Client
	StoreStorer        store_s.StoreStorer
	UserPurchaseStorer domain.UserPurchaseStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	org_storer store_s.StoreStorer,
	sub_storer domain.UserPurchaseStorer,
) UserPurchaseController {
	s := &UserPurchaseControllerImpl{
		Config:             appCfg,
		Logger:             loggerp,
		UUID:               uuidp,
		DbClient:           client,
		StoreStorer:        org_storer,
		UserPurchaseStorer: sub_storer,
	}
	s.Logger.Debug("store controller initialization started...")
	s.Logger.Debug("store controller initialized")
	return s
}
