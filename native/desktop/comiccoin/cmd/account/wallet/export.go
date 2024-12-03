package wallet

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
)

var (
	flagAccountAddress string
	flagExportFilepath string
)

func ExportWalletCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "export",
		Short: "Exports the wallet to specific filepath",
		Run: func(cmd *cobra.Command, args []string) {
			doRunExportWallet()
		},
	}

	cmd.Flags().StringVar(&flagAccountAddress, "address", "", "The address value to lookup the account by")
	cmd.MarkFlagRequired("address")
	cmd.Flags().StringVar(&flagExportFilepath, "filepath", "", "The location to save the wallet to")
	cmd.MarkFlagRequired("filepath")

	return cmd
}

func doRunExportWallet() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()

	accountAddress := common.HexToAddress(strings.ToLower(flagAccountAddress))

	if err := rpcClient.ExportWallet(ctx, &accountAddress, flagExportFilepath); err != nil {
		log.Fatalf("Failed to get account: %v\n", err)
	}

	logger.Info("Wallet exported",
		slog.Any("filepath", flagExportFilepath))
}
