package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
)

type ProofOfAuthorityTokenMintService struct {
	config                                *config.Config
	logger                                *slog.Logger
	kmutex                                kmutexutil.KMutexProvider
	loadGenesisBlockDataUseCase           *usecase.LoadGenesisBlockDataUseCase
	getWalletUseCase                      *usecase.GetWalletUseCase
	walletDecryptKeyUseCase               *usecase.WalletDecryptKeyUseCase
	getBlockchainLastestTokenIDUseCase    *usecase.GetBlockchainLastestTokenIDUseCase
	broadcastMempoolTransactionDTOUseCase *usecase.BroadcastMempoolTransactionDTOUseCase
}

func NewProofOfAuthorityTokenMintService(
	cfg *config.Config,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.LoadGenesisBlockDataUseCase,
	uc2 *usecase.GetWalletUseCase,
	uc3 *usecase.WalletDecryptKeyUseCase,
	uc4 *usecase.GetBlockchainLastestTokenIDUseCase,
	uc5 *usecase.BroadcastMempoolTransactionDTOUseCase,
) *ProofOfAuthorityTokenMintService {
	return &ProofOfAuthorityTokenMintService{cfg, logger, kmutex, uc1, uc2, uc3, uc4, uc5}
}

func (s *ProofOfAuthorityTokenMintService) Execute(
	ctx context.Context,
	poaAddr *common.Address,
	poaPassword string,
	toAddr *common.Address,
	metadataURI string,
) error {
	// Lock the mining service until it has completed executing (or errored).
	s.kmutex.Acquire("token-minting")
	defer s.kmutex.Release("token-minting")

	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if poaAddr == nil {
		e["poa_address"] = "missing value"
	}
	if poaPassword == "" {
		e["poa_password"] = "missing value"
	}
	if toAddr == nil {
		e["to"] = "missing value"
	}
	if metadataURI == "" {
		e["metadata_uri"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating token mint parameters",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get the account and extract the wallet private/public key.
	//

	wallet, err := s.getWalletUseCase.Execute(poaAddr)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed getting from database: %s", err)
	}
	if wallet == nil {
		s.logger.Error("failed getting from database",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", "d.n.e."))
		return fmt.Errorf("failed getting from database: %s", "wallet d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(wallet.Filepath, poaPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return fmt.Errorf("failed getting key: %s", "d.n.e.")
	}

	//
	// STEP 3:
	// Verify the account is validator.
	//

	genesisBlockData, err := s.loadGenesisBlockDataUseCase.Execute()
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed getting account: %s", err)
	}
	if genesisBlockData == nil {
		return fmt.Errorf("failed getting genesis block data: %s", "d.n.e.")
	}
	validator := genesisBlockData.Validator

	publicKeyECDSA, err := validator.GetPublicKeyECDSA()
	if err != nil {
		s.logger.Error("failed unmarshalling validator public key",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return fmt.Errorf("failed unmarshalling validator public key: %s", err)
	}

	// Developers Note:
	// This is a super important step to enforce the authority being used by
	// the correct party. This code verifies that the the public key of the
	// authority matches the public key set on the genesis block because the
	// user has opend up the actual authorities wallet.
	if key.PrivateKey.PublicKey.Equal(publicKeyECDSA) == false {
		return fmt.Errorf("failed comparing public keys: %s", "they do not match")
	}

	//
	// STEP 4
	// Authority generates the latest token ID value by taking the previous
	// token ID value and incrementing it by one.
	//

	latestTokenID, err := s.getBlockchainLastestTokenIDUseCase.Execute()
	if err != nil {
		s.logger.Error("failed getting latest token ID",
			slog.Any("from_account_address", poaAddr),
			slog.Any("error", err))
		return err
	}

	newTokenID := latestTokenID + 1

	//TODO: Add security feature of looking up the latest blockchain state and
	// compare the latest token id set there and the one here and error if
	// inconsistencies.

	//
	// STEP 5
	// Create our pending transaction and sign it with the accounts private key.
	//

	tx := &domain.Transaction{
		ChainID:          s.config.Blockchain.ChainID,
		Nonce:            uint64(time.Now().Unix()),
		From:             poaAddr,
		To:               toAddr,
		Value:            0, // Token have no value!
		Tip:              0,
		Data:             make([]byte, 0),
		Type:             domain.TransactionTypeToken,
		TokenID:          newTokenID,
		TokenMetadataURI: metadataURI,
		TokenNonce:       0, // Newly minted tokens always have their nonce start at value of zero.
	}

	stx, signingErr := tx.Sign(key.PrivateKey)
	if signingErr != nil {
		s.logger.Debug("Failed to sign the token mint transaction",
			slog.Any("error", signingErr))
		return signingErr
	}

	s.logger.Debug("Pending token mint transaction signed successfully",
		slog.Uint64("tx_token_id", stx.TokenID))

	//
	// STEP 6
	// Send our pending signed transaction to our distributed mempool nodes
	// in the blochcian network.
	//

	mempoolTx := &domain.MempoolTransaction{
		Transaction: stx.Transaction,
		V:           stx.V,
		R:           stx.R,
		S:           stx.S,
	}

	if err := s.broadcastMempoolTransactionDTOUseCase.Execute(ctx, mempoolTx); err != nil {
		s.logger.Error("Failed to broadcast to the blockchain network",
			slog.Any("error", err))
		return err
	}

	s.logger.Info("Pending signed token transaction submitted to blockchain",
		slog.Uint64("tx_token_id", stx.TokenID))

	//
	// STEP 7
	// Finish, however do not save the `token_id`, let our PoA validator
	// handle updating the database if the submission was correct. The reason
	// for this is because if we update the local database to early and the PoA
	// validator fails, then our local database will have incorrect latest
	// `token_id` values saved.
	//

	return nil
}
