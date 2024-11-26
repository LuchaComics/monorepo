package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/blacklist"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http"
	httphandler "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/task"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run the ComicCoin Faucet",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Running daemon......")
			doRunDaemon()
		},
	}
	return cmd
}

func doRunDaemon() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	// Common
	logger := logger.NewProvider()
	// kmutex := kmutexutil.NewKMutexProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	// keystore := keystore.NewAdapter()
	// passp := password.NewProvider()
	blackp := blacklist.NewProvider()

	//
	// Repository
	//

	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	_ = walletRepo
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	tenantRepo := repo.NewTenantRepository(cfg, logger, dbClient)
	_ = tenantRepo
	userRepo := repo.NewUserRepository(cfg, logger, dbClient)
	_ = userRepo
	tokRepo := repo.NewTokenRepository(cfg, logger, dbClient)
	blockchainStateRepo := repo.NewBlockchainStateRepository(cfg, logger, dbClient)
	_ = blockchainStateRepo

	blockchainStateChangeEventDTOConfigurationProvider := repo.NewBlockchainStateChangeEventDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockchainStateChangeEventDTORepo := repo.NewBlockchainStateChangeEventDTORepo(
		blockchainStateChangeEventDTOConfigurationProvider,
		logger)

	blockchainStateDTOConfigurationProvider := repo.NewBlockchainStateDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockchainStateDTORepo := repo.NewBlockchainStateDTORepo(
		blockchainStateDTOConfigurationProvider,
		logger)
	_ = blockchainStateDTORepo

	blockDataRepo := repo.NewBlockDataRepository(
		cfg,
		logger,
		dbClient)
	_ = blockDataRepo

	blockDataDTOConfigurationProvider := repo.NewBlockDataDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockDataDTORepo := repo.NewBlockDataDTORepository(
		blockDataDTOConfigurationProvider,
		logger)
	_ = blockDataDTORepo

	genesisBlockDataRepo := repo.NewGenesisBlockDataRepository(
		cfg,
		logger,
		dbClient)
	_ = genesisBlockDataRepo

	genesisBlockDataDTOConfigurationProvider := repo.NewGenesisBlockDataDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	genesisBlockDataDTORepo := repo.NewGenesisBlockDataDTORepository(
		genesisBlockDataDTOConfigurationProvider,
		logger)
	_ = genesisBlockDataDTORepo

	//
	// Use-case
	//

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo)
	_ = createAccountUseCase
	getAccountUseCase := usecase.NewGetAccountUseCase(
		logger,
		accountRepo)
	_ = getAccountUseCase
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		logger,
		accountRepo)
	_ = getAccountsHashStateUseCase
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo)
	_ = upsertAccountUseCase
	accountsFilterByAddressesUseCase := usecase.NewAccountsFilterByAddressesUseCase(
		logger,
		accountRepo,
	)
	_ = accountsFilterByAddressesUseCase
	_ = getAccountsHashStateUseCase

	// Token
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokRepo,
	)
	_ = upsertTokenIfPreviousTokenNonceGTEUseCase
	// listTokensByOwnerUseCase := usecase.NewListTokensByOwnerUseCase(
	// 	logger,
	// 	tokRepo,
	// )
	// countTokensByOwnerUseCase := usecase.NewCountTokensByOwnerUseCase(
	// 	logger,
	// 	tokRepo,
	// )

	// // Genesis Block Data
	// upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
	// 	logger,
	// 	genesisBlockDataRepo)
	// getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
	// 	logger,
	// 	genesisBlockDataRepo)

	// // Genesis Block Data DTO
	// getGenesisBlockDataDTOFromBlockchainAuthorityUseCase := auth_usecase.NewGetGenesisBlockDataDTOFromBlockchainAuthorityUseCase(
	// 	logger,
	// 	genesisBlockDataDTORepo)

	// // Block Data
	// upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
	// 	logger,
	// 	blockDataRepo)
	// getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
	// 	logger,
	// 	blockDataRepo)

	// // Block Data DTO
	// getBlockDataDTOFromBlockchainAuthorityUseCase := auth_usecase.NewGetBlockDataDTOFromBlockchainAuthorityUseCase(
	// 	logger,
	// 	blockDataDTORepo)

	// // Blockchain State
	// upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
	// 	logger,
	// 	blockchainStateRepo)
	// getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
	// 	logger,
	// 	blockchainStateRepo)

	// // Blockchain State DTO
	// getBlockchainStateDTOFromBlockchainAuthorityUseCase := auth_usecase.NewGetBlockchainStateDTOFromBlockchainAuthorityUseCase(
	// 	logger,
	// 	blockchainStateDTORepo)

	// // Blockchain State
	// upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
	// 	logger,
	// 	blockchainStateRepo)
	// getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
	// 	logger,
	// 	blockchainStateRepo)

	// Blockchain State DTO
	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase := usecase.NewSubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase(
		logger,
		blockchainStateChangeEventDTORepo)

	_ = subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase

	//
	// Service
	//

	// blockchainSyncService := service.NewBlockchainSyncWithBlockchainAuthorityService(
	// 	logger,
	// 	getGenesisBlockDataUseCase,
	// 	upsertGenesisBlockDataUseCase,
	// 	getGenesisBlockDataDTOFromBlockchainAuthorityUseCase,
	// 	getBlockchainStateUseCase,
	// 	upsertBlockchainStateUseCase,
	// 	getBlockchainStateDTOFromBlockchainAuthorityUseCase,
	// 	getBlockDataUseCase,
	// 	upsertBlockDataUseCase,
	// 	getBlockDataDTOFromBlockchainAuthorityUseCase,
	// 	getAccountUseCase,
	// 	upsertAccountUseCase,
	// 	upsertTokenIfPreviousTokenNonceGTEUseCase,
	// )

	// blockchainSyncManagerService := service.NewBlockchainSyncManagerService(
	// 	logger,
	// 	blockchainSyncService,
	// 	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase,
	// )

	//
	// Interface.
	//

	// --- Task Manager --- //
	// poaConsensusMechanismTask := taskhandler.NewProofOfFaucetConsensusMechanismTaskHandler(
	// 	cfg,
	// 	logger,
	// 	proofOfFaucetConsensusMechanismService,
	// )
	taskManager := task.NewTaskManager(
		cfg,
		logger,
	)

	// --- HTTP --- //
	getVersionHTTPHandler := httphandler.NewGetVersionHTTPHandler(
		logger)
	getHealthCheckHTTPHandler := httphandler.NewGetHealthCheckHTTPHandler(
		logger)
	httpMiddleware := httpmiddle.NewMiddleware(
		logger,
		blackp)
	httpServ := http.NewHTTPServer(
		cfg,
		logger,
		httpMiddleware,
		getVersionHTTPHandler,
		getHealthCheckHTTPHandler,
	)

	//
	// STEP X
	// Execute.
	//

	// Load up our operating system interaction handlers, more specifically
	// signals. The OS sends our application various signals based on the
	// OS's state, we want to listen into the termination signals.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

	// Run in background
	go httpServ.Run()
	defer httpServ.Shutdown()
	go taskManager.Run()
	defer taskManager.Shutdown()

	logger.Info("ComicCoin Faucet is running.")

	<-done
}
