package repo

import (
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage"
)

type BlockchainLastestHashRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockchainLastestHashRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *BlockchainLastestHashRepo {
	return &BlockchainLastestHashRepo{cfg, logger, db}
}

func (r *BlockchainLastestHashRepo) Set(hash string) error {
	hashBytes := []byte(hash)
	if err := r.dbClient.Set("lasthash", hashBytes); err != nil {
		r.logger.Error("failed setting last block data hash into database",
			slog.Any("error", err))
		return fmt.Errorf("failed setting last block data hash into database: %v", err)
	}
	return nil
}

func (r *BlockchainLastestHashRepo) Get() (string, error) {
	bin, err := r.dbClient.Get("lasthash")
	if err != nil {
		r.logger.Error("failed getting last block data hash from database",
			slog.Any("error", err))
		return string(""), fmt.Errorf("failed getting last block data hash from database: %v", err)
	}
	return string(bin), nil
}
