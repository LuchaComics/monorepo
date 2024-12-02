package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/memory/redis"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
)

type BlockchainStatePublishUseCase struct {
	logger *slog.Logger
	cache  redis.Cacher
}

func NewBlockchainStatePublishUseCase(
	logger *slog.Logger,
	cache redis.Cacher,
) *BlockchainStatePublishUseCase {
	return &BlockchainStatePublishUseCase{logger, cache}
}

func (uc *BlockchainStatePublishUseCase) Execute(ctx context.Context, bcState *domain.BlockchainState) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if bcState == nil {
		e["blockchain_state"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating blockchain state",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Publish to Redis.
	//

	bcStateBin, err := bcState.Serialize()
	if err != nil {
		uc.logger.Error("failed serializing publish to blockchain state", slog.Any("error", err))
		return err
	}

	if err := uc.cache.Publish(ctx, "blockchain_state", bcStateBin); err != nil {
		uc.logger.Error("failed publishing blockchain state", slog.Any("error", err))
		return err
	}
	uc.logger.Debug("Publishing latest blockchain to redis",
		slog.Any("chain_id", bcState.ChainID))
	return nil
}
