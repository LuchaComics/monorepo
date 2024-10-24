package repo

import (
	"fmt"
	"log/slog"

	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
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

func (r *BlockchainLastestHashRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *BlockchainLastestHashRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *BlockchainLastestHashRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
