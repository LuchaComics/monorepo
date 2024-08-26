package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
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
	d.Status = s_d.StatusArchived
	d.ModifiedAt = time.Now()

	// Save to the database the modified tenant.
	if err := impl.ProjectStorer.UpdateByID(ctx, d); err != nil {
		impl.Logger.Error("database update by id error", slog.Any("error", err))
		return err
	}

	return nil
}
