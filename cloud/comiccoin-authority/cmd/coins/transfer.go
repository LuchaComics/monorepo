package coins

import (
	"context"
	"log"
	"log/slog"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
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
		Short: "Submit a (pending) transaction to the ComicCoin blockchain network to transfer coins from coinbase account to another account",
		Run: func(cmd *cobra.Command, args []string) {
			doRunTransferCoinsCommand()
		},
	}

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
	keystore := keystore.NewAdapter()

	// ------ Repository ------
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	mempoolTxRepo := repo.NewMempoolTransactionRepo(cfg, logger, dbClient)
	// blockchainStateRepo := repo.NewBlockchainStateRepo(cfg, logger, dbClient)
	// tokRepo := repo.NewTokenRepo(cfg, logger, dbClient)
	// gbdRepo := repo.NewGenesisBlockDataRepo(cfg, logger, dbClient)
	// bdRepo := repo.NewBlockDataRepo(cfg, logger, dbClient)

	// ------ Use-case ------
	// Wallet
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		cfg,
		logger,
		walletRepo,
	)

	// Account
	getAccountUseCase := usecase.NewGetAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)

	// Mempool Transaction
	mempoolTransactionCreateUseCase := usecase.NewMempoolTransactionCreateUseCase(
		cfg,
		logger,
		mempoolTxRepo,
	)

	// ------ Service ------
	coinTransferService := service.NewCoinTransferService(
		cfg,
		logger,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		mempoolTransactionCreateUseCase,
	)

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

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		recAddr := common.HexToAddress(strings.ToLower(flagRecipientAddress))

		// Execution
		err := coinTransferService.Execute(
			sessCtx,
			cfg.Blockchain.ProofOfAuthorityAccountAddress,
			cfg.Blockchain.ProofOfAuthorityWalletPassword,
			&recAddr,
			flagQuantity,
			[]byte(flagData),
		)
		if err != nil {
			logger.Error("Failed transfering coins",
				slog.Any("error", err))
			return nil, err
		}
		return nil, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		logger.Error("session failed error",
			slog.Any("error", err))
		log.Fatalf("Failed transfering coins: %v\n", err)
	}

	logger.Debug("Coins transfered",
		slog.Any("from", "coinbase"),
		slog.Any("to", cfg.Blockchain.ProofOfAuthorityAccountAddress),
		slog.Any("quantity", flagQuantity),
	)
}
