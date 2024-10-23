package repo

import (
	_ "embed"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/static"
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
	// DEVELOPERS NOTE:
	// We don't want to be reading local files because if we are using this
	// code from an external package, the file reader code will error. Therefore
	// we will utilize a static embed reader.
	return static.GetGenesisBlockData()
}
