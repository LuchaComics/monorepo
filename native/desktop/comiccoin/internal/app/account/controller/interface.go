package controller

import (
	"context"
	"log/slog"

	a_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type AccountController interface {
	Create(ctx context.Context, b *AccountCreateRequestIDO) (*AccountDetailResponseIDO, error)
	GetByName(ctx context.Context, name string) (*AccountDetailResponseIDO, error)
	List(ctx context.Context) ([]*AccountDetailResponseIDO, error)
	DeleteByName(ctx context.Context, name string) error
}

type accountControllerImpl struct {
	config        *config.Config
	logger        *slog.Logger
	accountStorer a_ds.AccountStorer
}

func NewController(cfg *config.Config, logger *slog.Logger, as a_ds.AccountStorer) AccountController {
	return &accountControllerImpl{
		config:        cfg,
		logger:        logger,
		accountStorer: as,
	}
}
