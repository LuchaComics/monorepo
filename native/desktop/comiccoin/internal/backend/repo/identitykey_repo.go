package repo

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage"
)

type IdentityKeyRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewIdentityKeyRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *IdentityKeyRepo {
	return &IdentityKeyRepo{cfg, logger, db}
}

func (r *IdentityKeyRepo) Upsert(identityKey *domain.IdentityKey) error {
	bBytes, err := identityKey.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(identityKey.ID, bBytes); err != nil {
		return err
	}
	return nil
}

func (r *IdentityKeyRepo) GetByID(id string) (*domain.IdentityKey, error) {
	bBytes, err := r.dbClient.Get(id)
	if err != nil {
		return nil, err
	}
	b, err := domain.NewIdentityKeyFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.String("id", id),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}
