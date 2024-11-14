package account

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
)

var (
	flagDataDirectory    string
	flagLabel            string
	flagPassword         string
	flagPasswordRepeated string
)

func NewAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Creates a new wallet in your local filesystem and encrypts it using the inputted password for security",
		Run: func(cmd *cobra.Command, args []string) {
			if err := doRunNewAccountCmd(); err != nil {
				log.Fatal(err)
			}
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().StringVar(&flagPassword, "wallet-password", "", "The password to encrypt the new wallet with")
	cmd.MarkFlagRequired("wallet-password")
	cmd.Flags().StringVar(&flagPasswordRepeated, "wallet-password-repeated", "", "The password repeated to verify your password is correct")
	cmd.MarkFlagRequired("wallet-password-repeated")
	cmd.Flags().StringVar(&flagLabel, "wallet-label", "", "The (optional) label to describe the new wallet with")

	return cmd
}

func doRunNewAccountCmd() error {
	logger := logger.NewProvider()

	// logger := logger.NewProvider()
	logger.Debug("Creating new account...",
		slog.Any("wallet_password", flagPassword),
		slog.Any("wallet_password_repeated", flagPasswordRepeated),
		slog.Any("wallet_label", flagLabel),
	)

	// ------ Common ------
	keystore := keystore.NewAdapter()
	walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)
	genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)
	blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)

	// ------ Repo ------
	walletRepo := repo.NewWalletRepo(
		logger,
		walletDB)
	accountRepo := repo.NewAccountRepo(
		logger,
		accountDB)
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		logger,
		genesisBlockDataDB)
	blockchainStateRepo := repo.NewBlockchainStateRepo(
		logger,
		blockchainStateDB)
	blockDataRepo := repo.NewBlockDataRepo(
		logger,
		blockDataDB)

	// ------ Use-case ------

	// Storage Transaction
	storageTransactionOpenUseCase := usecase.NewStorageTransactionOpenUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo)
	storageTransactionCommitUseCase := usecase.NewStorageTransactionCommitUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo)
	storageTransactionDiscardUseCase := usecase.NewStorageTransactionDiscardUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo)

	// Wallet
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		logger,
		keystore,
		walletRepo)
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		logger,
		keystore,
		walletRepo)
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		logger,
		walletRepo)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		logger,
		walletRepo)
	listAllWalletUseCase := usecase.NewListAllWalletUseCase(
		logger,
		walletRepo)

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		logger,
		accountRepo)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		logger,
		accountRepo)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		logger,
		accountRepo)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		logger,
		accountRepo)

	// ------ Service ------

	createAccountService := service.NewCreateAccountService(
		logger,
		walletEncryptKeyUseCase,
		walletDecryptKeyUseCase,
		createWalletUseCase,
		createAccountUseCase,
		getAccountUseCase,
	)

	_ = getAccountsHashStateUseCase
	_ = upsertAccountUseCase

	_ = getWalletUseCase
	_ = listAllWalletUseCase

	// ------ Execute ------
	ctx := context.Background()

	if err := storageTransactionOpenUseCase.Execute(); err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	account, err := createAccountService.Execute(ctx, flagPassword, flagPasswordRepeated, flagLabel)
	if err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to encrypt wallet: %v", err)
	}
	if account == nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatal("Account does not exist.")
	}

	if err := storageTransactionCommitUseCase.Execute(); err != nil {
		storageTransactionDiscardUseCase.Execute()
		log.Fatalf("Failed to open storage transaction: %v\n", err)
	}

	logger.Info("Account created",
		slog.Any("nonce", account.GetNonce()),
		slog.Uint64("balance", account.Balance),
		slog.String("address", account.Address.Hex()),
	)

	return nil
}
