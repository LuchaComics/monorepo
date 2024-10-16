package repo

import (
	"encoding/json"
	"log/slog"
	"os"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
)

type genesisBlockDataRepoImpl struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewGenesisBlockDataRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) domain.GenesisBlockDataRepository {
	return &genesisBlockDataRepoImpl{cfg, logger, db}
}

func (r *genesisBlockDataRepoImpl) LoadGenesisData() (*domain.GenesisBlockData, error) {
	path := "static/genesis.json"
	content, err := os.ReadFile(path)
	if err != nil {
		r.logger.Error("failed reading file",
			slog.String("path", path),
			slog.Any("error", err))
		return &domain.GenesisBlockData{}, err
	}

	var genesis domain.GenesisBlockData
	err = json.Unmarshal(content, &genesis)
	if err != nil {
		r.logger.Error("failed unmarshalling file",
			slog.String("path", path),
			slog.Any("error", err))
		return &domain.GenesisBlockData{}, err
	}

	return &genesis, nil
}
