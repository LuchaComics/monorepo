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

var (
	flagAccountAddress string
)

func GetAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get account details",
		Run: func(cmd *cobra.Command, args []string) {
			doRunGetAccount()
		},
	}

	cmd.Flags().StringVar(&flagAccountAddress, "address", "", "The address value to lookup the account by")
	cmd.MarkFlagRequired("address")

	return cmd
}

func doRunGetAccount() {
	logger := logger.NewProvider()

	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()

	accountAddress := common.HexToAddress(strings.ToLower(flagAccountAddress))

	account, err := rpcClient.GetAccount(ctx, &accountAddress)
	if err != nil {
		log.Fatalf("Failed to get account: %v\n", err)
	}

	if account == nil {
		err := errors.New("Account does not exist")
		log.Fatalf("Failed to get account: %v\n", err)
	}

	logger.Info("Account retrieved",
		slog.Any("nonce", account.GetNonce()),
		slog.Uint64("balance", account.Balance),
		slog.String("address", account.Address.Hex()),
	)
}
