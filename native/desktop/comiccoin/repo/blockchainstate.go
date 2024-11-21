package repo

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type BlockchainStateRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockchainStateRepo(logger *slog.Logger, db disk.Storage) domain.BlockchainStateRepository {
	return &BlockchainStateRepo{logger, db}
}

func (r *BlockchainStateRepo) UpsertByChainID(ctx context.Context, genesis *domain.BlockchainState) error {
	bBytes, err := genesis.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", genesis.ChainID), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *BlockchainStateRepo) GetByChainID(ctx context.Context, chainID uint16) (*domain.BlockchainState, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", chainID))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewBlockchainStateFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("chainID", chainID),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *BlockchainStateRepo) DeleteByChainID(ctx context.Context, chainID uint16) error {
	err := r.dbClient.Delete(fmt.Sprintf("%v", chainID))
	if err != nil {
		return err
	}
	return nil
}

func (r *BlockchainStateRepo) ListAll(ctx context.Context) ([]*domain.BlockchainState, error) {
	res := make([]*domain.BlockchainState, 0)
	err := r.dbClient.Iterate(func(key, value []byte) error {
		account, err := domain.NewBlockchainStateFromDeserialize(value)
		if err != nil {
			r.logger.Error("failed to deserialize",
				slog.String("key", string(key)),
				slog.String("value", string(value)),
				slog.Any("error", err))
			return err
		}

		res = append(res, account)

		// Return nil to indicate success
		return nil
	})

	return res, err
}

func (r *BlockchainStateRepo) GetUpdateChangeStreamChannel(context.Context) (<-chan domain.BlockchainState, chan struct{}, error) {
	defer log.Fatal("Unsupported feature - only supported in `comiccoin-authority` application.")
	return nil, nil, nil
}

func (r *BlockchainStateRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *BlockchainStateRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *BlockchainStateRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
