package controller

import (
	"context"
	"log/slog"
)

func (impl *accountControllerImpl) DeleteByName(ctx context.Context, name string) error {
	err := impl.accountStorer.DeleteByName(ctx, name)
	if err != nil {
		impl.logger.Error("failed getting by name",
			slog.Any("error", err))
		return err
	}
	return nil
}
