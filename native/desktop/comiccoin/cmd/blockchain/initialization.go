package blockchain

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/memory"
)

func InitCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "(Developer only) Initialize the ComicCoin blockchain for the very first time by creating the genesis block and the coinbase account.",
		Run: func(cmd *cobra.Command, args []string) {
			doRunInitBlockchain()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagPassword, "coinbase-password", "", "The password to encrypt the cointbase's account wallet")
	cmd.MarkFlagRequired("coinbase-password")

	return cmd
}

func doRunInitBlockchain() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

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
	logger := logger.NewLogger()
	walletDB := disk.NewDiskStorage(cfg.DB.DataDir+"/wallet", logger)
	blockDataDB := disk.NewDiskStorage(cfg.DB.DataDir+"/block_data", logger)
	latestHashDB := disk.NewDiskStorage(cfg.DB.DataDir+"/latest_hash", logger)
	// db := disk.NewDiskStorage(cfg.DB.DataDir, logger)
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

	account, err := createAccountService.Execute(flagDataDir, flagPassword)
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
		setBlockchainLastestHashUseCase,
		createBlockDataUseCase,
		proofOfWorkUseCase,
		upsertAccountUseCase,
	)

	ctx := context.Background()
	if err := createGenesisBlockDataService.Execute(ctx); err != nil {
		log.Fatalf("failed creating genesis blockdata: %v", err)
	}

	logger.Info("Blockchain successfully initialized",
		slog.Any("coinbase_address", account.Address),
	)
}
