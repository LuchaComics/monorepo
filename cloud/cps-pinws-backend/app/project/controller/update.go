package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	u_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (c *ProjectControllerImpl) UpdateByID(ctx context.Context, ns *domain.Project) (*domain.Project, error) {
	// Extract from our session the following data.
	urole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	oname, _ := ctx.Value(constants.SessionUserTenantName).(string)
	otz, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		c.Logger.Debug("access granted")
	default:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Fetch the original tenant.
	os, err := c.ProjectStorer.GetByID(ctx, ns.ID)
	if err != nil {
		c.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	if os == nil {
		return nil, httperror.NewForBadRequestWithSingleField("id", "project type does not exist")
	}

	// Modify our original tenant.
	os.TenantID = oid
	os.TenantName = oname
	os.TenantTimezone = otz
	os.ModifiedAt = time.Now()
	os.Status = ns.Status
	// os.Name = ns.Name

	// Save to the database the modified tenant.
	if err := c.ProjectStorer.UpdateByID(ctx, os); err != nil {
		c.Logger.Error("database update by id error", slog.Any("error", err))
		return nil, err
	}

	return os, nil
}
