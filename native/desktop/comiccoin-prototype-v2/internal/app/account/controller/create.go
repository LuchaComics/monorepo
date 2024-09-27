package controller

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *accountControllerImpl) validateCreateRequest(ctx context.Context, dirtyData *AccountCreateRequestIDO) error {
	e := make(map[string]string)

	if dirtyData == nil {
		e["name"] = "missing value"
	} else {
		if dirtyData.Name == "" {
			e["name"] = "missing value"
		} else {
			account, err := impl.accountStorer.GetByName(context.Background(), dirtyData.Name)
			if err != nil {
				e["name"] = fmt.Sprintf("failed getting account: %v", err)
			}
			if account != nil {
				e["name"] = "already exists"
			}
		}
		if dirtyData.WalletPassword == "" {
			e["wallet_password"] = "missing value"
		}
	}

	if len(e) != 0 {
		impl.logger.Debug("Failed creating new wallet",
			slog.Any("e", e))
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *accountControllerImpl) Create(ctx context.Context, req *AccountCreateRequestIDO) (*AccountDetailResponseIDO, error) {
	if err := impl.validateCreateRequest(ctx, req); err != nil {
		return nil, err
	}
	impl.logger.Debug("Creating new wallet...",
		slog.Any("req", req))

	account, insertErr := impl.accountStorer.Insert(ctx, req.Name, req.WalletPassword)
	if insertErr != nil {
		impl.logger.Error("failed inserting new account into database",
			slog.Any("name", req.Name),
			slog.Any("error", insertErr))
		return nil, fmt.Errorf("failed inserting into database: %s", insertErr)
	}

	impl.logger.Debug("New wallet created.",
		slog.String("account", account.WalletAddress.Hex()),
		slog.String("wallet_filepath", account.WalletFilepath))

	return &AccountDetailResponseIDO{
		Name:          req.Name,
		WalletAddress: account.WalletAddress.Hex(),
	}, nil
}
