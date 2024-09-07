package controller

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"

	s_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/collection/datastore"
	u_d "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (impl *CollectionControllerImpl) Create(ctx context.Context, m *s_d.Collection) (*s_d.Collection, error) {
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
	randomSecretStr, err := impl.Password.GenerateSecureRandomString(64)
	if err != nil {
		impl.Logger.Error("hashing error",
			slog.Any("error", err))
		return nil, err
	}
	randomSecretHash, err := impl.Password.GenerateHashFromPassword(randomSecretStr)
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
	m.SecretHash = randomSecretHash

	// Save to our database.
	if err := impl.CollectionStorer.Create(ctx, m); err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}

	// Generate our one-time API key and attach it to the response. What is
	// important here is that we share the plaintext secret to the user to
	// keep but we do not keep the plaintext value in our system, we only
	// keep the hash, so we keep the value safe.
	apiKeyPayload := fmt.Sprintf("%v@%v", m.ID.Hex(), randomSecretStr)
	atExpiry := 250 * 24 * time.Hour // Duration: 250 years.
	apiKey, _, err := impl.JWT.GenerateJWTToken(apiKeyPayload, atExpiry)
	if err != nil {
		impl.Logger.Error("jwt generate pairs error",
			slog.Any("err", err))
		return nil, err
	}

	m.ApiKey = apiKey

	return m, nil
}
