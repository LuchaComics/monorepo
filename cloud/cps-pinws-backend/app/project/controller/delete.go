package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (impl *ProjectControllerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	d, err := impl.GetByID(ctx, id)
	if err != nil {
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return err
	}
	if d == nil {
		impl.Logger.Error("database returns nothing from get by id")
		return httperror.NewForBadRequestWithSingleField("id", "project does not exist")
	}

	// Get all the pinned objects and delete them from our local IPFS node's file system and unpin from sharing in IPFS network.
	pp, err := impl.PinObjectStorer.ListByProjectID(ctx, id)
	for _, p := range pp.Results {
		if err := impl.IPFS.DeleteContent(ctx, p.CID); err != nil {
			impl.Logger.Error("ipfs failed to delete",
				slog.String("cid", p.CID),
				slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("deleted pin from ipfs", slog.String("cid", p.CID))
		if err := impl.PinObjectStorer.DeleteByRequestID(ctx, p.RequestID); err != nil {
			impl.Logger.Error("database failed to delete pin",
				slog.String("cid", p.CID),
				slog.String("requestid", p.RequestID.Hex()),
				slog.Any("error", err))
			return err
		}
		impl.Logger.Debug("deleted pin from database", slog.String("requestid", p.RequestID.Hex()))
	}

	if err := impl.ProjectStorer.DeleteByID(ctx, id); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
