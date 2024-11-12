package coins

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
)

// Command line argument flags
var (
	flagKeystoreFile                  string // Location of the wallet keystore
	flagDataDir                       string // Location of the database directory
	flagLabel                         string
	flagSenderAccountPassword         string
	flagSenderAccountPasswordRepeated string
	flagCoinbaseAddress               string
	flagRecipientAddress              string
	flagQuantity                      uint64
	flagKeypairName                   string
	flagSenderAccountAddress          string
	flagData                          string

	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string

	flagListenHTTPPort       int
	flagListenHTTPIP         string
	flagListenPeerToPeerPort int

	flagListenHTTPAddress string

	flagIdentityKeyID string
)

func TransferCoinsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "transfer",
		Short: "Submit a (pending) transaction to the ComicCoin blockchain network to transfer coins from your account to another account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunTransferCoinsCommand()
		},
	}

	cmd.Flags().StringVar(&flagSenderAccountAddress, "sender-account-address", "", "The address of the account we will use in our coin transfer")
	cmd.MarkFlagRequired("sender-account-address")

	cmd.Flags().StringVar(&flagSenderAccountPassword, "sender-account-password", "", "The password to unlock the account which will transfer the coin")
	cmd.MarkFlagRequired("sender-account-password")

	cmd.Flags().Uint64Var(&flagQuantity, "value", 0, "The amount of coins to send")
	cmd.MarkFlagRequired("value")

	cmd.Flags().StringVar(&flagData, "data", "", "Optional data to include with this transaction")

	cmd.Flags().StringVar(&flagRecipientAddress, "recipient-address", "", "The address of the account whom will receive this coin")
	cmd.MarkFlagRequired("recipient-address")

	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "The HTTP JSON API server's port")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "The HTTP JSON API server's ip-address")

	return cmd
}

func doRunTransferCoinsCommand() {
	//
	// Load up dependencies.
	//

	// ------ Common ------
	logger := logger.NewProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter(cfg, logger)

	// ------ Repository ------
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	// blockchainStateRepo := repo.NewBlockchainStateRepo(cfg, logger, dbClient)
	// tokRepo := repo.NewTokenRepo(cfg, logger, dbClient)
	// gbdRepo := repo.NewGenesisBlockDataRepo(cfg, logger, dbClient)
	// bdRepo := repo.NewBlockDataRepo(cfg, logger, dbClient)

	_ = keystore
	_ = walletRepo
	_ = accountRepo

	// // ------ Use-case ------
	// // Wallet
	// walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
	// 	cfg,
	// 	logger,
	// 	keystore,
	// 	walletRepo,
	// )
	// walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
	// 	cfg,
	// 	logger,
	// 	keystore,
	// 	walletRepo,
	// )
	// createWalletUseCase := usecase.NewCreateWalletUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo,
	// )
	// getWalletUseCase := usecase.NewGetWalletUseCase(
	// 	cfg,
	// 	logger,
	// 	walletRepo,
	// )
	//
	// // Account
	// createAccountUseCase := usecase.NewCreateAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	// getAccountUseCase := usecase.NewGetAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	// upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )
	// getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
	// 	cfg,
	// 	logger,
	// 	accountRepo,
	// )

	// // Blockchain State
	// getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
	// 	cfg,
	// 	logger,
	// 	blockchainStateRepo,
	// )
	// upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
	// 	cfg,
	// 	logger,
	// 	blockchainStateRepo,
	// )
	//
	// // Token
	// upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
	// 	cfg,
	// 	logger,
	// 	tokRepo,
	// )
	// getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
	// 	cfg,
	// 	logger,
	// 	tokRepo,
	// )
	//
	// // Genesis BlockData
	// upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
	// 	cfg,
	// 	logger,
	// 	gbdRepo,
	// )
	//
	// // BlockData
	// upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
	// 	cfg,
	// 	logger,
	// 	bdRepo,
	// )
	//
	// // Proof of Work
	// proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(
	// 	cfg,
	// 	logger,
	// )

	// ------ Service ------
	// createAccountService := service.NewCreateAccountService(
	// 	cfg,
	// 	logger,
	// 	walletEncryptKeyUseCase,
	// 	walletDecryptKeyUseCase,
	// 	createWalletUseCase,
	// 	createAccountUseCase,
	// 	getAccountUseCase,
	// )

	////
	//// Start the transaction.
	////
	ctx := context.Background()

	session, err := dbClient.StartSession()
	if err != nil {
		logger.Error("start session error",
			slog.Any("error", err))
		log.Fatalf("Failed executing: %v\n", err)
	}
	defer session.EndSession(ctx)

	logger.Debug("Creating new account...",
		slog.Any("wallet_password", flagSenderAccountPassword),
		slog.Any("wallet_password_repeated", flagSenderAccountPasswordRepeated),
		slog.Any("wallet_label", flagLabel),
	)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// // Execution
		// blockchainState, err := createGenesisBlockDataService.Execute(sessCtx, flagPassword, flagPasswordRepeated)
		// if err != nil {
		// 	logger.Error("Failed initializing new blockchain from new genesis block",
		// 		slog.Any("error", err))
		// 	return nil, err
		// }
		// if blockchainState == nil {
		// 	err := errors.New("Blockchain state does not exist")
		// 	return nil, err
		// }
		//
		// return blockchainState, nil
		return nil, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		logger.Error("session failed error",
			slog.Any("error", err))
		log.Fatalf("Failed creating account: %v\n", err)
	}

	blockchainState := res.(*domain.BlockchainState)

	logger.Debug("Coins transfered",
		slog.Any("chain_id", blockchainState.ChainID),
	)

}
