package service

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type BlockchainSyncClientService struct {
	config                               *config.Config
	logger                               *slog.Logger
	blockchainSyncSendRequestUseCase     *usecase.BlockchainSyncSendRequestUseCase
	blockchainSyncReceiveResponseUseCase *usecase.BlockchainSyncReceiveResponseUseCase
}

func NewBlockchainSyncClientService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.BlockchainSyncSendRequestUseCase,
	uc2 *usecase.BlockchainSyncReceiveResponseUseCase,
) *BlockchainSyncClientService {
	return &BlockchainSyncClientService{cfg, logger, uc1, uc2}
}

func (s *BlockchainSyncClientService) Execute(ctx context.Context) error {
	err := s.blockchainSyncSendRequestUseCase.Execute()
	log.Println(err)
	return nil
}
