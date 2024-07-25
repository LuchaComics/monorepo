package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	domain "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/provider/uuid"
)

// ReceiptController Interface for store business logic controller.
type ReceiptController interface {
	Create(ctx context.Context, m *domain.Receipt) (*domain.Receipt, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Receipt, error)
	UpdateByID(ctx context.Context, m *domain.Receipt) (*domain.Receipt, error)
	ListByFilter(ctx context.Context, f *domain.ReceiptPaginationListFilter) (*domain.ReceiptPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.ReceiptPaginationListFilter) ([]*domain.ReceiptAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type ReceiptControllerImpl struct {
	Config        *config.Conf
	Logger        *slog.Logger
	UUID          uuid.Provider
	DbClient      *mongo.Client
	StoreStorer   store_s.StoreStorer
	ReceiptStorer domain.ReceiptStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	client *mongo.Client,
	org_storer store_s.StoreStorer,
	sub_storer domain.ReceiptStorer,
) ReceiptController {
	s := &ReceiptControllerImpl{
		Config:        appCfg,
		Logger:        loggerp,
		UUID:          uuidp,
		DbClient:      client,
		StoreStorer:   org_storer,
		ReceiptStorer: sub_storer,
	}
	s.Logger.Debug("store controller initialization started...")
	s.Logger.Debug("store controller initialized")
	return s
}
