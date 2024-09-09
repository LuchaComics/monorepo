package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *NFTControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	d, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if d == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "nft does not exist")
	}

	if err := impl.NFTStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
