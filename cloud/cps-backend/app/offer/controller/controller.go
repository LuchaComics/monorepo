package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	pm "github.com/LuchaComics/monorepo/cloud/cps-backend/adapter/paymentprocessor/stripe"
	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// Offerontroller Interface for store business logic controller.
type Offerontroller interface {
	// Create(ctx context.Context, m *domain.Offer) (*domain.Offer, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Offer, error)
	GetByServiceType(ctx context.Context, serviceType int8) (*domain.Offer, error)
	UpdateByID(ctx context.Context, m *domain.Offer) (*domain.Offer, error)
	ListByFilter(ctx context.Context, f *domain.OfferPaginationListFilter) (*domain.OfferPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.OfferPaginationListFilter) ([]*domain.OfferAsSelectOption, error)
	// DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type OfferControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	DbClient         *mongo.Client
	PaymentProcessor pm.PaymentProcessor
	StoreStorer      store_s.StoreStorer
	OfferStorer      domain.OfferStorer
	UserStorer       user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	paymentProcessor pm.PaymentProcessor,
	org_storer store_s.StoreStorer,
	sub_storer domain.OfferStorer,
	usr_storer user_s.UserStorer,
) Offerontroller {
	s := &OfferControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		DbClient:         client,
		PaymentProcessor: paymentProcessor,
		StoreStorer:      org_storer,
		OfferStorer:      sub_storer,
		UserStorer:       usr_storer,
	}
	s.Logger.Debug("offer controller initialization started...")
	// if err := s.createDefaults(context.Background()); err != nil {
	// 	log.Fatal(err)
	// }
	s.Logger.Debug("offer controller initialized")
	return s
}
