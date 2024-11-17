package repo

import (
	"context"
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type GenesisBlockDataRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewGenesisBlockDataRepo(logger *slog.Logger, db disk.Storage) domain.GenesisBlockDataRepository {
	return &GenesisBlockDataRepo{logger, db}
}

func (r *GenesisBlockDataRepo) UpsertByChainID(ctx context.Context, genesis *domain.GenesisBlockData) error {
	bBytes, err := genesis.Serialize()
	if err != nil {
		return err
	}
	if err := r.dbClient.Set(fmt.Sprintf("%v", genesis.Header.ChainID), bBytes); err != nil {
		return err
	}
	return nil
}

func (r *GenesisBlockDataRepo) GetByChainID(ctx context.Context, chainID uint16) (*domain.GenesisBlockData, error) {
	bBytes, err := r.dbClient.Get(fmt.Sprintf("%v", chainID))
	if err != nil {
		return nil, err
	}
	b, err := domain.NewGenesisBlockDataFromDeserialize(bBytes)
	if err != nil {
		r.logger.Error("failed to deserialize",
			slog.Any("chainID", chainID),
			slog.String("bin", string(bBytes)),
			slog.Any("error", err))
		return nil, err
	}
	return b, nil
}

func (r *GenesisBlockDataRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *GenesisBlockDataRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *GenesisBlockDataRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
