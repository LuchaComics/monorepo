package repo

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
)

type PendingBlockTransactionRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewPendingBlockTransactionRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *PendingBlockTransactionRepo {
	return &PendingBlockTransactionRepo{cfg, logger, db}
}

func (r *PendingBlockTransactionRepo) Upsert(stx *domain.PendingBlockTransaction) error {
	bBytes, err := stx.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", stx.Nonce), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *PendingBlockTransactionRepo) ListAll() ([]*domain.PendingBlockTransaction, error) {
	res := make([]*domain.PendingBlockTransaction, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		stx, err := domain.NewPendingBlockTransactionFromDeserialize(value)
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

func (r *PendingBlockTransactionRepo) DeleteAll() error {
	res := make([]*domain.PendingBlockTransaction, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		stx, err := domain.NewPendingBlockTransactionFromDeserialize(value)
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
		err := r.dbClient.Delete(fmt.Sprintf("%v", item.Nonce))
		if err != nil {
			return err
		}
		return nil
	}

	return err
}
