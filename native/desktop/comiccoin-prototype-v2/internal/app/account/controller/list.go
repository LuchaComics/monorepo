package controller

import (
	"context"
	"log/slog"
)

func (impl *accountControllerImpl) List(ctx context.Context) ([]*AccountDetailResponseIDO, error) {
	aa, err := impl.accountStorer.List(ctx)
	if err != nil {
		impl.logger.Error("failed getting by name",
			slog.Any("error", err))
		return nil, err
	}

	res := make([]*AccountDetailResponseIDO, len(aa))
	for _, a := range aa {
		datum := &AccountDetailResponseIDO{
			Name:          a.Name,
			WalletAddress: a.WalletAddress.Hex(),
		}
		res = append(res, datum)
	}

	return res, nil
}
