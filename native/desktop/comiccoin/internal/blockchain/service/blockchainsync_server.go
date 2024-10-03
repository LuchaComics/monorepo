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
	lastBlockDataHashDTOReceiveP2PRequestUseCase *usecase.BlockchainLastestHashDTOReceiveP2PRequestUseCase
	getBlockchainLastestHashUseCase              *usecase.GetBlockchainLastestHashUseCase
	lastBlockDataHashDTOSendP2PResponseUseCase   *usecase.BlockchainLastestHashDTOSendP2PResponseUseCase
}

func NewBlockchainSyncServerService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.BlockchainLastestHashDTOReceiveP2PRequestUseCase,
	uc2 *usecase.GetBlockchainLastestHashUseCase,
	uc3 *usecase.BlockchainLastestHashDTOSendP2PResponseUseCase,
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

	localHash, err := s.getBlockchainLastestHashUseCase.Execute()
	if err != nil {
		s.logger.Error("failed getting latest local hash",
			slog.Any("error", err))
		return err
	}
	if localHash == "" {
		s.logger.Warn("blockchain has no data to server, exiting...")
		return nil
	}

	//
	// STEP 3
	// Send to the peer the hash we have locally.
	//

	lastHash := domain.BlockchainLastestHashDTO(localHash)

	if err := s.lastBlockDataHashDTOSendP2PResponseUseCase.Execute(ctx, peerID, lastHash); err != nil {
		s.logger.Error("failed sending response",
			slog.Any("local_hash", lastHash),
			slog.Any("peer_id", peerID),
			slog.Any("error", err))
		return err
	}

	return nil
}
