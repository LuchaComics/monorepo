package blockchain

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/memory"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

func InitCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "(Developer only) Initialize the ComicCoin blockchain for the very first time by creating the genesis block and the coinbase account.",
		Run: func(cmd *cobra.Command, args []string) {
			doRunInitBlockchain()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", config.GetDefaultDataDirectory(), "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagPassword, "coinbase-password", "", "The password to encrypt the coinbase's account wallet")
	cmd.MarkFlagRequired("coinbase-password")
	cmd.Flags().StringVar(&flagPasswordRepeated, "coinbase-password-repeated", "", "The password (again) to verify you are entering the correct password")
	cmd.MarkFlagRequired("coinbase-password-repeated")

	return cmd
}

func doRunInitBlockchain() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	logger := logger.NewProvider()
	logger.Debug("Excuting...",
		slog.String("data_dir", flagDataDir))
	if flagDataDir == "./data" {
		log.Fatal("cannot be `./data`")
	}

	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID:                        constants.ChainIDMainNet,
			TransPerBlock:                  1,
			Difficulty:                     2,
			ConsensusPollingDelayInMinutes: 1,
		},
		App: config.AppConfig{
			DirPath:     flagDataDir,
			HTTPAddress: flagListenHTTPAddress,
		},
		DB: config.DBConfig{
			DataDir: flagDataDir,
		},
		Peer: config.PeerConfig{
			ListenPort: flagListenPeerToPeerPort,
			// KeyName:          flagKeypairName,
			RendezvousString: flagRendezvousString,
		},
	}
	walletDB := disk.NewDiskStorage(cfg.DB.DataDir, "wallet", logger)
	blockDataDB := disk.NewDiskStorage(cfg.DB.DataDir, "block_data", logger)
	latestHashDB := disk.NewDiskStorage(cfg.DB.DataDir, "latest_hash", logger)
	latestTokenIDDB := disk.NewDiskStorage(cfg.DB.DataDir, "latest_token_id", logger)
	tokenDB := disk.NewDiskStorage(cfg.DB.DataDir, "token", logger)
	memdb := memory.NewInMemoryStorage(logger)

	// ------------ Repo ------------
	walletRepo := repo.NewWalletRepo(
		cfg,
		logger,
		walletDB)
	latestBlockDataHashRepo := repo.NewBlockchainLastestHashRepo(
		cfg,
		logger,
		latestHashDB)
	blockDataRepo := repo.NewBlockDataRepo(
		cfg,
		logger,
		blockDataDB)
	accountRepo := repo.NewAccountRepo(
		cfg,
		logger,
		memdb) // Do not store on disk, only in-memory.
	latestBlockDataTokenIDRepo := repo.NewBlockchainLastestTokenIDRepo(
		cfg,
		logger,
		latestTokenIDDB)
	tokenRepo := repo.NewTokenRepo(
		cfg,
		logger,
		tokenDB)

	// ------------ Use-case ------------

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		cfg,
		logger,
		accountRepo)

	// Token
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokenRepo)
	getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
		cfg,
		logger,
		tokenRepo)

	// Wallet
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		cfg,
		logger,
		walletRepo)
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		walletRepo)
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		cfg,
		logger,
		walletRepo)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		cfg,
		logger,
		walletRepo)

	setBlockchainLastestHashUseCase := usecase.NewSetBlockchainLastestHashUseCase(
		cfg,
		logger,
		latestBlockDataHashRepo)
	createBlockDataUseCase := usecase.NewCreateBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo)
	proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(cfg, logger)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		cfg,
		logger,
		accountRepo)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo)
	setBlockchainLastestTokenIDIfGreatestUseCase := usecase.NewSetBlockchainLastestTokenIDIfGreatestUseCase(
		cfg,
		logger,
		latestBlockDataTokenIDRepo)

	// ------------ Service ------------

	createAccountService := service.NewCreateAccountService(
		cfg,
		logger,
		walletEncryptKeyUseCase,
		walletDecryptKeyUseCase,
		createWalletUseCase,
		createAccountUseCase,
		getAccountUseCase)

	getKeyService := service.NewGetKeyService(
		cfg,
		logger,
		getWalletUseCase,
		walletDecryptKeyUseCase)

	//
	// STEP 2:
	// Create our coinbase account.
	//

	account, err := createAccountService.Execute(flagDataDir, flagPassword, flagPasswordRepeated, flagLabel)
	if err != nil {
		log.Fatalf("failed creating account: %v", err)
	}

	coinbaseAccountKey, err := getKeyService.Execute(account.Address, flagPassword)
	if err != nil {
		log.Fatalf("failed getting account wallet key: %v", err)
	}

	// DEVELOPERS NOTE:
	// Since we are using in-memory database, we'll need to manually create
	// the coinbase account before proceeding. This is not a mistake, remember
	// in-memory data get's lost on app shutdown, so when the app starts up
	// again you'll need to populate the accounts database again.
	if err := upsertAccountUseCase.Execute(account.Address, 0, 0); err != nil {
		log.Fatalf("failed upserting account: %v", err)
	}

	//
	// STEP 3:
	// Execute our genesis creation.
	//

	createGenesisBlockDataService := service.NewCreateGenesisBlockDataService(
		cfg,
		logger,
		coinbaseAccountKey,
		getAccountsHashStateUseCase,
		getTokensHashStateUseCase,
		setBlockchainLastestHashUseCase,
		setBlockchainLastestTokenIDIfGreatestUseCase,
		createBlockDataUseCase,
		proofOfWorkUseCase,
		upsertAccountUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
	)

	ctx := context.Background()
	if err := createGenesisBlockDataService.Execute(ctx); err != nil {
		log.Fatalf("failed creating genesis blockdata: %v", err)
	}

	logger.Info("Blockchain successfully initialized",
		slog.Any("coinbase_address", account.Address),
	)
}
