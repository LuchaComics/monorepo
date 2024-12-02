package repo

import (
	"fmt"
	"log/slog"
	"math/big"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type PendingSignedTransactionRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewPendingSignedTransactionRepo(logger *slog.Logger, dbClient disk.Storage) *PendingSignedTransactionRepo {
	return &PendingSignedTransactionRepo{logger, dbClient}
}

func (r *PendingSignedTransactionRepo) Upsert(pstx *domain.PendingSignedTransaction) error {
	bBytes, err := pstx.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", pstx.Transaction.GetNonce()), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *PendingSignedTransactionRepo) GetByNonce(pstxNonce *big.Int) (*domain.PendingSignedTransaction, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", pstxNonce))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewPendingSignedTransactionFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("pstx_id", pstxNonce),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *PendingSignedTransactionRepo) ListAll() ([]*domain.PendingSignedTransaction, error) {
	res := make([]*domain.PendingSignedTransaction, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		pstx, err := domain.NewPendingSignedTransactionFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, pstx)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

// func (r *PendingSignedTransactionRepo) ListWithFilterByNonces(tokNonces []*big.Int) ([]*domain.PendingSignedTransaction, error) {
// 	res := make([]*domain.PendingSignedTransaction, 0)
// 	err := r.dbClient.Iterate(func(key, value []byte) error {
// 		pstx, err := domain.NewPendingSignedTransactionFromDeserialize(value)
// 		if err != nil {
// 			r.logger.Error("failed to deserialize",
// 				slog.String("key", string(key)),
// 				slog.String("value", string(value)),
// 				slog.Any("error", err))
// 			return err
// 		}
//
// 		// Apply our filter to the results.
// 		for _, tokNonce := range tokNonces {
// 			if tokNonce == pstx.Nonce {
// 				res = append(res, pstx)
// 			}
// 		}
//
// 		// Return nil to indicate success
// 		return nil
// 	})
//
// 	return res, err
// }

func (r *PendingSignedTransactionRepo) DeleteByNonce(pstxNonce *big.Int) error {
	if err := r.dbClient.Delete(fmt.Sprintf("%v", pstxNonce)); err != nil {
		return err
	}
	return nil
}

func (r *PendingSignedTransactionRepo) OpenTransaction() error {
	if err := r.dbClient.OpenTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *PendingSignedTransactionRepo) CommitTransaction() error {
	if err := r.dbClient.CommitTransaction(); err != nil {
		return err
	}
	return nil
}

func (r *PendingSignedTransactionRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
