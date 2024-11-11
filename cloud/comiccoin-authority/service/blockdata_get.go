package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type GetBlockDataService struct {
	config              *config.Configuration
	logger              *slog.Logger
	GetBlockDataUseCase *usecase.GetBlockDataUseCase
}

func NewGetBlockDataService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.GetBlockDataUseCase,
) *GetBlockDataService {
	return &GetBlockDataService{cfg, logger, uc}
}

func (s *GetBlockDataService) Execute(ctx context.Context, hash string) (*domain.BlockData, error) {
	data, err := s.GetBlockDataUseCase.Execute(ctx, hash)
	if err != nil {
		s.logger.Error("Failed getting block data", slog.Any("error", err))
		return nil, err
	}
	if data == nil {
		errStr := fmt.Sprintf("Block data does not exist for hash: %v", hash)
		s.logger.Error("Failed getting block data", slog.Any("error", errStr))
		return nil, httperror.NewForNotFoundWithSingleField("hash", errStr)
	}
	return data, nil
}
