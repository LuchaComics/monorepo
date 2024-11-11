package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type BlockDataListAllUnorderedHashesService struct {
	config                                 *config.Configuration
	logger                                 *slog.Logger
	listBlockDataUnorderedHashArrayUseCase *usecase.ListBlockDataUnorderedHashArrayUseCase
}

func NewBlockDataListAllUnorderedHashesService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.ListBlockDataUnorderedHashArrayUseCase,
) *BlockDataListAllUnorderedHashesService {
	return &BlockDataListAllUnorderedHashesService{cfg, logger, uc}
}

func (s *BlockDataListAllUnorderedHashesService) Execute(ctx context.Context, chainID uint16) ([]string, error) {
	return s.listBlockDataUnorderedHashArrayUseCase.Execute(ctx, chainID)
}
