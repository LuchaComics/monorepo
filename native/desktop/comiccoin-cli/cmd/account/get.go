package account

import (
	"context"
	"errors"
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/usecase"
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
	// ------ Common ------
	logger := logger.NewProvider()
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)

	// ------ Repo ------
	accountRepo := repo.NewAccountRepo(
		logger,
		accountDB)

	// ------ Use-case ------
	getAccountUseCase := usecase.NewGetAccountUseCase(
		logger,
		accountRepo,
	)

	// ------ Service ------
	getAccountService := service.NewGetAccountService(
		logger,
		getAccountUseCase,
	)

	// ------ Execute ------
	ctx := context.Background()

	accountAddress := common.HexToAddress(strings.ToLower(flagAccountAddress))

	account, err := getAccountService.Execute(ctx, &accountAddress)
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
