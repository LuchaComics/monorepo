package repo

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type BlockDataRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockDataRepo(logger *slog.Logger, db disk.Storage) *BlockDataRepo {
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

func (r *BlockDataRepo) ListAll(ctx context.Context) ([]*domain.BlockData, error) {
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

		res = append(res, blockdata)

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

func (r *BlockDataRepo) ListAllBlockTransactionsByAddress(ctx context.Context, address *common.Address) ([]*domain.BlockTransaction, error) {
	res := make([]*domain.BlockTransaction, 0)
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
			}
		}

		// Return nil to indicate success because non-nil's indicate error.
		return nil
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

func (r *BlockDataRepo) ListInHashes(ctx context.Context, hashes []string) ([]*domain.BlockData, error) {
	log.Fatal("TODO: ListInHashes")
	return nil, nil
}
func (r *BlockDataRepo) ListInBetweenBlockNumbersForChainID(ctx context.Context, startBlockNumber, finishBlockNumber uint64, chainID uint16) ([]*domain.BlockData, error) {
	log.Fatal("TODO: ListInBetweenBlockNumbersForChainID")
	return nil, nil
}
func (r *BlockDataRepo) ListBlockNumberByHashArrayForChainID(ctx context.Context, chainID uint16) ([]domain.BlockNumberByHash, error) {
	log.Fatal("TODO: ListBlockNumberByHashArrayForChainID")
	return nil, nil
}
func (r *BlockDataRepo) ListUnorderedHashArrayForChainID(ctx context.Context, chainID uint16) ([]string, error) {
	log.Fatal("TODO: ListUnorderedHashArrayForChainID")
	return nil, nil
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
