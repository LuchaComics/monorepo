package tokens

import (
	"github.com/spf13/cobra"
)

func GetTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Get account details",
		Run: func(cmd *cobra.Command, args []string) {
			// doRunGetToken()
		},
	}

	cmd.Flags().StringVar(&flagTokenID, "token_id", "", "The token ID value to lookup the account by")
	cmd.MarkFlagRequired("address")

	return cmd
}

func doRunGetToken() {
	// // Common
	// logger := logger.NewProvider()
	// cfg := config.NewProvider()
	// dbClient := mongodb.NewProvider(cfg, logger)
	//
	// // Repository
	// accountRepo := repo.NewTokenRepo(cfg, logger, dbClient)
	//
	// // Use-case
	// getTokenUseCase := usecase.NewGetTokenUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	//
	// // Service
	// getTokenService := service.NewGetTokenService(
	// 	logger,
	// 	getTokenUseCase,
	// )
	//
	// ////
	// //// Start the transaction.
	// ////
	// ctx := context.Background()
	//
	// session, err := dbClient.StartSession()
	// if err != nil {
	// 	logger.Error("start session error",
	// 		slog.Any("error", err))
	// 	log.Fatalf("Failed executing: %v\n", err)
	// }
	// defer session.EndSession(ctx)
	//
	// logger.Debug("Creating new account...",
	// 	slog.Any("wallet_password", flagPassword),
	// 	slog.Any("wallet_password_repeated", flagPasswordRepeated),
	// 	slog.Any("wallet_label", flagLabel),
	// )
	//
	// // Define a transaction function with a series of operations
	// transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
	// 	accountAddress := common.HexToAddress(strings.ToLower(flagTokenAddress))
	//
	// 	account, err := getTokenService.Execute(sessCtx, &accountAddress)
	// 	if err != nil {
	// 		return nil, err
	// 	}
	// 	if account == nil {
	// 		err := errors.New("Token does not exist")
	// 		return nil, err
	// 	}
	//
	// 	return account, nil
	// }
	//
	// // Start a transaction
	// res, err := session.WithTransaction(ctx, transactionFunc)
	// if err != nil {
	// 	logger.Error("session failed error",
	// 		slog.Any("error", err))
	// 	log.Fatalf("Failed creating account: %v\n", err)
	// }
	//
	// account := res.(*domain.Token)
	//
	// logger.Debug("Token retrieved",
	// 	slog.Any("nonce", account.GetNonce()),
	// 	slog.Uint64("balance", account.Balance),
	// 	slog.String("address", account.Address.Hex()),
	// )
}
