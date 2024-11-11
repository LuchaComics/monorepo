package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type ListBlockDataFilteredInHashesService struct {
	config                               *config.Configuration
	logger                               *slog.Logger
	listBlockDataFilteredInHashesUseCase *usecase.ListBlockDataFilteredInHashesUseCase
}

func NewListBlockDataFilteredInHashesService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.ListBlockDataFilteredInHashesUseCase,
) *ListBlockDataFilteredInHashesService {
	return &ListBlockDataFilteredInHashesService{cfg, logger, uc}
}

func (s *ListBlockDataFilteredInHashesService) Execute(ctx context.Context, hashes []string) ([]*domain.BlockData, error) {
	data, err := s.listBlockDataFilteredInHashesUseCase.Execute(ctx, hashes)
	if err != nil {
		s.logger.Error("Failed getting block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
