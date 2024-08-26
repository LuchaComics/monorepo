package controller

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"

	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (c *ProjectControllerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Project, error) {
	// Retrieve from our database the record for the specific id.
	m, err := c.ProjectStorer.GetByID(ctx, id)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if m == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "project does not exist")
	}
	return m, err
}
