package controller

import (
	"context"
	"log"
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

func NewController(cfg *config.Config, logger *slog.Logger, accountStorer a_ds.AccountStorer) AccountController {
	// For debugging purposes only, we want to print all the accounts (and
	// their wallet addresses) in the console in case the programmer wants
	// to use any of the data.
	aa, err := accountStorer.List(context.Background())
	if err != nil {
		log.Fatalf("failed to list accounts: %v", err)
	}
	for _, a := range aa {
		logger.Debug("account",
			slog.String("name", a.Name),
			slog.String("address", a.WalletAddress.Hex()))
	}

	return &accountControllerImpl{
		config:        cfg,
		logger:        logger,
		accountStorer: accountStorer,
	}
}
