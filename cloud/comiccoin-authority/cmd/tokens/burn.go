package tokens

import (
	"context"
	"log"
	"log/slog"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

func BurnTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "burn",
		Short: "Burns a non-fungible token you own",
		Run: func(cmd *cobra.Command, args []string) {
			doRunBurnToken()
		},
	}

	cmd.Flags().StringVar(&flagBurnTokenID, "token-id", "", "The ID of the token that you own")
	cmd.MarkFlagRequired("token-id")

	return cmd
}

func doRunBurnToken() {
	//
	// Load up dependencies.
	//

	// ------ Common ------
	logger := logger.NewProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter()
	kmutex := kmutexutil.NewKMutexProvider()

	// ------ Repository ------
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	// accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	blockchainStateRepo := repo.NewBlockchainStateRepo(cfg, logger, dbClient)
	tokRepo := repo.NewTokenRepo(cfg, logger, dbClient)
	// gbdRepo := repo.NewGenesisBlockDataRepo(cfg, logger, dbClient)
	bdRepo := repo.NewBlockDataRepo(cfg, logger, dbClient)
	mempoolTxRepo := repo.NewMempoolTransactionRepo(cfg, logger, dbClient)

	// // ------ Use-case ------
	// // Wallet
	// walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
	// 	cfg,
	// 	logger,
	// 	keystore,
	// 	walletRepo,
	// )
	// _ = walletEncryptKeyUseCase
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	// createWalletUseCase := usecase.NewCreateWalletUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo,
	// )
	// _ = createWalletUseCase
	getWalletUseCase := usecase.NewGetWalletUseCase(
		cfg,
		logger,
		walletRepo,
	)

	// // Account
	// createAccountUseCase := usecase.NewCreateAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	// _ = createAccountUseCase
	// getAccountUseCase := usecase.NewGetAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	// _ = getAccountUseCase
	// upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	// getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	//
	// Blockchain State
	getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
		cfg,
		logger,
		blockchainStateRepo,
	)
	upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
		cfg,
		logger,
		blockchainStateRepo,
	)

	// Token
	getTokenUseCase := usecase.NewGetTokenUseCase(
		cfg,
		logger,
		tokRepo,
	)
	// upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
	// 	cfg,
	// 	logger,
	// 	tokRepo,
	// )
	// getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
	// 	cfg,
	// 	logger,
	// 	tokRepo,
	// )
	//
	// // Genesis BlockData
	// upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
	// 	cfg,
	// 	logger,
	// 	gbdRepo,
	// )
	//
	// BlockData
	getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
		cfg,
		logger,
		bdRepo,
	)
	// upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
	// 	cfg,
	// 	logger,
	// 	bdRepo,
	// )
	//
	// // Proof of Work
	// proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(
	// 	cfg,
	// 	logger,
	// )

	// Mempool Transaction
	mempoolTransactionCreateUseCase := usecase.NewMempoolTransactionCreateUseCase(
		cfg,
		logger,
		mempoolTxRepo,
	)

	// ------ Service ------
	tokenBurnService := service.NewTokenBurnService(
		cfg,
		logger,
		kmutex,
		dbClient, // Note: Used for mongoDB transaction handling.
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getBlockchainStateUseCase,
		upsertBlockchainStateUseCase,
		getBlockDataUseCase,
		getTokenUseCase,
		mempoolTransactionCreateUseCase,
	)

	// Execution
	ctx := context.Background()
	burnRecipientAddress := common.HexToAddress(strings.ToLower(flagBurnRecipientAddress))
	tokenID, ok := new(big.Int).SetString(flagBurnTokenID, 10)
	if !ok {
		log.Fatalf("Failed converting to big.Int: %v", flagTokenID)
	}
	logger.Debug("burning token...",
		slog.Any("address", burnRecipientAddress),
		slog.Any("token_id", tokenID.Uint64()))

	err := tokenBurnService.Execute(
		ctx,
		tokenID,
		cfg.Blockchain.ProofOfAuthorityAccountAddress,
		cfg.Blockchain.ProofOfAuthorityWalletPassword,
	)
	if err != nil {
		logger.Error("Failed executing",
			slog.Any("error", err))
		log.Fatalf("Failed executing: %v\n", err)
	}
	logger.Debug("Token burned")
}
