package repo

import (
	"context"
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/domain"
)

type BlockchainStateRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockchainStateRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) domain.BlockchainStateRepository {
	return &BlockchainStateRepo{cfg, logger, db}
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

func (r *BlockchainStateRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *BlockchainStateRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *BlockchainStateRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
