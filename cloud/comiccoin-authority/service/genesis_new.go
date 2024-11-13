package service

import (
	"crypto/ecdsa"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/crypto"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/merkle"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/signature"

	// "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type CreateGenesisBlockDataService struct {
	config                                    *config.Configuration
	logger                                    *slog.Logger
	getProofOfAuthorityPrivateKeyService      *GetProofOfAuthorityPrivateKeyService
	getAccountUseCase                         *usecase.GetAccountUseCase
	upsertAccountUseCase                      *usecase.UpsertAccountUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase
	getAccountsHashStateUseCase               *usecase.GetAccountsHashStateUseCase
	getTokensHashStateUseCase                 *usecase.GetTokensHashStateUseCase
	proofOfWorkUseCase                        *usecase.ProofOfWorkUseCase
	upsertGenesisBlockDataUseCase             *usecase.UpsertGenesisBlockDataUseCase
	upsertBlockDataUseCase                    *usecase.UpsertBlockDataUseCase
	upsertBlockchainStateUseCase              *usecase.UpsertBlockchainStateUseCase
	getBlockchainStateUseCase                 *usecase.GetBlockchainStateUseCase
}

func NewCreateGenesisBlockDataService(
	config *config.Configuration,
	logger *slog.Logger,
	s1 *GetProofOfAuthorityPrivateKeyService,
	uc1 *usecase.GetAccountUseCase,
	uc2 *usecase.UpsertAccountUseCase,
	uc3 *usecase.UpsertTokenIfPreviousTokenNonceGTEUseCase,
	uc4 *usecase.GetAccountsHashStateUseCase,
	uc5 *usecase.GetTokensHashStateUseCase,
	uc6 *usecase.ProofOfWorkUseCase,
	uc7 *usecase.UpsertGenesisBlockDataUseCase,
	uc8 *usecase.UpsertBlockDataUseCase,
	uc9 *usecase.UpsertBlockchainStateUseCase,
	uc10 *usecase.GetBlockchainStateUseCase,
) *CreateGenesisBlockDataService {
	return &CreateGenesisBlockDataService{config, logger, s1, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9, uc10}
}

