package account

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

func ListAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "list",
		Short: "List all the accounts that belong to you that exist locally on this machine.",
		Run: func(cmd *cobra.Command, args []string) {
			doRunListAccount()
		},
	}

	return cmd
}

func doRunListAccount() {
	// ------ Common ------
	logger := logger.NewProvider()
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)
	walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)

	// ------ Repo ------
	accountRepo := repo.NewAccountRepo(
		logger,
		accountDB)
	walletRepo := repo.NewWalletRepo(
		logger,
		walletDB)

	// ------ Use-case ------
	listAllAddressesWalletUseCase := usecase.NewListAllAddressesWalletUseCase(
		logger,
		walletRepo,
	)
	accountsFilterByAddressesUseCase := usecase.NewAccountsFilterByAddressesUseCase(
		logger,
		accountRepo,
	)

	// ------ Service ------
	accountListingByLocalWalletsService := service.NewAccountListingByLocalWalletsService(
		logger,
		listAllAddressesWalletUseCase,
		accountsFilterByAddressesUseCase,
	)

	// ------ Execute ------
	ctx := context.Background()

	accounts, err := accountListingByLocalWalletsService.Execute(ctx)
	if err != nil {
		log.Fatalf("Failed to listing my account: %v\n", err)
	}

	for _, account := range accounts {
		logger.Debug("Local account retrieved",
			slog.Any("nonce", account.GetNonce()),
			slog.Uint64("balance", account.Balance),
			slog.String("address", account.Address.Hex()),
		)
	}
}
