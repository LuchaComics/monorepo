package account

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
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
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()
	accounts, err := rpcClient.AccountListingByLocalWallets(ctx)
	if err != nil {
		log.Fatalf("Failed to get local accounts: %v\n", err)
	}
	for _, account := range accounts {
		logger.Debug("Local account retrieved",
			slog.Any("nonce", account.GetNonce()),
			slog.Uint64("balance", account.Balance),
			slog.String("address", account.Address.Hex()),
		)
	}
}
