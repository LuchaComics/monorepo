package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type MajorityVoteConsensusServerService struct {
	config                                             *config.Config
	logger                                             *slog.Logger
	consensusMechanismReceiveRequestFromNetworkUseCase *usecase.ConsensusMechanismReceiveRequestFromNetworkUseCase
	getBlockchainLastestHashUseCase                    *usecase.GetBlockchainLastestHashUseCase
	consensusMechanismSendResponseToPeerUseCase        *usecase.ConsensusMechanismSendResponseToPeerUseCase
}

func NewMajorityVoteConsensusServerService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ConsensusMechanismReceiveRequestFromNetworkUseCase,
	uc2 *usecase.GetBlockchainLastestHashUseCase,
	uc3 *usecase.ConsensusMechanismSendResponseToPeerUseCase,
) *MajorityVoteConsensusServerService {
	return &MajorityVoteConsensusServerService{cfg, logger, uc1, uc2, uc3}
}

func (s *MajorityVoteConsensusServerService) Execute(ctx context.Context) error {

	//
	// STEP 1:
	// Wait to receive any request from the peer-to-peer network.
	//

	peerID, err := s.consensusMechanismReceiveRequestFromNetworkUseCase.Execute(ctx)
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

	//
	// STEP 3
	// Send to the peer the hash we have locally.
	//

	if err := s.consensusMechanismSendResponseToPeerUseCase.Execute(ctx, peerID, localHash); err != nil {
		s.logger.Error("failed sending response",
			slog.Any("local_hash", localHash),
			slog.Any("peer_id", peerID),
			slog.Any("error", err))
		return err
	}

	return nil
}
