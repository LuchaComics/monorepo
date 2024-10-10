package repo

import (
	"fmt"
	"log/slog"
	"strconv"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage"
)

type BlockchainLastestTokenIDRepo struct {
	config   *config.Config
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewBlockchainLastestTokenIDRepo(cfg *config.Config, logger *slog.Logger, db disk.Storage) *BlockchainLastestTokenIDRepo {
	return &BlockchainLastestTokenIDRepo{cfg, logger, db}
}

func (r *BlockchainLastestTokenIDRepo) Set(tokenID uint64) error {
	tokenIDBytes := []byte(strconv.FormatUint(tokenID, 10))
	if err := r.dbClient.Set("last_token_id", tokenIDBytes); err != nil {
		r.logger.Error("failed setting last token ID into database",
			slog.Any("error", err))
		return fmt.Errorf("failed setting last block data token ID into database: %v", err)
	}
	return nil
}

func (r *BlockchainLastestTokenIDRepo) Get() (uint64, error) {
	bin, err := r.dbClient.Get("last_token_id")
	if err != nil {
		r.logger.Error("failed getting last token ID from database",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed getting last block data token ID from database: %v", err)
	}
	tokenID, err := strconv.ParseUint(string(bin), 10, 64)
	if err != nil {
		r.logger.Error("failed parsing token ID",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed parsing token ID: %v", err)
	}
	return tokenID, nil
}
