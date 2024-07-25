package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (c *StoreControllerImpl) Create(ctx context.Context, m *s_d.Store) (*s_d.Store, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	userName := ctx.Value(constants.SessionUserName).(string)
	userRole := ctx.Value(constants.SessionUserRole).(int8)

	// Apply protection based on ownership and role.
	if userRole != user_d.UserRoleRoot {
		c.Logger.Error("authenticated user is not staff role error",
			slog.Any("role", userRole),
			slog.Any("userID", userID))
		return nil, httperror.NewForForbiddenWithSingleField("message", "you role does not grant you access to this")
	}

	// Add defaults.
	m.ID = primitive.NewObjectID()
	m.CreatedByUserID = userID
	m.CreatedByUserName = userName
	m.CreatedAt = time.Now()
	m.ModifiedByUserID = userID
	m.ModifiedByUserName = userName
	m.ModifiedAt = time.Now()

	// Save to our database.
	err := c.StoreStorer.Create(ctx, m)
	if err != nil {
		c.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	// Send notifications in the background.
	if m.Status == s_d.StoreActiveStatus {
		c.Logger.Debug("store became active, sending email to retailer staff")
		go func(m *s_d.Store) {
			res, err := c.UserStorer.ListAllRetailerStaffForStoreID(context.Background(), m.ID)
			if err != nil {
				c.Logger.Error("list store error", slog.Any("error", err))
				return
			}
			var retailerEmails []string
			for _, u := range res.Results {
				retailerEmails = append(retailerEmails, u.Email)
			}
			if err := c.TemplatedEmailer.SendRetailerStoreActiveEmailToRetailers(retailerEmails, m.Name); err != nil {
				c.Logger.Error("failed sending templated error", slog.Any("error", err))
				return
			}

		}(m)
	}

	return m, nil
}
