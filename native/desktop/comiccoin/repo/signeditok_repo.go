package repo

import (
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type SignedIssuedTokenRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewSignedIssuedTokenRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) domain.SignedIssuedTokenRepository {
	return &SignedIssuedTokenRepo{cfg, logger, db}
}

func (r *SignedIssuedTokenRepo) Upsert(stx *domain.SignedIssuedToken) error {
	bBytes, err := stx.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", stx.ID), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *SignedIssuedTokenRepo) GetByID(id uint64) (*domain.SignedIssuedToken, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", id))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewSignedIssuedTokenFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Uint64("id", id),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *SignedIssuedTokenRepo) ListAll() ([]*domain.SignedIssuedToken, error) {
	res := make([]*domain.SignedIssuedToken, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		stx, err := domain.NewSignedIssuedTokenFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, stx)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *SignedIssuedTokenRepo) DeleteAll() error {
	res := make([]*domain.SignedIssuedToken, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		stx, err := domain.NewSignedIssuedTokenFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, stx)

		// Return nil to indicate success
		return nil
	})

	for _, item := range res {
		err := r.dbClient.Delete(fmt.Sprintf("%v", item.ID))
		if err != nil {
			return err
		}
		return nil
	}

	return err
}

func (r *SignedIssuedTokenRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *SignedIssuedTokenRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *SignedIssuedTokenRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
