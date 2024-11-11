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

type GetBlockchainStateService struct {
	config                    *config.Configuration
	logger                    *slog.Logger
	getBlockchainStateUseCase *usecase.GetBlockchainStateUseCase
}

func NewGetBlockchainStateService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc *usecase.GetBlockchainStateUseCase,
) *GetBlockchainStateService {
	return &GetBlockchainStateService{cfg, logger, uc}
}

func (s *GetBlockchainStateService) Execute(ctx context.Context, chainID uint16) (*domain.BlockchainState, error) {
	blkchState, err := s.getBlockchainStateUseCase.Execute(ctx, chainID)
	if err != nil {
		s.logger.Error("Failed getting blockchain state", slog.Any("error", err))
		return nil, err
	}
	if blkchState == nil {
		errStr := fmt.Sprintf("Blockchain state does not exist for chain ID: %v", chainID)
		s.logger.Error("Failed getting blockchain state", slog.Any("error", errStr))
		return nil, httperror.NewForNotFoundWithSingleField("chain_id", errStr)
	}
	return blkchState, nil
}
