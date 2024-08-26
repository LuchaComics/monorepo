package controller

import (
	"context"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (impl *ProjectControllerImpl) Create(ctx context.Context, m *s_d.Project) (*s_d.Project, error) {
	// Extract from our session the following data.
	urole, _ := ctx.Value(constants.SessionUserRole).(int8)
	// uid := ctx.Value(constants.SessionUserID).(primitive.ObjectID)
	// uname := ctx.Value(constants.SessionUserName).(string)
	oid, _ := ctx.Value(constants.SessionUserTenantID).(primitive.ObjectID)
	oname, _ := ctx.Value(constants.SessionUserTenantName).(string)
	otz, _ := ctx.Value(constants.SessionUserTenantTimezone).(string)

	switch urole { // Security.
	case u_d.UserRoleRoot:
		impl.Logger.Debug("access granted")
	default:
		return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
	}

	// Generate hash for the secret.
	randomStr, err := impl.Password.GenerateSecureRandomString(64)
	if err != nil {
		impl.Logger.Error("hashing error",
			slog.Any("error", err))
		return nil, err
	}
	secretHash, err := impl.Password.GenerateHashFromPassword(randomStr)
	if err != nil {
		impl.Logger.Error("hashing error",
			slog.Any("error", err))
		return nil, err
	}

	// Add defaults.
	m.TenantID = oid
	m.TenantName = oname
	m.TenantTimezone = otz
	m.ID = primitive.NewObjectID()
	m.CreatedAt = time.Now()
	// m.CreatedByUserID = uid
	// m.CreatedByUserName = uname
	m.ModifiedAt = time.Now()
	// m.ModifiedByUserID = uid
	// m.ModifiedByUserName = uname
	m.Status = s_d.StatusActive
	m.SecretHashAlgorithm = impl.Password.AlgorithmName()
	m.SecretHash = secretHash

	// Save to our database.
	if err := impl.ProjectStorer.Create(ctx, m); err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	// Attach our secret for one-time.
	m.Secret = randomStr

	return m, nil
}
