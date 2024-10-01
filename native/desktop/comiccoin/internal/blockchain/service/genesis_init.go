package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

type CreateGenesisBlockDataService struct {
	config                      *config.Config
	logger                      *slog.Logger
	coinbaseAccountKey          *keystore.Key
	setLastBlockDataHashUseCase *usecase.SetLastBlockDataHashUseCase
	createBlockDataUseCase      *usecase.CreateBlockDataUseCase
	proofOfWorkUseCase          *usecase.ProofOfWorkUseCase
}

func NewCreateGenesisBlockDataService(
	config *config.Config,
	logger *slog.Logger,
	coinbaseAccKey *keystore.Key,
	uc1 *usecase.SetLastBlockDataHashUseCase,
	uc2 *usecase.CreateBlockDataUseCase,
	uc3 *usecase.ProofOfWorkUseCase,
) *CreateGenesisBlockDataService {
	return &CreateGenesisBlockDataService{config, logger, coinbaseAccKey, uc1, uc2, uc3}
}

func (s *CreateGenesisBlockDataService) Execute(ctx context.Context) error {
	s.logger.Debug("starting genesis creation service...")
	defer s.logger.Debug("finished  genesis creation service")

	//
	// STEP 1:
	// Setup our very first (signed) transaction: i.e. coinbase giving coins
	// onto the blockchain ... from nothing.
	//

	initialSupply := uint64(5000000000000000000)
	tx := &domain.Transaction{
		ChainID: s.config.Blockchain.ChainID,
		Nonce:   0, // Will be calculated later.
		From:    &s.coinbaseAccountKey.Address,
		To:      &s.coinbaseAccountKey.Address,
		Value:   initialSupply,
		Tip:     0,
		Data:    make([]byte, 0),
	}
	signedTx, err := tx.Sign(s.coinbaseAccountKey.PrivateKey)
	if err != nil {
		return fmt.Errorf("Failed to sign transaction: %v", err)
	}

	// Note: Genesis block has no previous hash
	prevBlockHash := signature.ZeroHash

	gasPrice := uint64(s.config.Blockchain.GasPrice)
	unitsOfGas := uint64(s.config.Blockchain.UnitsOfGas)
	blockTx := domain.BlockTransaction{
		SignedTransaction: signedTx,
		TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
		GasPrice:          gasPrice,
		GasUnits:          unitsOfGas,
	}
	trans := make([]domain.BlockTransaction, 1)
	trans = append(trans, blockTx)

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		return fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	// Construct the genesis block.
	block := domain.Block{
		Header: &domain.BlockHeader{
			Number:        0, // Genesis always starts at zero
			PrevBlockHash: prevBlockHash,
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			Beneficiary:   s.coinbaseAccountKey.Address,
			Difficulty:    s.config.Blockchain.Difficulty,
			MiningReward:  s.config.Blockchain.MiningReward,
			// StateRoot:     "",             //args.StateRoot, // SKIP!
			TransRoot: tree.RootHex(), //
			Nonce:     0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	//
	// STEP 2:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		return fmt.Errorf("Failed to mine block: %v", powErr)
	}

	block.Header.Nonce = nonce

	s.logger.Debug("mining completed",
		slog.Uint64("nonce", block.Header.Nonce))

	//
	// STEP 3
	// Save to database.
	//

	genesisBlockData := domain.NewBlockData(block)

	if err := s.createBlockDataUseCase.Execute(genesisBlockData.Hash, genesisBlockData.Header, genesisBlockData.Trans); err != nil {
		return fmt.Errorf("Failed to create genesis block data: %v", err)
	}

	s.logger.Debug("genesis block created",
		slog.Uint64("nonce", block.Header.Nonce))

	if err := s.setLastBlockDataHashUseCase.Execute(domain.LastBlockDataHash(genesisBlockData.Hash)); err != nil {
		return fmt.Errorf("Failed to save last hash of genesis block data: %v", err)
	}

	return nil
}
