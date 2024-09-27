package datastore

import (
	"context"
	"fmt"
	"log/slog"
)

func (impl *lastHashStorerImpl) Get(ctx context.Context) (string, error) {
	bin, err := impl.dbClient.Get("constant", "lasthash")
	if err != nil {
		impl.logger.Error("failed getting last hash from database", slog.Any("error", err))
		return "", fmt.Errorf("failed getting last hash from database: %v", err)
	}
	return string(bin), nil
}

func (impl *lastHashStorerImpl) Set(ctx context.Context, hash string) error {
	hashBytes := []byte(hash)
	if err := impl.dbClient.Set("constant", "lasthash", hashBytes); err != nil {
		impl.logger.Error("failed setting last hash into database", slog.Any("error", err))
		return fmt.Errorf("failed setting last hash into database: %v", err)
	}
	return nil
}
