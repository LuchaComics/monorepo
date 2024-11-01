package repo

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type BlockDataRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockDataRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *BlockDataRepo {
	return &BlockDataRepo{cfg, logger, db}
}

func (r *BlockDataRepo) Upsert(blockdata *domain.BlockData) error {
	bBytes, err := blockdata.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(blockdata.Hash, bBytes); err != nil {
		return err
	}
	return nil
}

func (r *BlockDataRepo) GetByHash(hash string) (*domain.BlockData, error) {
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

func (r *BlockDataRepo) ListAll() ([]*domain.BlockData, error) {
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

func (r *BlockDataRepo) DeleteByHash(hash string) error {
	err := r.dbClient.Delete(hash)
	if err != nil {
		return err
	}
	return nil
}

func (r *BlockDataRepo) ListAllBlockTransactionsByAddress(address *common.Address) ([]*domain.BlockTransaction, error) {
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
			if tx.To.String() == address.String() || tx.From.String() == address.String() {
				res = append(res, &tx)
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
