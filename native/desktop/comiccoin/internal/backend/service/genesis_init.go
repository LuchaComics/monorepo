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

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/backend/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/merkle"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/blockchain/signature"
)

type CreateGenesisBlockDataService struct {
	config                                       *config.Config
	logger                                       *slog.Logger
	coinbaseAccountKey                           *keystore.Key
	getAccountsHashStateUseCase                  *usecase.GetAccountsHashStateUseCase
	getTokensHashStateUseCase                    *usecase.GetTokensHashStateUseCase
	setBlockchainLastestHashUseCase              *usecase.SetBlockchainLastestHashUseCase
	setBlockchainLastestTokenIDIfGreatestUseCase *usecase.SetBlockchainLastestTokenIDIfGreatestUseCase
	createBlockDataUseCase                       *usecase.CreateBlockDataUseCase
	proofOfWorkUseCase                           *usecase.ProofOfWorkUseCase
	upsertAccountUseCase                         *usecase.UpsertAccountUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase    *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
}

func NewCreateGenesisBlockDataService(
	config *config.Config,
	logger *slog.Logger,
	coinbaseAccKey *keystore.Key,
	uc1 *usecase.GetAccountsHashStateUseCase,
	uc2 *usecase.GetTokensHashStateUseCase,
	uc3 *usecase.SetBlockchainLastestHashUseCase,
	uc4 *usecase.SetBlockchainLastestTokenIDIfGreatestUseCase,
	uc5 *usecase.CreateBlockDataUseCase,
	uc6 *usecase.ProofOfWorkUseCase,
	uc7 *usecase.UpsertAccountUseCase,
	uc8 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
) *CreateGenesisBlockDataService {
	return &CreateGenesisBlockDataService{config, logger, coinbaseAccKey, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8}
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

	// DEVELOPERS NOTE:
	// During genesis block creation, the account's nonce value is indeed 0.
	//
	// After the genesis block is mined, the account's nonce value is
	// incremented to 1.
	//
	// This makes sense because the genesis block is the first block in the
	// blockchain, and the account's nonce value is used to track the number of
	// transactions sent from that account.
	//
	// Since the genesis block is the first transaction sent from the account,
	// the nonce value is incremented from 0 to 1 after the block is mined.
	//
	// Here's a step-by-step breakdown:
	//
	// 1. Genesis block creation:
	// --> Account's nonce value is 0.
	// 2. Genesis block mining:
	// --> Account's nonce value is still 0.
	// 3. Genesis block is added to the blockchain:
	// --> Account's nonce value is now 1.
	//
	// From this point on, every time a transaction is sent from the account, the nonce value is incremented by 1.

	if err := s.upsertAccountUseCase.Execute(&s.coinbaseAccountKey.Address, initialSupply, 0); err != nil {
		return fmt.Errorf("Failed to upsert account: %v", err)
	}

	//
	// STEP 2:
	// Setup our very first (signed) transaction: i.e. coinbase giving coins
	// onto the blockchain ... from nothing.
	//

	coinTx := &domain.Transaction{
		ChainID: s.config.Blockchain.ChainID,
		Nonce:   0, // Will be calculated later.
		From:    &s.coinbaseAccountKey.Address,
		To:      &s.coinbaseAccountKey.Address,
		Value:   initialSupply,
		Tip:     0,
		Data:    make([]byte, 0),
		Type:    domain.TransactionTypeCoin,
	}
	signedCoinTx, err := coinTx.Sign(s.coinbaseAccountKey.PrivateKey)
	if err != nil {
		return fmt.Errorf("Failed to sign coin transaction: %v", err)
	}

	//
	// STEP 3:
	// Setup our very first (signed) token transaction.
	//

	tokenTx := &domain.Transaction{
		ChainID:          s.config.Blockchain.ChainID,
		Nonce:            0, // Will be calculated later.
		From:             &s.coinbaseAccountKey.Address,
		To:               &s.coinbaseAccountKey.Address,
		Value:            0, //Note: Tokens don't have coin value.
		Tip:              0,
		Data:             make([]byte, 0),
		Type:             domain.TransactionTypeToken,
		TokenID:          0, // The very first token in our entire blockchain starts at the value of zero.
		TokenMetadataURI: "https://cpscapsule.com/comiccoin/tokens/0/metadata.json",
		TokenNonce:       0, // Newly minted tokens always have their nonce start at value of zero.
	}
	signedTokenTx, err := tokenTx.Sign(s.coinbaseAccountKey.PrivateKey)
	if err != nil {
		return fmt.Errorf("Failed to sign token transaction: %v", err)
	}

	nftFromAddr, err := signedTokenTx.FromAddress()
	if err != nil {
		s.logger.Error("Failed getting from address",
			slog.Any("chain_id", s.config.Blockchain.ChainID),
			slog.Any("error", err))
		return err
	}

	s.logger.Info("Created first token",
		slog.Any("from", signedTokenTx.From),
		slog.Any("from_via_sig", nftFromAddr),
		slog.Any("to", signedTokenTx.To),
		slog.Any("tx_sig_v", signedTokenTx.V),
		slog.Any("tx_sig_r", signedTokenTx.R),
		slog.Any("tx_sig_s", signedTokenTx.S),
		slog.Uint64("tx_token_id", signedTokenTx.TokenID))

	// Defensive code: Run this code to ensure this transaction is
	// properly structured for our blockchain.
	if err := signedTokenTx.Validate(s.config.Blockchain.ChainID, true); err != nil {
		s.logger.Error("Failed token transaction.",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 4:
	// Save our token to our database.
	//

	if err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(tokenTx.TokenID, tokenTx.To, tokenTx.TokenMetadataURI, tokenTx.TokenNonce); err != nil {
		return err
	}
	if err := s.setBlockchainLastestTokenIDIfGreatestUseCase.Execute(tokenTx.TokenID); err != nil {
		return fmt.Errorf("Failed to save last token ID of genesis block data: %v", err)
	}

	//
	// STEP 5:
	// Create our first block, i.e. also called "Genesis block".
	//

	// Note: Genesis block has no previous hash
	prevBlockHash := signature.ZeroHash

	gasPrice := uint64(s.config.Blockchain.GasPrice)
	unitsOfGas := uint64(s.config.Blockchain.UnitsOfGas)
	coinBlockTx := domain.BlockTransaction{
		SignedTransaction: signedCoinTx,
		TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
		GasPrice:          gasPrice,
		GasUnits:          unitsOfGas,
	}
	tokenBlockTx := domain.BlockTransaction{
		SignedTransaction: signedTokenTx,
		TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
		GasPrice:          gasPrice,
		GasUnits:          unitsOfGas,
	}
	trans := make([]domain.BlockTransaction, 0)
	trans = append(trans, coinBlockTx)
	trans = append(trans, tokenBlockTx)

	// Construct a merkle tree from the transaction for this block. The root
	// of this tree will be part of the block to be mined.
	tree, err := merkle.NewTree(trans)
	if err != nil {
		return fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to get hash of all accounts",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get hash of all accounts: %v", err)
	}

	// Running this code get's a hash of all the tokens, thus making the
	// tokens tamper proof.
	tokensRoot, err := s.getTokensHashStateUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed to get hash of all tokens",
			slog.Any("error", err))
		return fmt.Errorf("Failed to get hash of all tokens: %v", err)
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
			LatestTokenID: 0,              // ComicCoin: Token ID values start at zero.
			TokensRoot:    tokensRoot,
		},
		MerkleTree: tree,
	}

	genesisBlockData := domain.NewBlockData(block)

	//
	// STEP 6:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		return fmt.Errorf("Failed to mine block: %v", powErr)
	}

	block.Header.Nonce = nonce

	s.logger.Debug("genesis mining completed",
		slog.Uint64("nonce", block.Header.Nonce))

	// STEP 7:
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
	// STEP 8:
	// Sign our genesis block's header with our proof-of-authority validator.
	// Note: Signing always happens after the miner has found the `nonce` in
	// the block header.
	//

	genesisBlockData.Validator = poaValidator
	genesisBlockHeaderSignatureBytes, err := poaValidator.Sign(coinbasePrivateKey, genesisBlockData.Header)
	if err != nil {
		return fmt.Errorf("Failed to sign block header: %v", err)
	}
	genesisBlockData.HeaderSignatureBytes = genesisBlockHeaderSignatureBytes

	//
	// STEP 9:
	// Save genesis block to a JSON file.
	//

	genesisBlockDataBytes, err := json.MarshalIndent(genesisBlockData, "", "    ")
	if err != nil {
		return fmt.Errorf("Failed to serialize genesis block: %v", err)
	}

	if err := os.WriteFile("static/genesis.json", genesisBlockDataBytes, 0644); err != nil {
		return fmt.Errorf("Failed to write genesis block data to file: %v", err)
	}

	//
	// STEP 10
	// Save genesis block to a database.
	//

	if err := s.createBlockDataUseCase.Execute(genesisBlockData.Hash, genesisBlockData.Header, genesisBlockData.HeaderSignatureBytes, genesisBlockData.Trans, genesisBlockData.Validator); err != nil {
		return fmt.Errorf("Failed to write genesis block data to file: %v", err)
	}
	if err := s.setBlockchainLastestHashUseCase.Execute(genesisBlockData.Hash); err != nil {
		return fmt.Errorf("Failed to save last hash of genesis block data: %v", err)
	}

	s.logger.Debug("genesis block created, finished running service",
		slog.String("hash", genesisBlockData.Hash))

	return nil
}
