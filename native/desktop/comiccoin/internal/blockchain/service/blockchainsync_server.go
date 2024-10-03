package service

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type BlockchainSyncServerService struct {
	config                              *config.Config
	logger                              *slog.Logger
	blockchainSyncReceiveRequestUseCase *usecase.BlockchainSyncReceiveRequestUseCase
	blockchainSyncSendResponseUseCase   *usecase.BlockchainSyncSendResponseUseCase
}

func NewBlockchainSyncServerService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.BlockchainSyncReceiveRequestUseCase,
	uc2 *usecase.BlockchainSyncSendResponseUseCase,
) *BlockchainSyncServerService {
	return &BlockchainSyncServerService{cfg, logger, uc1, uc2}
}

func (s *BlockchainSyncServerService) Execute(ctx context.Context) error {

	err := s.blockchainSyncReceiveRequestUseCase.Execute()
	log.Println(err)
	return nil
}
