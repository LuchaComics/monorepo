package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type BlockDataDTOServerService struct {
	config                                *config.Config
	logger                                *slog.Logger
	blockDataDTOReceiveP2PRequesttUseCase *usecase.BlockDataDTOReceiveP2PRequesttUseCase
	getBlockDataUseCase                   *usecase.GetBlockDataUseCase
	blockDataDTOSendP2PResponsetUseCase   *usecase.BlockDataDTOSendP2PResponsetUseCase
}

func NewBlockDataDTOServerService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.BlockDataDTOReceiveP2PRequesttUseCase,
	uc2 *usecase.GetBlockDataUseCase,
	uc3 *usecase.BlockDataDTOSendP2PResponsetUseCase,
) *BlockDataDTOServerService {
	return &BlockDataDTOServerService{cfg, logger, uc1, uc2, uc3}
}

func (s *BlockDataDTOServerService) Execute(ctx context.Context) error {
	// s.logger.Debug("block data dto server running...")
	// defer s.logger.Debug("block data dto server ran")

	//
	// STEP 1:
	// Wait to receive any request from the peer-to-peer network.
	//

	peerID, blockDataHash, err := s.blockDataDTOReceiveP2PRequesttUseCase.Execute(ctx)
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
	if blockDataHash == "" {
		err := fmt.Errorf("failed receiving request: %v", "hash empty")
		s.logger.Error("failed receiving request",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2:
	// Lookup the hash we have locally.
	//

	blockData, err := s.getBlockDataUseCase.Execute(blockDataHash)
	if err != nil {
		s.logger.Error("failed getting latest local hash",
			slog.Any("error", err))
		return err
	}
	if blockData == nil {
		s.logger.Warn("blockchain has no data to server, exiting...")
		return nil
	}

	//
	// STEP 3
	// Send to the peer the local data we have.
	//

	blockDataDTO := &domain.BlockDataDTO{
		Hash:   blockData.Hash,
		Header: blockData.Header,
		Trans:  blockData.Trans,
	}

	if err := s.blockDataDTOSendP2PResponsetUseCase.Execute(ctx, peerID, blockDataDTO); err != nil {
		s.logger.Error("failed sending response",
			slog.Any("peer_id", peerID),
			slog.Any("error", err))
		return err
	}

	return nil
}
