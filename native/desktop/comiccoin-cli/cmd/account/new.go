package account

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	"github.com/spf13/cobra"

	// "github.com/LuchaComics/monorepo/cloud/comiccoin-authority-cli/config"
	// "github.com/LuchaComics/monorepo/cloud/comiccoin-authority-cli/service"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/repo"
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

	// cfg := &config.Config{}
	//
	// ------ Common ------

	walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)

	// ------ Repo ------
	walletRepo := repo.NewWalletRepo(
		logger,
		walletDB)
	accountRepo := repo.NewAccountRepo(
		logger,
		accountDB)

	// ------ Use-case ------

	// Wallet
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		walletRepo)
	_ = walletDecryptKeyUseCase
	// walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo)
	// createWalletUseCase := usecase.NewCreateWalletUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo)
	// getWalletUseCase := usecase.NewGetWalletUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo)
	// listAllWalletUseCase := usecase.NewListAllWalletUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo)
	//
	// // Account
	// createAccountUseCase := usecase.NewCreateAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo)
	// getAccountUseCase := usecase.NewGetAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo)
	// getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo)
	// upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo)
	//
	// // ------ Service ------
	//
	// createAccountService := service.NewCreateAccountService(
	// 	logger,
	// 	walletEncryptKeyUseCase,
	// 	walletDecryptKeyUseCase,
	// 	createWalletUseCase,
	// 	createAccountUseCase,
	// 	getAccountUseCase,
	// )
	//
	// _ = getAccountsHashStateUseCase
	// _ = upsertAccountUseCase
	//
	// _ = getWalletUseCase
	// _ = listAllWalletUseCase
	//
	// // Execute
	//
	// account, err := createAccountService.Execute(flagDataDirectory, flagPassword, flagPasswordRepeated, flagLabel)
	// if err != nil {
	// 	log.Fatalf("Failed to encrypt wallet: %v", err)
	// }
	// if account == nil {
	// 	log.Fatal("Account does not exist.")
	// }
	// logger.Debug("Account created",
	// 	slog.Uint64("nonce", account.Nonce),
	// 	slog.Uint64("balance", account.Balance),
	// 	slog.String("address", account.Address.Hex()),
	// )

	return nil
}
