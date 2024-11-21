package account

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
)

func ListBlockTransactionsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "transactions",
		Short: "Get all transactions for an account address",
		Run: func(cmd *cobra.Command, args []string) {
			doRunListBlockTransactions()
		},
	}

	cmd.Flags().StringVar(&flagAccountAddress, "address", "", "The address value to lookup the account transactions by")
	cmd.MarkFlagRequired("address")

	return cmd
}

func doRunListBlockTransactions() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()

	accountAddress := common.HexToAddress(strings.ToLower(flagAccountAddress))

	txs, err := rpcClient.ListBlockTransactionsByAddress(ctx, &accountAddress)
	if err != nil {
		log.Fatalf("Failed to get account: %v\n", err)
	}

	if txs == nil {
		err := errors.New("No transactions exist for the account address")
		log.Fatalf("Failed to get block transactions: %v\n", err)
	}

	logger.Info("Block transactions retrieved")

	for _, tx := range txs {
		logger.Info("Transaction Detail",
			slog.Any("From", tx.From),
			slog.Any("To", tx.To),
			slog.Any("Value", tx.Value),
			slog.Any("Data", tx.Data),
			slog.Any("TokenID", tx.GetTokenID()),
			slog.Any("TokenMetadataURI", tx.TokenMetadataURI),
		)
	}

}
