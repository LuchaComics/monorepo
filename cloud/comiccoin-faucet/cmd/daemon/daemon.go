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

	tenantRepo := repo.NewTenantRepository(cfg, logger, dbClient)
	_ = tenantRepo
	userRepo := repo.NewUserRepository(cfg, logger, dbClient)
	_ = userRepo

	blockchainStateChangeEventDTOConfigurationProvider := repo.NewBlockchainStateChangeEventDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockchainStateChangeEventDTORepo := repo.NewBlockchainStateChangeEventDTORepo(
		blockchainStateChangeEventDTOConfigurationProvider,
		logger)

	//
	// Use-case
	//

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
