package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	a_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/wallet"
)

func (impl *accountControllerImpl) validateCreateRequest(ctx context.Context, dirtyData *AccountCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData == nil {
		e["name"] = "missing value"
	} else {
		if dirtyData.Name == "" {
			e["name"] = "missing value"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *accountControllerImpl) Create(ctx context.Context, req *AccountCreateRequestIDO) (*AccountDetailResponseIDO, error) {
	if err := impl.validateCreateRequest(ctx, req); err != nil {
		return nil, err
	}

	impl.logger.Debug("Creating new wallet...")
	acc, err := wallet.NewKeystoreAccount(impl.config.App.DirPath, req.WalletPassword)
	if err != nil {
		impl.logger.Error("failed creating new keystore account",
			slog.Any("name", req.Name),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed creating new keystore account: %s", err)
	}

	account := &a_ds.Account{
		Name:          req.Name,
		WalletAddress: acc,
	}
	if insertErr := impl.accountStorer.Insert(ctx, account); insertErr != nil {
		impl.logger.Error("failed inserting new account into database",
			slog.Any("name", req.Name),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed inserting into database: %s", err)
	}

	impl.logger.Debug("New wallet created.", slog.String("account", acc.Hex()))

	return &AccountDetailResponseIDO{
		Name:          req.Name,
		WalletAddress: acc.Hex(),
	}, nil
}
