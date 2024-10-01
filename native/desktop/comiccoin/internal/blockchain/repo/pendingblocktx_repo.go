package repo

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type PendingBlockTransactionRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient dbase.Database
}

func NewPendingBlockTransactionRepo(cfg *config.Config, logger *slog.Logger, db dbase.Database) *PendingBlockTransactionRepo {
	return &PendingBlockTransactionRepo{cfg, logger, db}
}

func (r *PendingBlockTransactionRepo) Upsert(stx *domain.PendingBlockTransaction) error {
	bBytes, err := stx.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Setf(bBytes, "pending-block-transaction-%v", stx.Nonce); err != nil {
		return err
	}
	return nil
}

func (r *PendingBlockTransactionRepo) ListAll() ([]*domain.PendingBlockTransaction, error) {
	res := make([]*domain.PendingBlockTransaction, 0)
	seekThenIterateKey := ""
	err := r.dbClient.Iterate("pending-block-transaction-", seekThenIterateKey, func(key, value []byte) error {
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
	seekThenIterateKey := ""
	err := r.dbClient.Iterate("pending-block-transaction-", seekThenIterateKey, func(key, value []byte) error {
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
		err := r.dbClient.Deletef("pending-block-transaction-%v", item.Nonce)
		if err != nil {
			return err
		}
		return nil
	}

	return err
}
