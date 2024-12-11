package repo

import (
	"context"
	"log/slog"
	"sort"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type BlockDataRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockDataRepo(logger *slog.Logger, db disk.Storage) domain.BlockDataRepository {
	return &BlockDataRepo{logger, db}
}

func (r *BlockDataRepo) Upsert(ctx context.Context, blockdata *domain.BlockData) error {
	bBytes, err := blockdata.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(blockdata.Hash, bBytes); err != nil {
		return err
	}
	return nil
}

func (r *BlockDataRepo) GetByHash(ctx context.Context, hash string) (*domain.BlockData, error) {
	bBytes, err := r.dbClient.Get(hash)
	if err != nil {
		return nil, err
	}
	b, err := domain.NewBlockDataFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.String("hash", hash),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *BlockDataRepo) ListByChainID(ctx context.Context, chainID uint16) ([]*domain.BlockData, error) {
	res := make([]*domain.BlockData, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		blockdata, err := domain.NewBlockDataFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		if blockdata.Header.ChainID == chainID {
			res = append(res, blockdata)
		}

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *BlockDataRepo) DeleteByHash(ctx context.Context, hash string) error {
	err := r.dbClient.Delete(hash)
	if err != nil {
		return err
	}
	return nil
}

func (r *BlockDataRepo) ListBlockTransactionsByAddress(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	res := make([]*domain.BlockTransaction, 0)
	var counter int64
	err := r.dbClient.Iterate(func(key, value []byte) error {
		blockdata, err := domain.NewBlockDataFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		for _, tx := range blockdata.Trans {
			if strings.ToLower(tx.To.String()) == strings.ToLower(address.String()) || strings.ToLower(tx.From.String()) == strings.ToLower(address.String()) {
				res = append(res, &tx)
				counter++
			}
		}
		return nil
	})

	sort.Slice(res, func(i, j int) bool {
		return res[i].TimeStamp >= res[j].TimeStamp
	})

	return res, err
}

func (r *BlockDataRepo) ListWithLimitForBlockTransactionsByAddress(ctx context.Context, address *common.Address, limit int64) ([]*domain.BlockTransaction, error) {
	res := make([]*domain.BlockTransaction, 0)
	var counter int64
	err := r.dbClient.Iterate(func(key, value []byte) error {
		blockdata, err := domain.NewBlockDataFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		for _, tx := range blockdata.Trans {
			if strings.ToLower(tx.To.String()) == strings.ToLower(address.String()) || strings.ToLower(tx.From.String()) == strings.ToLower(address.String()) {
				res = append(res, &tx)
				counter++
			}
			if counter > limit {
				// Return nil to indicate success because non-nil's indicate error.
				return nil
			}
		}

		// Return nil to indicate success because non-nil's indicate error.
		return nil
	})

	// Sort the results by timestamp
	sort.Slice(res, func(i, j int) bool {
		return res[i].TimeStamp >= res[j].TimeStamp
	})

	return res, err
}

func (r *BlockDataRepo) GetByBlockTransactionTimestamp(ctx context.Context, timestamp uint64) (*domain.BlockData, error) {
	var res *domain.BlockData
	err := r.dbClient.Iterate(func(key, value []byte) error {
		blockdata, err := domain.NewBlockDataFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.Any("timestamp", timestamp),
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		for _, tx := range blockdata.Trans {
			if tx.TimeStamp == timestamp {
				res = blockdata
				return nil // Complete early the loop iteration.
			}
		}

		// Return nil to indicate success because non-nil's indicate error.
		return nil
	})
	return res, err
}

func (r *BlockDataRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *BlockDataRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *BlockDataRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
