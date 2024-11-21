package service

import (
	"context"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type BlockDataGetByHashService struct {
	logger              *slog.Logger
	GetBlockDataUseCase *usecase.GetBlockDataUseCase
}

func NewBlockDataGetByHashService(
	logger *slog.Logger,
	uc1 *usecase.GetBlockDataUseCase,
) *BlockDataGetByHashService {
	return &BlockDataGetByHashService{logger, uc1}
}

func (s *BlockDataGetByHashService) Execute(ctx context.Context, hash string) (*domain.BlockData, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if hash == "" {
		e["hash"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get
	//

	blockData, err := s.GetBlockDataUseCase.Execute(ctx, hash)
	if err != nil {
		if !strings.Contains(err.Error(), "does not exist") {
			s.logger.Error("failed getting block data",
				slog.Any("hash", hash),
				slog.Any("error", err))
			return nil, err
		}
	}
	if blockData != nil {
		return blockData, nil
	}

	return nil, nil
}
