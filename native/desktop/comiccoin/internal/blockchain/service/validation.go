package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

// ValidationService represents (TODO)
type ValidationService struct {
	config                             *config.Config
	logger                             *slog.Logger
	kmutex                             kmutexutil.KMutexProvider
	receiveProposedBlockDataDTOUseCase *usecase.ReceiveProposedBlockDataDTOUseCase
	getLastBlockDataHashUseCase        *usecase.GetLastBlockDataHashUseCase
	getBlockDataUseCase                *usecase.GetBlockDataUseCase
	createBlockDataUseCase             *usecase.CreateBlockDataUseCase
	setLastBlockDataHashUseCase        *usecase.SetLastBlockDataHashUseCase
}

func NewValidationService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ReceiveProposedBlockDataDTOUseCase,
	uc2 *usecase.GetLastBlockDataHashUseCase,
	uc3 *usecase.GetBlockDataUseCase,
	uc4 *usecase.CreateBlockDataUseCase,
	uc5 *usecase.SetLastBlockDataHashUseCase,
) *ValidationService {
	return &ValidationService{cfg, logger, kmutex, uc1, uc2, uc3, uc4, uc5}
}

func (s *ValidationService) Execute(ctx context.Context) error {
	s.logger.Debug("starting validation service...")
	defer s.logger.Debug("finished validation service")

	//
	// STEP 1
	// Wait to receive data (which also was validated) from the P2P network.
	//

	proposedBlockData, err := s.receiveProposedBlockDataDTOUseCase.Execute(ctx)
	if err != nil {
		s.logger.Error("validator failed receiving dto",
			slog.Any("error", err))
		return err
	}
	if proposedBlockData == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do at the moment is to
		// pause the execution for 1 second and then retry again.
		time.Sleep(1 * time.Second)
		return nil
	}

	s.logger.Debug("received dto from network",
		slog.Any("hash", proposedBlockData.Hash),
	)

	// Lock the validator's database so we coordinate when we receive, validate
	// and/or save to the database.
	s.kmutex.Acquire("validator-service")
	defer s.kmutex.Release("validator-service")

	//
	// STEP 2:
	// Fetch the previous block we have.
	//

	prevBlockDataHash, err := s.getLastBlockDataHashUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to get last hash in database",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get last hash in database: %v", err)
	}
	if prevBlockDataHash == "" {
		s.logger.Error("Blockchain not initialized error")
		return fmt.Errorf("Error: %v", "blockchain not initialized")
	}
	prevBlockData, err := s.getBlockDataUseCase.Execute(string(prevBlockDataHash))
	if err != nil {
		s.logger.Error("Error getting block data frin database",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get lastest block data from database: %v", err)
	}
	if prevBlockData == nil {
		s.logger.Error("Block data does not exist in database",
			slog.Any("hash", prevBlockDataHash))
		return fmt.Errorf("Block data does not exist in database for hash: %v", prevBlockDataHash)
	}
	previousBlock, err := domain.ToBlock(prevBlockData)
	if err != nil {
		s.logger.Error("Error converting block data to block",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3
	// Validate our proposed block data to our blockchain.
	//

	blockData := &domain.BlockData{
		Hash:   proposedBlockData.Hash,
		Header: proposedBlockData.Header,
		Trans:  proposedBlockData.Trans,
	}
	block, err := domain.ToBlock(blockData)
	if err != nil {
		s.logger.Error("validator failed converting block data into a block",
			slog.Any("error", err))
		return err
	}
	stateRoot := "" //TODO: Impl.
	if err := block.ValidateBlock(previousBlock, stateRoot); err != nil {
		s.logger.Error("validator failed validating the proposed block with the previous block",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3:
	// Save to the blockchain database.
	//

	if err := s.createBlockDataUseCase.Execute(blockData.Hash, blockData.Header, blockData.Trans); err != nil {
		s.logger.Error("validator failed saving block data",
			slog.Any("error", err))
		return err
	}

	s.logger.Debug("validator saved proposed block data to blockchain",
		slog.Any("hash", proposedBlockData.Hash),
	)

	return nil
}
