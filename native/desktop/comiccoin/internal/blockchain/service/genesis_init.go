package service

import (
	"context"
	"crypto/ecdsa"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

type CreateGenesisBlockDataService struct {
	config                          *config.Config
	logger                          *slog.Logger
	coinbaseAccountKey              *keystore.Key
	getAccountsHashStateUseCase     *usecase.GetAccountsHashStateUseCase
	setBlockchainLastestHashUseCase *usecase.SetBlockchainLastestHashUseCase
	createBlockDataUseCase          *usecase.CreateBlockDataUseCase
	proofOfWorkUseCase              *usecase.ProofOfWorkUseCase
	upsertAccountUseCase            *usecase.UpsertAccountUseCase
}

func NewCreateGenesisBlockDataService(
	config *config.Config,
	logger *slog.Logger,
	coinbaseAccKey *keystore.Key,
	uc1 *usecase.GetAccountsHashStateUseCase,
	uc2 *usecase.SetBlockchainLastestHashUseCase,
	uc3 *usecase.CreateBlockDataUseCase,
	uc4 *usecase.ProofOfWorkUseCase,
	uc5 *usecase.UpsertAccountUseCase,
) *CreateGenesisBlockDataService {
	return &CreateGenesisBlockDataService{config, logger, coinbaseAccKey, uc1, uc2, uc3, uc4, uc5}
}

func (s *CreateGenesisBlockDataService) Execute(ctx context.Context) error {
	s.logger.Debug("starting genesis creation service...")

	// DEVELOPERS NOTE:
	// Here is where we initialize the total supply of coins for the entire
	// blockchain, so adjust accordingly.
	initialSupply := uint64(5000000000000000000)

	//
	// STEP 1
	// Initialize our coinbase account in our in-memory database.
	//

	if err := s.upsertAccountUseCase.Execute(&s.coinbaseAccountKey.Address, initialSupply, 0); err != nil {
		return fmt.Errorf("Failed to upsert account: %v", err)
	}

	//
	// STEP 2:
	// Setup our very first (signed) transaction: i.e. coinbase giving coins
	// onto the blockchain ... from nothing.
	//

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

	//
	// STEP 3:
	// Create our first block, i.e. also called "Genesis block".
	//

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
	trans := make([]domain.BlockTransaction, 0)
	trans = append(trans, blockTx)

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		return fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to create merkle tree",
			slog.Any("error", err))
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
			StateRoot:     stateRoot,
			TransRoot:     tree.RootHex(), //
			Nonce:         0,              // Will be identified by the POW algorithm.
		},
		MerkleTree: tree,
	}

	//
	// STEP 4:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		return fmt.Errorf("Failed to mine block: %v", powErr)
	}

	block.Header.Nonce = nonce

	s.logger.Debug("mining completed",
		slog.Uint64("nonce", block.Header.Nonce))

	genesisBlockData := domain.NewBlockData(block)

	// STEP 5:
	// Create our single proof-of-authority validator via coinbase account.
	//

	coinbasePrivateKey := s.coinbaseAccountKey.PrivateKey
	// Extract the bytes for the original public key.
	coinbasePublicKey := coinbasePrivateKey.Public()
	publicKeyECDSA, ok := coinbasePublicKey.(*ecdsa.PublicKey)
	if !ok {
		return errors.New("error casting public key to ECDSA")
	}
	publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)

	poaValidator := &domain.Validator{
		ID:             "ComicCoin Blockchain Authority",
		PublicKeyBytes: publicKeyBytes,
	}

	//
	// STEP 6:
	// Sign our genesis block's header with our proof-of-authority validator.
	//

	genesisBlockData.Validator = poaValidator
	genesisBlockHeaderSignature, err := poaValidator.Sign(coinbasePrivateKey, genesisBlockData.Header)
	if err != nil {
		return fmt.Errorf("Failed to sign block header: %v", err)
	}
	genesisBlockData.HeaderSignature = genesisBlockHeaderSignature

	//
	// STEP 4:
	// Save to JSON file.
	//

	genesisBlockDataBytes, err := json.MarshalIndent(genesisBlockData, "", "    ")
	if err != nil {
		return fmt.Errorf("Failed to serialize genesis block: %v", err)
	}

	if err := os.WriteFile("static/genesis.json", genesisBlockDataBytes, 0644); err != nil {
		return fmt.Errorf("Failed to write genesis block data to file: %v", err)
	}

	//
	// STEP 5
	// Save to database.
	//

	if err := s.createBlockDataUseCase.Execute(genesisBlockData.Hash, genesisBlockData.Header, genesisBlockData.Trans); err != nil {
		return fmt.Errorf("Failed to write genesis block data to file: %v", err)
	}

	s.logger.Debug("genesis block created",
		slog.String("hash", genesisBlockData.Hash))

	if err := s.setBlockchainLastestHashUseCase.Execute(genesisBlockData.Hash); err != nil {
		return fmt.Errorf("Failed to save last hash of genesis block data: %v", err)
	}

	s.logger.Debug("finished genesis creation service",
		slog.Any("hash", genesisBlockData.Hash))
	return nil
}
