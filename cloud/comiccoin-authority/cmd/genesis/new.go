package genesis

import (
	"context"
	"errors"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

func NewGenesistCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Initializes a new blockchain by creating a genesis block",
		Run: func(cmd *cobra.Command, args []string) {
			doRunNewAccount()
		},
	}

	cmd.Flags().StringVar(&flagPassword, "coinbase-password", "", "The password to encrypt the coinbase's account wallet")
	cmd.MarkFlagRequired("coinbase-password")
	cmd.Flags().StringVar(&flagPasswordRepeated, "coinbase-password-repeated", "", "The password (again) to verify you are entering the correct password")
	cmd.MarkFlagRequired("coinbase-password-repeated")

	return cmd
}

func doRunNewAccount() {
	// Common
	logger := logger.NewProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter(cfg, logger)

	// Repository
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	blockchainStateRepo := repo.NewBlockchainStateRepo(cfg, logger, dbClient)

	// Use-case
	// Wallet
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		cfg,
		logger,
		walletRepo,
	)

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)

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

	// Proof of Work
	proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(
		cfg,
		logger,
	)

	// Service
	createAccountService := service.NewCreateAccountService(
		cfg,
		logger,
		walletEncryptKeyUseCase,
		walletDecryptKeyUseCase,
		createWalletUseCase,
		createAccountUseCase,
		getAccountUseCase,
	)
	createGenesisBlockDataService := service.NewCreateGenesisBlockDataService(
		cfg,
		logger,
		createAccountService,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		proofOfWorkUseCase,
		upsertBlockchainStateUseCase,
		getBlockchainStateUseCase,
	)

	////
	//// Start the transaction.
	////
	ctx := context.Background()

	session, err := dbClient.StartSession()
	if err != nil {
		logger.Error("start session error",
			slog.Any("error", err))
		log.Fatalf("Failed executing: %v\n", err)
	}
	defer session.EndSession(ctx)

	logger.Debug("Creating new account...",
		slog.Any("wallet_password", flagPassword),
		slog.Any("wallet_password_repeated", flagPasswordRepeated),
		slog.Any("wallet_label", flagLabel),
	)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Execution
		blockchainState, err := createGenesisBlockDataService.Execute(sessCtx, flagPassword, flagPasswordRepeated)
		if err != nil {
			logger.Error("Failed initializing new blockchain from new genesis block",
				slog.Any("error", err))
			return nil, err
		}
		if blockchainState == nil {
			err := errors.New("Blockchain state does not exist")
			return nil, err
		}

		return blockchainState, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		logger.Error("session failed error",
			slog.Any("error", err))
		log.Fatalf("Failed creating account: %v\n", err)
	}

	blockchainState := res.(*domain.BlockchainState)

	logger.Debug("Genesis block created",
		slog.Any("chain_id", blockchainState.ChainID),
		// slog.Uint64("balance", account.Balance),
		// slog.String("address", account.Address.Hex()),
	)
}
