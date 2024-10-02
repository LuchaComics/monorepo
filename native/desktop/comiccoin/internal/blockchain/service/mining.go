package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type MiningService struct {
	config                                  *config.Config
	logger                                  *slog.Logger
	kmutex                                  kmutexutil.KMutexProvider
	listAllPendingBlockTransactionUseCase   *usecase.ListAllPendingBlockTransactionUseCase
	getLastBlockDataHashUseCase             *usecase.GetLastBlockDataHashUseCase
	setLastBlockDataHashUseCase             *usecase.SetLastBlockDataHashUseCase
	getBlockDataUseCase                     *usecase.GetBlockDataUseCase
	createBlockDataUseCase                  *usecase.CreateBlockDataUseCase
	proofOfWorkUseCase                      *usecase.ProofOfWorkUseCase
	broadcastProposedBlockDataDTOUseCase    *usecase.BroadcastProposedBlockDataDTOUseCase
	deleteAllPendingBlockTransactionUseCase *usecase.DeleteAllPendingBlockTransactionUseCase
}

func NewMiningService(
	config *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.ListAllPendingBlockTransactionUseCase,
	uc2 *usecase.GetLastBlockDataHashUseCase,
	uc3 *usecase.SetLastBlockDataHashUseCase,
	uc4 *usecase.GetBlockDataUseCase,
	uc5 *usecase.CreateBlockDataUseCase,
	uc6 *usecase.ProofOfWorkUseCase,
	uc7 *usecase.BroadcastProposedBlockDataDTOUseCase,
	uc8 *usecase.DeleteAllPendingBlockTransactionUseCase,
) *MiningService {
	return &MiningService{config, logger, kmutex, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8}
}

func (s *MiningService) Execute(ctx context.Context) error {
	s.logger.Debug("starting mining service...")
	defer s.logger.Debug("finished mining service")

	//
	// STEP 1:
	// Lock this function - this is important - because it will fetch all the
	// latest pending block transactions, so there needs to be a lockdown of
	// this service that when it runs it will no longer accept any more calls
	// from the system. Therefore we are using a key-based mutex to lock down
	// this service to act as a singleton runtime usage.
	//

	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("mining-service")
	defer s.kmutex.Release("mining-service")

	pendingBlockTxs, err := s.listAllPendingBlockTransactionUseCase.Execute()
	if err != nil {
		s.logger.Error("failed listing pending block transactions",
			slog.Any("error", err))
		return nil
	}
	if len(pendingBlockTxs) <= 0 {
		s.logger.Debug("skipped mining: has no pending block transactions")
		return nil
	}

	s.logger.Debug("executing mining for pending block transactions",
		slog.Int("count", len(pendingBlockTxs)),
	)

	//
	// STEP 2:
	// Lookup the most recent block (data) in our blockchain
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

	//
	// STEP 3:
	// Setup our new block.
	//

	gasPrice := uint64(s.config.Blockchain.GasPrice)
	unitsOfGas := uint64(s.config.Blockchain.UnitsOfGas)
	trans := make([]domain.BlockTransaction, 1)
	for _, pendingBlockTx := range pendingBlockTxs {
		// Create our block.
		blockTx := domain.BlockTransaction{
			SignedTransaction: *pendingBlockTx.MempoolTransaction.ToSignedTransaction(),
			TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
			GasPrice:          gasPrice,
			GasUnits:          unitsOfGas,
		}
		trans = append(trans, blockTx)
	}

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		s.logger.Error("Failed to create merkle tree",
			slog.Any("error", err))
		return fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	// Construct the genesis block.
	block := domain.Block{
		Header: &domain.BlockHeader{
			Number:        prevBlockData.Header.Number + 1,
			PrevBlockHash: string(prevBlockDataHash),
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			Beneficiary:   prevBlockData.Header.Beneficiary,
			Difficulty:    s.config.Blockchain.Difficulty,
			MiningReward:  s.config.Blockchain.MiningReward,
			// StateRoot:     "",             //args.StateRoot, // SKIP!
			TransRoot: tree.RootHex(), //
			Nonce:     0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	//
	// STEP 4:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		s.logger.Error("Failed to mine block",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to mine block: %v", powErr)
	}

	block.Header.Nonce = nonce

	s.logger.Debug("mining completed",
		slog.Uint64("nonce", block.Header.Nonce))

	//
	// STEP 5:
	// Handle case of another miner executing the mine before this
	// current mine and hence our blockchain is out of sync and thus
	// invalidating this current mining operation.

	//TODO: IMPL.

	//
	// STEP 6:
	// Save to database.
	//

	// Convert into saving for our database and transmitting over network.
	blockData := domain.NewBlockData(block)

	if err := s.createBlockDataUseCase.Execute(blockData.Hash, blockData.Header, blockData.Trans); err != nil {
		s.logger.Error("Failed to save block data to blockchain",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to save block data to blockchain: %v", powErr)
	}

	//
	// STEP 7:
	// Broadcast to the distributed / P2P blockchain network our new proposed
	// block data.
	//

	// Convert to our new datastructure.
	purposeBlockData := &domain.ProposedBlockData{
		Hash:   blockData.Hash,
		Header: blockData.Header,
		Trans:  blockData.Trans,
	}

	if err := s.broadcastProposedBlockDataDTOUseCase.Execute(ctx, purposeBlockData); err != nil {
		s.logger.Error("Failed to broadcast proposed block data",
			slog.Any("error", powErr))
		return fmt.Errorf("Failed to broadcast proposed block data: %v", powErr)
	}

	return nil
}
