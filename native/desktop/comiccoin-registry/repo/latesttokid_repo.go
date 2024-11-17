package repo

import (
	"fmt"
	"log/slog"
	"strconv"
	"strings"

	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage"
)

type LastestTokenIDRepo struct {
	logger   *slog.Logger
	dbClient disk.Storage
}

func NewLastestTokenIDRepo(logger *slog.Logger, db disk.Storage) *LastestTokenIDRepo {
	return &LastestTokenIDRepo{logger, db}
}

func (r *LastestTokenIDRepo) Set(tokenID uint64) error {
	tokenIDBytes := []byte(strconv.FormatUint(tokenID, 10))
	if err := r.dbClient.Set("last_token_id", tokenIDBytes); err != nil {
		r.logger.Error("failed setting last token ID into database",
			slog.Any("error", err))
		return fmt.Errorf("failed setting last block data token ID into database: %v", err)
	}
	return nil
}

func (r *LastestTokenIDRepo) Get() (uint64, error) {
	bin, err := r.dbClient.Get("last_token_id")
	if err != nil {
		r.logger.Error("failed getting last token ID from database",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed getting last block data token ID from database: %v", err)
	}

	// If our database is empty then we just return `0` as result.
	if string(bin) == "" {
		return 0, nil
	}

	tokenID, err := strconv.ParseUint(string(bin), 10, 64)
	if err != nil {
		// CASE 1 OF 2: Empty db.
		if strings.Contains(err.Error(), "strconv.ParseUint: parsing \"\"") {
			return 0, nil
		}

		// CASE 2 OF 2: Full db.
		r.logger.Error("failed parsing token ID",
			slog.Any("error", err))
		return 0, fmt.Errorf("failed parsing token ID: %v", err)
	}
	return tokenID, nil
}

func (r *LastestTokenIDRepo) OpenTransaction() error {
	return r.dbClient.OpenTransaction()
}

func (r *LastestTokenIDRepo) CommitTransaction() error {
	return r.dbClient.CommitTransaction()
}

func (r *LastestTokenIDRepo) DiscardTransaction() {
	r.dbClient.DiscardTransaction()
}
