package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type BlockchainSyncServerService struct {
	config                                       *config.Config
	logger                                       *slog.Logger
	lastBlockDataHashDTOReceiveP2PRequestUseCase *usecase.LastBlockDataHashDTOReceiveP2PRequestUseCase
	getLastBlockDataHashUseCase                  *usecase.GetLastBlockDataHashUseCase
	lastBlockDataHashDTOSendP2PResponseUseCase   *usecase.LastBlockDataHashDTOSendP2PResponseUseCase
}

func NewBlockchainSyncServerService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.LastBlockDataHashDTOReceiveP2PRequestUseCase,
	uc2 *usecase.GetLastBlockDataHashUseCase,
	uc3 *usecase.LastBlockDataHashDTOSendP2PResponseUseCase,
) *BlockchainSyncServerService {
	return &BlockchainSyncServerService{cfg, logger, uc1, uc2, uc3}
}

func (s *BlockchainSyncServerService) Execute(ctx context.Context) error {
	s.logger.Debug("blockchain server running...")
	defer s.logger.Debug("blockchain server ran")

	//
	// STEP 1:
	// Wait to receive any request from the peer-to-peer network.
	//

	peerID, err := s.lastBlockDataHashDTOReceiveP2PRequestUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("failed receiving request",
			slog.Any("error", err))
		return err
	}
	if peerID == "" {
		err := fmt.Errorf("failed receiving request: %v", "peer id d.n.e.")
		s.logger.Error("failed receiving request",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2:
	// Lookup the hash we have locally.
	//

	hash, err := s.getLastBlockDataHashUseCase.Execute()
	if err != nil {
		s.logger.Error("failed getting latest local hash",
			slog.Any("error", err))
		return err
	}
	if hash == "" {
		err := fmt.Errorf("failed getting latest local hash: %v", "hash d.n.e.")
		s.logger.Error("failed getting latest hash",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3
	// Send to the peer the hash we have locally.
	//

	lastHash := domain.LastBlockDataHashDTO(hash)

	if err := s.lastBlockDataHashDTOSendP2PResponseUseCase.Execute(ctx, peerID, lastHash); err != nil {
		s.logger.Error("failed sending response",
			slog.Any("hash", lastHash),
			slog.Any("peer_id", peerID),
			slog.Any("error", err))
		return err
	}

	return nil
}
