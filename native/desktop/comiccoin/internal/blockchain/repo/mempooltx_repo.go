package repo

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage"
)

type MempoolTransactionRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewMempoolTransactionRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *MempoolTransactionRepo {
	return &MempoolTransactionRepo{cfg, logger, db}
}

func (r *MempoolTransactionRepo) Upsert(stx *domain.MempoolTransaction) error {
	bBytes, err := stx.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", stx.Nonce), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *MempoolTransactionRepo) ListAll() ([]*domain.MempoolTransaction, error) {
	res := make([]*domain.MempoolTransaction, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		stx, err := domain.NewMempoolTransactionFromDeserialize(value)
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

func (r *MempoolTransactionRepo) DeleteAll() error {
	res := make([]*domain.MempoolTransaction, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		stx, err := domain.NewMempoolTransactionFromDeserialize(value)
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
