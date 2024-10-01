package repo

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db"
)

type LastBlockDataHashRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient dbase.Database
}

func NewLastBlockDataHashRepo(cfg *config.Config, logger *slog.Logger, db dbase.Database) *LastBlockDataHashRepo {
	return &LastBlockDataHashRepo{cfg, logger, db}
}

func (r *LastBlockDataHashRepo) Set(hash domain.LastBlockDataHash) error {
	hashBytes := []byte(hash)
	if err := r.dbClient.Set("blockdata", "lasthash", hashBytes); err != nil {
		r.logger.Error("failed setting last block data hash into database",
			slog.Any("error", err))
		return fmt.Errorf("failed setting last block data hash into database: %v", err)
	}
	return nil
}

func (r *LastBlockDataHashRepo) Get() (domain.LastBlockDataHash, error) {
	bin, err := r.dbClient.Get("blockdata", "lasthash")
	if err != nil {
		r.logger.Error("failed getting last block data hash from database",
			slog.Any("error", err))
		return domain.LastBlockDataHash(""), fmt.Errorf("failed getting last block data hash from database: %v", err)
	}
	return domain.LastBlockDataHash(bin), nil
}
