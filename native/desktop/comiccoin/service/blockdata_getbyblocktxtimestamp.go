package service

import (
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type GetBlockDataByBlockTransactionTimestampService struct {
	config                                         *config.Config
	logger                                         *slog.Logger
	getBlockDataByBlockTransactionTimestampUseCase *usecase.GetBlockDataByBlockTransactionTimestampUseCase
}

func NewGetBlockDataByBlockTransactionTimestampService(
	cfg *config.Config,
	logger *slog.Logger,
	uc *usecase.GetBlockDataByBlockTransactionTimestampUseCase,
) *GetBlockDataByBlockTransactionTimestampService {
	return &GetBlockDataByBlockTransactionTimestampService{cfg, logger, uc}
}

func (s *GetBlockDataByBlockTransactionTimestampService) Execute(timestamp uint64) (*domain.BlockData, error) {
	return s.getBlockDataByBlockTransactionTimestampUseCase.Execute(timestamp)
}