func (s *CreateGenesisBlockDataService) Execute(sessCtx mongo.SessionContext) (*domain.BlockchainState, error) {
	s.logger.Debug("starting genesis creation service...")

	//
	// STEP 1:
	// Get our coinbase account and private key.
	//

	coinbaseKey, err := s.getProofOfAuthorityPrivateKeyService.Execute(sessCtx)
	if err != nil {
		s.logger.Error("Failed getting proof of authority private key", slog.Any("error", err))
		return nil, err
	}
	account, err := s.getAccountUseCase.Execute(sessCtx, s.config.Blockchain.ProofOfAuthorityAccountAddress)
	if err != nil {
		s.logger.Error("Failed getting proof of authority account", slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 2:
	// Set coinbase with all the coins.
	//

	s.logger.Debug("starting genesis creation service...")

	// DEVELOPERS NOTE:
	// Here is where we initialize the total supply of coins for the entire
	// blockchain, so adjust accordingly. We will set the maximum value possible
	// for the unsigned 64-bit integer in a computer. That's a big number!
	initialSupply := uint64(math.MaxInt64) // Note: 9223372036854775807

	// DEVELOPERS NOTE:
	// Also here are some additional notes on the order of magnitude for powers
	// of 10:
	// 10^0 = 1
	// 10^3 = thousand
	// 10^6 = million
	// 10^9 = billion
	// 10^12 = trillion
	// 10^15 = quadrillion
	// 10^18 = quintillion
	// 10^21 = sextillion
	// 10^24 = septillion

	//
	// STEP 3
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

	if err := s.upsertAccountUseCase.Execute(sessCtx, account.Address, initialSupply, big.NewInt(0)); err != nil {
		s.logger.Error("Failed upserting account", slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 4:
	// Setup our very first (signed) transaction: i.e. coinbase giving coins
	// onto the blockchain ... from nothing.
	//

	coinTx := &domain.Transaction{
		ChainID: s.config.Blockchain.ChainID,
		Nonce:   big.NewInt(0).Bytes(), // Will be calculated later.
		From:    account.Address,
		To:      account.Address,
		Value:   initialSupply,
		Tip:     0,
		Data:    make([]byte, 0),
		Type:    domain.TransactionTypeCoin,
	}
	signedCoinTx, err := coinTx.Sign(coinbaseKey.PrivateKey)
	if err != nil {
		s.logger.Error("Failed signing coin transaction", slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 5:
	// Setup our very first (signed) token transaction.
	//

	tokenTx := &domain.Transaction{
		ChainID:          s.config.Blockchain.ChainID,
		Nonce:            big.NewInt(0).Bytes(), // Will be calculated later.
		From:             account.Address,
		To:               account.Address,
		Value:            0, //Note: Tokens don't have coin value.
		Tip:              0,
		Data:             make([]byte, 0),
		Type:             domain.TransactionTypeToken,
		TokenID:          big.NewInt(0).Bytes(), // The very first token in our entire blockchain starts at the value of zero.
		TokenMetadataURI: "https://cpscapsule.com/comiccoin/tokens/0/metadata.json",
		TokenNonce:       big.NewInt(0).Bytes(), // Newly minted tokens always have their nonce start at value of zero.
	}
	signedTokenTx, err := tokenTx.Sign(coinbaseKey.PrivateKey)
	if err != nil {
		s.logger.Error("Failed signing token transaction", slog.Any("error", err))
		return nil, err
	}

	nftFromAddr, err := signedTokenTx.FromAddress()
	if err != nil {
		s.logger.Error("Failed getting from address",
			slog.Any("chain_id", s.config.Blockchain.ChainID),
			slog.Any("error", err))
		return nil, err
	}

	s.logger.Info("Created first token",
		slog.Any("from", signedTokenTx.From),
		slog.Any("from_via_sig", nftFromAddr),
		slog.Any("to", signedTokenTx.To),
		slog.Any("tx_sig_v", signedTokenTx.V),
		slog.Any("tx_sig_r", signedTokenTx.R),
		slog.Any("tx_sig_s", signedTokenTx.S),
		slog.Any("tx_token_id", signedTokenTx.GetTokenID()))

	// Defensive code: Run this code to ensure this transaction is
	// properly structured for our blockchain.
	if err := signedTokenTx.Validate(s.config.Blockchain.ChainID, true); err != nil {
		s.logger.Error("Failed token transaction.",
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 6:
	// Save our token to our database.
	//

	if err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(sessCtx, tokenTx.GetTokenID(), tokenTx.To, tokenTx.TokenMetadataURI, tokenTx.GetTokenNonce()); err != nil {
		return nil, err
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
		return nil, fmt.Errorf("Failed to create merkle tree: %v", err)
	}

	stateRoot, err := s.getAccountsHashStateUseCase.Execute(sessCtx)
	if err != nil {
		s.logger.Error("Failed to get hash of all accounts",
			slog.Any("error", err))
		return nil, fmt.Errorf("Failed to get hash of all accounts: %v", err)
	}

	// Running this code get's a hash of all the tokens, thus making the
	// tokens tamper proof.
	tokensRoot, err := s.getTokensHashStateUseCase.Execute(sessCtx)
	if err != nil {
		s.logger.Error("Failed to get hash of all tokens",
			slog.Any("error", err))
		return nil, fmt.Errorf("Failed to get hash of all tokens: %v", err)
	}

	// Construct the genesis block.
	block := domain.Block{
		Header: &domain.BlockHeader{
			ChainID:       uint16(s.config.Blockchain.ChainID),
			Number:        big.NewInt(0).Bytes(), // Genesis always starts at zero
			PrevBlockHash: prevBlockHash,
			TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
			Beneficiary:   *account.Address,
			Difficulty:    s.config.Blockchain.Difficulty,
			MiningReward:  s.config.Blockchain.MiningReward,
			StateRoot:     stateRoot,
			TransRoot:     tree.RootHex(),        //
			Nonce:         big.NewInt(0).Bytes(), // Will be identified by the POW algorithm.
			LatestTokenID: big.NewInt(0).Bytes(), // ComicCoin: Token ID values start at zero.
			TokensRoot:    tokensRoot,
		},
		MerkleTree: tree,
	}

	genesisBlockData := domain.NewBlockData(block)

	//
	// STEP 6:
	// Execute the proof of work to find our nounce to meet the hash difficulty.
	//

	nonce, powErr := s.proofOfWorkUseCase.Execute(sessCtx, &block, s.config.Blockchain.Difficulty)
	if powErr != nil {
		return nil, fmt.Errorf("Failed to mine block: %v", powErr)
	}

	block.Header.Nonce = nonce.Bytes()

	s.logger.Debug("genesis mining completed",
		slog.Any("nonce", block.Header.GetNonce()))

	// STEP 7:
	// Create our single proof-of-authority validator via coinbase account.
	//

	coinbasePrivateKey := coinbaseKey.PrivateKey
	// Extract the bytes for the original public key.
	coinbasePublicKey := coinbasePrivateKey.Public()
	publicKeyECDSA, ok := coinbasePublicKey.(*ecdsa.PublicKey)
	if !ok {
		return nil, errors.New("error casting public key to ECDSA")
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
		return nil, fmt.Errorf("Failed to sign block header: %v", err)
	}
	genesisBlockData.HeaderSignatureBytes = genesisBlockHeaderSignatureBytes

	// //
	// // STEP 9:
	// // Save genesis block to a JSON file.
	// //
	//
	// genesisBlockDataBytes, err := json.MarshalIndent(genesisBlockData, "", "    ")
	// if err != nil {
	// 	return nil, fmt.Errorf("Failed to serialize genesis block: %v", err)
	// }
	//
	// if err := os.WriteFile("static/genesis.json", genesisBlockDataBytes, 0644); err != nil {
	// 	return nil, fmt.Errorf("Failed to write genesis block data to file: %v", err)
	// }

	//
	// STEP 10
	// Save genesis block to a database.
	//

	if err := s.upsertBlockDataUseCase.Execute(sessCtx, genesisBlockData.Hash, genesisBlockData.Header, genesisBlockData.HeaderSignatureBytes, genesisBlockData.Trans, genesisBlockData.Validator); err != nil {
		return nil, fmt.Errorf("Failed to write genesis block data to file: %v", err)
	}

	if err := s.upsertGenesisBlockDataUseCase.Execute(sessCtx, genesisBlockData.Hash, genesisBlockData.Header, genesisBlockData.HeaderSignatureBytes, genesisBlockData.Trans, genesisBlockData.Validator); err != nil {
		return nil, fmt.Errorf("Failed to write genesis block data to file: %v", err)
	}

	s.logger.Debug("genesis block created, finished running service",
		slog.String("hash", genesisBlockData.Hash))

	//
	// STEP 11:
	// Save our blockchain state.
	//

	blockchainState := &domain.BlockchainState{
		ChainID:           s.config.Blockchain.ChainID,
		LatestBlockNumber: genesisBlockData.Header.Number,
		LatestHash:        genesisBlockData.Hash,
		LatestTokenID:     tokenTx.TokenID,
		AccountHashState:  stateRoot,
		TokenHashState:    tokensRoot,
	}

	if err := s.upsertBlockchainStateUseCase.Execute(sessCtx, blockchainState); err != nil {
		s.logger.Error("Failed to save blockchain state",
			slog.Any("error", err))
		return nil, fmt.Errorf("Failed to save blockchain state: %v", err)
	}

	return blockchainState, nil
}
