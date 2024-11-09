package service

import (
	"context"
	"log/slog"
	"math"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type CreateGenesisBlockDataService struct {
	config                       *config.Configuration
	logger                       *slog.Logger
	createAccountService         *CreateAccountService
	getWalletUseCase             *usecase.GetWalletUseCase
	walletDecryptKeyUseCase      *usecase.WalletDecryptKeyUseCase
	upsertAccountUseCase         *usecase.UpsertAccountUseCase
	proofOfWorkUseCase           *usecase.ProofOfWorkUseCase
	upsertBlockchainStateUseCase *usecase.UpsertBlockchainStateUseCase
	getBlockchainStateUseCase    *usecase.GetBlockchainStateUseCase
}

func NewCreateGenesisBlockDataService(
	config *config.Configuration,
	logger *slog.Logger,
	s1 *CreateAccountService,
	uc1 *usecase.GetWalletUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
	uc3 *usecase.UpsertAccountUseCase,
	uc4 *usecase.ProofOfWorkUseCase,
	uc5 *usecase.UpsertBlockchainStateUseCase,
	uc6 *usecase.GetBlockchainStateUseCase,
) *CreateGenesisBlockDataService {
	return &CreateGenesisBlockDataService{config, logger, s1, uc1, uc2, uc3, uc4, uc5, uc6}
}

func (s *CreateGenesisBlockDataService) Execute(ctx context.Context, walletPassword string, walletPasswordRepeated string) (*domain.BlockchainState, error) {
	s.logger.Debug("starting genesis creation service...")

	//
	// STEP 1:
	// Validation.
	//

	e := make(map[string]string)
	if walletPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if walletPasswordRepeated == "" {
		e["wallet_password_repeated"] = "missing value"
	}
	if walletPassword != walletPasswordRepeated {
		e["wallet_password"] = "do not match"
		e["wallet_password_repeated"] = "do not match"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2:
	// Create our coinbase account and get the key.
	//

	account, err := s.createAccountService.Execute(ctx, walletPassword, walletPasswordRepeated, "Coinbase")
	if err != nil {
		s.logger.Error("Failed creating account", slog.Any("error", err))
		return nil, err
	}
	accountWallet, err := s.getWalletUseCase.Execute(ctx, account.Address)
	if err != nil {
		s.logger.Error("Failed creating wallet", slog.Any("error", err))
		return nil, err
	}
	coinbaseKey, err := s.walletDecryptKeyUseCase.Execute(ctx, accountWallet.FilePath, walletPassword)
	if err != nil {
		s.logger.Error("Failed decrypting wallet", slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 3:
	// Set coinbase with all the coins.
	//

	s.logger.Debug("starting genesis creation service...")

	// DEVELOPERS NOTE:
	// Here is where we initialize the total supply of coins for the entire
	// blockchain, so adjust accordingly. We will set the maximum value possible
	// for the unsigned 64-bit integer in a computer. That's a big number!
	initialSupply := uint64(math.MaxUint64) // Note: 18446744073709551615

	// DEVELOPERS NOTE:
	// The value of `18446744073709551615` using scientific notation is
	// 1.8446744 × 10^19 which can be read as approximately `1.8 quintillion`.
	//
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

	if err := s.upsertAccountUseCase.Execute(ctx, account.Address, initialSupply, 0); err != nil {
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
		Nonce:   0, // Will be calculated later.
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
		Nonce:            0, // Will be calculated later.
		From:             account.Address,
		To:               account.Address,
		Value:            0, //Note: Tokens don't have coin value.
		Tip:              0,
		Data:             make([]byte, 0),
		Type:             domain.TransactionTypeToken,
		TokenID:          0, // The very first token in our entire blockchain starts at the value of zero.
		TokenMetadataURI: "https://cpscapsule.com/comiccoin/tokens/0/metadata.json",
		TokenNonce:       0, // Newly minted tokens always have their nonce start at value of zero.
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
		slog.Uint64("tx_token_id", signedTokenTx.TokenID))

	// Defensive code: Run this code to ensure this transaction is
	// properly structured for our blockchain.
	if err := signedTokenTx.Validate(s.config.Blockchain.ChainID, true); err != nil {
		s.logger.Error("Failed token transaction.",
			slog.Any("error", err))
		return nil, err
	}

	_ = signedCoinTx

	return nil, nil

	//
	// STEP 6:
	// Save our token to our database.
	//

	//
	// if err := s.upsertTokenIfPreviousTokenNonceGTEUseCase.Execute(tokenTx.TokenID, tokenTx.To, tokenTx.TokenMetadataURI, tokenTx.TokenNonce); err != nil {
	// 	return err
	// }
	// if err := s.setBlockchainLastestTokenIDIfGTEUseCase.Execute(tokenTx.TokenID); err != nil {
	// 	return fmt.Errorf("Failed to save last token ID of genesis block data: %v", err)
	// }
	//
	// //
	// // STEP 5:
	// // Create our first block, i.e. also called "Genesis block".
	// //
	//
	// // Note: Genesis block has no previous hash
	// prevBlockHash := signature.ZeroHash
	//
	// gasPrice := uint64(s.config.Blockchain.GasPrice)
	// unitsOfGas := uint64(s.config.Blockchain.UnitsOfGas)
	// coinBlockTx := domain.BlockTransaction{
	// 	SignedTransaction: signedCoinTx,
	// 	TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
	// 	GasPrice:          gasPrice,
	// 	GasUnits:          unitsOfGas,
	// }
	// tokenBlockTx := domain.BlockTransaction{
	// 	SignedTransaction: signedTokenTx,
	// 	TimeStamp:         uint64(time.Now().UTC().UnixMilli()),
	// 	GasPrice:          gasPrice,
	// 	GasUnits:          unitsOfGas,
	// }
	// trans := make([]domain.BlockTransaction, 0)
	// trans = append(trans, coinBlockTx)
	// trans = append(trans, tokenBlockTx)
	//
	// // Construct a merkle tree from the transaction for this block. The root
	// // of this tree will be part of the block to be mined.
	// tree, err := merkle.NewTree(trans)
	// if err != nil {
	// 	return fmt.Errorf("Failed to create merkle tree: %v", err)
	// }
	//
	// stateRoot, err := s.getAccountsHashStateUseCase.Execute()
	// if err != nil {
	// 	s.logger.Error("Failed to get hash of all accounts",
	// 		slog.Any("error", err))
	// 	return fmt.Errorf("Failed to get hash of all accounts: %v", err)
	// }
	//
	// // Running this code get's a hash of all the tokens, thus making the
	// // tokens tamper proof.
	// tokensRoot, err := s.getTokensHashStateUseCase.Execute()
	// if err != nil {
	// 	s.logger.Error("Failed to get hash of all tokens",
	// 		slog.Any("error", err))
	// 	return fmt.Errorf("Failed to get hash of all tokens: %v", err)
	// }
	//
	// // Construct the genesis block.
	// block := domain.Block{
	// 	Header: &domain.BlockHeader{
	// 		Number:        0, // Genesis always starts at zero
	// 		PrevBlockHash: prevBlockHash,
	// 		TimeStamp:     uint64(time.Now().UTC().UnixMilli()),
	// 		Beneficiary:   s.coinbaseAccountKey.Address,
	// 		Difficulty:    s.config.Blockchain.Difficulty,
	// 		MiningReward:  s.config.Blockchain.MiningReward,
	// 		StateRoot:     stateRoot,
	// 		TransRoot:     tree.RootHex(), //
	// 		Nonce:         0,              // Will be identified by the POW algorithm.
	// 		LatestTokenID: 0,              // ComicCoin: Token ID values start at zero.
	// 		TokensRoot:    tokensRoot,
	// 	},
	// 	MerkleTree: tree,
	// }
	//
	// genesisBlockData := domain.NewBlockData(block)
	//
	// //
	// // STEP 6:
	// // Execute the proof of work to find our nounce to meet the hash difficulty.
	// //
	//
	// nonce, powErr := s.proofOfWorkUseCase.Execute(ctx, &block, s.config.Blockchain.Difficulty)
	// if powErr != nil {
	// 	return fmt.Errorf("Failed to mine block: %v", powErr)
	// }
	//
	// block.Header.Nonce = nonce
	//
	// s.logger.Debug("genesis mining completed",
	// 	slog.Uint64("nonce", block.Header.Nonce))
	//
	// // STEP 7:
	// // Create our single proof-of-authority validator via coinbase account.
	// //
	//
	// coinbasePrivateKey := s.coinbaseAccountKey.PrivateKey
	// // Extract the bytes for the original public key.
	// coinbasePublicKey := coinbasePrivateKey.Public()
	// publicKeyECDSA, ok := coinbasePublicKey.(*ecdsa.PublicKey)
	// if !ok {
	// 	return errors.New("error casting public key to ECDSA")
	// }
	// publicKeyBytes := crypto.FromECDSAPub(publicKeyECDSA)
	//
	// poaValidator := &domain.Validator{
	// 	ID:             "ComicCoin Blockchain Authority",
	// 	PublicKeyBytes: publicKeyBytes,
	// }
	//
	// //
	// // STEP 8:
	// // Sign our genesis block's header with our proof-of-authority validator.
	// // Note: Signing always happens after the miner has found the `nonce` in
	// // the block header.
	// //
	//
	// genesisBlockData.Validator = poaValidator
	// genesisBlockHeaderSignatureBytes, err := poaValidator.Sign(coinbasePrivateKey, genesisBlockData.Header)
	// if err != nil {
	// 	return fmt.Errorf("Failed to sign block header: %v", err)
	// }
	// genesisBlockData.HeaderSignatureBytes = genesisBlockHeaderSignatureBytes
	//
	// //
	// // STEP 9:
	// // Save genesis block to a JSON file.
	// //
	//
	// genesisBlockDataBytes, err := json.MarshalIndent(genesisBlockData, "", "    ")
	// if err != nil {
	// 	return fmt.Errorf("Failed to serialize genesis block: %v", err)
	// }
	//
	// if err := os.WriteFile("static/genesis.json", genesisBlockDataBytes, 0644); err != nil {
	// 	return fmt.Errorf("Failed to write genesis block data to file: %v", err)
	// }
	//
	// //
	// // STEP 10
	// // Save genesis block to a database.
	// //
	//
	// if err := s.createBlockDataUseCase.Execute(genesisBlockData.Hash, genesisBlockData.Header, genesisBlockData.HeaderSignatureBytes, genesisBlockData.Trans, genesisBlockData.Validator); err != nil {
	// 	return fmt.Errorf("Failed to write genesis block data to file: %v", err)
	// }
	// if err := s.setBlockchainLastestHashUseCase.Execute(genesisBlockData.Hash); err != nil {
	// 	return fmt.Errorf("Failed to save last hash of genesis block data: %v", err)
	// }
	//
	// s.logger.Debug("genesis block created, finished running service",
	// 	slog.String("hash", genesisBlockData.Hash))
	//
	// return nil
}
