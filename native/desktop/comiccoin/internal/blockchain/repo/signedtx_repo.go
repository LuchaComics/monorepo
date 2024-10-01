package repo

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type SignedTransactionRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient dbase.Database
}

func NewSignedTransactionRepo(cfg *config.Config, logger *slog.Logger, db dbase.Database) *SignedTransactionRepo {
	return &SignedTransactionRepo{cfg, logger, db}
}

func (r *SignedTransactionRepo) Upsert(stx *domain.SignedTransaction) error {
	bBytes, err := stx.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Setf(bBytes, "signed-transaction-%v", stx.Nonce); err != nil {
		return err
	}
	return nil
}

func (r *SignedTransactionRepo) GetByNonce(nonce uint64) (*domain.SignedTransaction, error) {
	bBytes, err := r.dbClient.Getf("signed-transaction-%v", nonce)
	if err != nil {
		return nil, err
	}
	b, err := domain.NewSignedTransactionFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Uint64("nonce", nonce),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *SignedTransactionRepo) ListAll() ([]*domain.SignedTransaction, error) {
	res := make([]*domain.SignedTransaction, 0)
	seekThenIterateKey := ""
	err := r.dbClient.Iterate("signed-transaction-", seekThenIterateKey, func(key, value []byte) error {
		stx, err := domain.NewSignedTransactionFromDeserialize(value)
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

func (r *SignedTransactionRepo) DeleteByNonce(nonce uint64) error {
	err := r.dbClient.Deletef("signed-transaction-%v", nonce)
	if err != nil {
		return err
	}
	return nil
}

func (r *SignedTransactionRepo) DeleteAll() error {
	res := make([]*domain.SignedTransaction, 0)
	seekThenIterateKey := ""
	err := r.dbClient.Iterate("signed-transaction-", seekThenIterateKey, func(key, value []byte) error {
		stx, err := domain.NewSignedTransactionFromDeserialize(value)
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
		if err := r.DeleteByNonce(item.Nonce); err != nil {
			r.logger.Error("failed to delete",
				slog.Any("error", err))
			return err
		}
	}

	return err
}
