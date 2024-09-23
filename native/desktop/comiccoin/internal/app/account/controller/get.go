package controller

import (
	"context"
	"fmt"
	"log/slog"
)

func (impl *accountControllerImpl) GetByName(ctx context.Context, name string) (*AccountDetailResponseIDO, error) {
	a, err := impl.accountStorer.GetByName(ctx, name)
	if err != nil {
		impl.logger.Error("failed getting by name",
			slog.Any("error", err))
		return nil, err
	}
	if a == nil {
		impl.logger.Error("failed getting by name",
			slog.Any("name", name),
			slog.Any("error", "does not exist"))
		return nil, fmt.Errorf("does not exist for name: %s", name)
	}
	return &AccountDetailResponseIDO{
		Name:          a.Name,
		WalletAddress: a.WalletAddress.Hex(),
	}, nil
}
