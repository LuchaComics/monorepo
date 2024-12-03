package wallet

import (
	"context"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
)

var (
	flagImportWalletFilepath string
)

func ImportWalletCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "import",
		Short: "Imports the wallet to specific filepath",
		Run: func(cmd *cobra.Command, args []string) {
			doRunImportWallet()
		},
	}

	cmd.Flags().StringVar(&flagImportWalletFilepath, "filepath", "", "The location of the wallet")
	cmd.MarkFlagRequired("filepath")

	return cmd
}

func doRunImportWallet() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()

	if err := rpcClient.ImportWallet(ctx, flagImportWalletFilepath); err != nil {
		log.Fatalf("Failed to get account: %v\n", err)
	}

	logger.Info("Wallet imported",
		slog.Any("filepath", flagImportWalletFilepath))
}
