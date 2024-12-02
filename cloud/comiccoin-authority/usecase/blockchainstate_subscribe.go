package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/memory/redis"
)

type BlockchainStateSubscribeUseCase struct {
	logger *slog.Logger
	cache  redis.Cacher
}

func NewBlockchainStateSubscribeUseCase(
	logger *slog.Logger,
	cache redis.Cacher,
) *BlockchainStateSubscribeUseCase {
	return &BlockchainStateSubscribeUseCase{logger, cache}
}

func (uc *BlockchainStateSubscribeUseCase) Execute(ctx context.Context) redis.RedisSubscriber {
	return uc.cache.Subscribe(ctx, "blockchain_state")
}
