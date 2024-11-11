package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type BlockDataListAllOrderedHashesService struct {
	config                               *config.Configuration
	logger                               *slog.Logger
	listAllBlockNumberByHashArrayUseCase *usecase.ListAllBlockNumberByHashArrayUseCase
}

func NewBlockDataListAllOrderedHashesService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.ListAllBlockNumberByHashArrayUseCase,
) *BlockDataListAllOrderedHashesService {
	return &BlockDataListAllOrderedHashesService{cfg, logger, uc}
}

func (s *BlockDataListAllOrderedHashesService) Execute(ctx context.Context, chainID uint16) ([]domain.BlockNumberByHash, error) {
	return s.listAllBlockNumberByHashArrayUseCase.Execute(ctx, chainID)
}
