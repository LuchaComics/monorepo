package repo

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type BlockDataRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient dbase.Database
}

func NewBlockDataRepo(cfg *config.Config, logger *slog.Logger, db dbase.Database) *BlockDataRepo {
	return &BlockDataRepo{cfg, logger, db}
}

func (r *BlockDataRepo) Upsert(blockdata *domain.BlockData) error {
	bBytes, err := blockdata.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Setf(bBytes, "blockdata-%v", blockdata.Hash); err != nil {
		return err
	}
	return nil
}

func (r *BlockDataRepo) GetByHash(hash string) (*domain.BlockData, error) {
	bBytes, err := r.dbClient.Getf("blockdata-%v", hash)
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

func (r *BlockDataRepo) List() ([]*domain.BlockData, error) {
	res := make([]*domain.BlockData, 0)
	seekThenIterateKey := ""
	err := r.dbClient.Iterate("blockdata-", seekThenIterateKey, func(key, value []byte) error {
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
	err := r.dbClient.Deletef("blockdata-%v", hash)
	if err != nil {
		return err
	}
	return nil
}
