package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type BlockchainSyncClientService struct {
	config                                        *config.Config
	logger                                        *slog.Logger
	lastBlockDataHashDTOSendP2PRequestUseCase     *usecase.LastBlockDataHashDTOSendP2PRequestUseCase
	lastBlockDataHashDTOReceiveP2PResponseUseCase *usecase.LastBlockDataHashDTOReceiveP2PResponseUseCase
	getLastBlockDataHashUseCase                   *usecase.GetLastBlockDataHashUseCase
	setLastBlockDataHashUseCase                   *usecase.SetLastBlockDataHashUseCase
}

func NewBlockchainSyncClientService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.LastBlockDataHashDTOSendP2PRequestUseCase,
	uc2 *usecase.LastBlockDataHashDTOReceiveP2PResponseUseCase,
	uc3 *usecase.GetLastBlockDataHashUseCase,
	uc4 *usecase.SetLastBlockDataHashUseCase,
) *BlockchainSyncClientService {
	return &BlockchainSyncClientService{cfg, logger, uc1, uc2, uc3, uc4}
}

func (s *BlockchainSyncClientService) Execute(ctx context.Context) error {
	s.logger.Debug("blockchain sync client running...")
	defer s.logger.Debug("blockchain sync client ran")

	//
	// STEP 1:
	// Send a request over the peer-to-peer network.
	//

	err := s.lastBlockDataHashDTOSendP2PRequestUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("failed sending request",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 2:
	// Wait to receive request from the peer-to-peer network.
	//

	receivedHash, err := s.lastBlockDataHashDTOReceiveP2PResponseUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("failed receiving response",
			slog.Any("error", err))
		return err
	}
	if receivedHash == "" {
		s.logger.Warn("returned hash is empty")
		return nil
	}

	s.logger.Debug("blockchain sync client received from network",
		slog.Any("network_hash", receivedHash))

	//
	// STEP 3:
	// Get the latest blockchain hash we have in our local database and compare
	// with the received hash and if there are differences then we will need to
	// download from the network the latest blockdata.
	//

	// Note: Do not handle any errors, if we have any errors then continue and
	// fetch the latest hash from network anyway.
	localHash, _ := s.getLastBlockDataHashUseCase.Execute()

	s.logger.Debug("blockchain sync client looked up local hash",
		slog.Any("local_hash", localHash))

	if localHash != string(receivedHash) {
		s.logger.Warn("local blockchain is outdated and needs updating from network",
			slog.Any("network_hash", receivedHash),
			slog.Any("local_hash", localHash))

		//TODO: IMPLEMENT THE CODE HERE...

		return nil
	} else {
		s.logger.Debug("local blockchain is up-to-date")
	}

	return nil
}
