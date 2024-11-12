package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/blacklist"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/password"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http"
	httphandler "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/service"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run the comiccoin authority",
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
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter(cfg, logger)
	passp := password.NewProvider()
	jwtp := jwt.NewProvider(cfg)
	blackp := blacklist.NewProvider()

	_ = passp
	_ = jwtp

	// Repository
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	bdRepo := repo.NewBlockDataRepo(cfg, logger, dbClient)
	gbdRepo := repo.NewGenesisBlockDataRepo(cfg, logger, dbClient)
	bcStateRepo := repo.NewBlockchainStateRepo(cfg, logger, dbClient)

	_ = keystore
	_ = walletRepo
	_ = accountRepo

	// Genesis Block Data
	getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
		cfg,
		logger,
		gbdRepo,
	)

	// Blockchain State
	getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
		cfg,
		logger,
		bcStateRepo,
	)

	// Block Data
	listAllBlockNumberByHashArrayUseCase := usecase.NewListAllBlockNumberByHashArrayUseCase(
		cfg,
		logger,
		bdRepo,
	)
	listBlockDataUnorderedHashArrayUseCase := usecase.NewListBlockDataUnorderedHashArrayUseCase(
		cfg,
		logger,
		bdRepo,
	)
	getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
		cfg,
		logger,
		bdRepo,
	)
	listBlockDataFilteredInHashesUseCase := usecase.NewListBlockDataFilteredInHashesUseCase(
		cfg,
		logger,
		bdRepo,
	)
	listBlockDataFilteredBetweenBlockNumbersUseCase := usecase.NewListBlockDataFilteredBetweenBlockNumbersUseCase(
		cfg,
		logger,
		bdRepo,
	)

	// --- Service

	// Genesis
	getGenesisBlockDataService := service.NewGetGenesisBlockDataService(
		cfg,
		logger,
		getGenesisBlockDataUseCase,
	)

	// Blockchain State
	getBlockchainStateService := service.NewGetBlockchainStateService(
		cfg,
		logger,
		getBlockchainStateUseCase,
	)

	// Block Data
	getBlockDataService := service.NewGetBlockDataService(
		cfg,
		logger,
		getBlockDataUseCase,
	)
	blockDataListAllOrderedHashesService := service.NewBlockDataListAllOrderedHashesService(
		cfg,
		logger,
		listAllBlockNumberByHashArrayUseCase,
	)
	blockDataListAllUnorderedHashesService := service.NewBlockDataListAllUnorderedHashesService(
		cfg,
		logger,
		listBlockDataUnorderedHashArrayUseCase,
	)
	listBlockDataFilteredInHashesService := service.NewListBlockDataFilteredInHashesService(
		cfg,
		logger,
		listBlockDataFilteredInHashesUseCase,
	)
	listBlockDataFilteredBetweenBlockNumbersInChainIDService := service.NewListBlockDataFilteredBetweenBlockNumbersInChainIDService(
		cfg,
		logger,
		listBlockDataFilteredBetweenBlockNumbersUseCase,
	)

	// Coins
	signedTransactionSubmissionService := service.NewSignedTransactionSubmissionService(
		cfg,
		logger,
	)

	//
	// Interface.
	//

	// --- HTTP --- //
	getVersionHTTPHandler := httphandler.NewGetVersionHTTPHandler(
		logger)
	getGenesisBlockDataHTTPHandler := httphandler.NewGetGenesisBlockDataHTTPHandler(
		logger,
		getGenesisBlockDataService)
	getBlockDataHTTPHandler := httphandler.NewGetBlockDataHTTPHandler(
		logger,
		getBlockDataService)
	getBlockchainStateHTTPHandler := httphandler.NewGetBlockchainStateHTTPHandler(
		logger,
		getBlockchainStateService)
	listAllBlockDataOrderedHashesHTTPHandler := httphandler.NewListAllBlockDataOrderedHashesHTTPHandler(
		logger,
		blockDataListAllOrderedHashesService)
	listAllBlockDataUnorderedHashesHTTPHandler := httphandler.NewListAllBlockDataUnorderedHashesHTTPHandler(
		logger,
		blockDataListAllUnorderedHashesService)
	listBlockDataFilteredInHashesHTTPHandler := httphandler.NewListBlockDataFilteredInHashesHTTPHandler(
		logger,
		listBlockDataFilteredInHashesService)
	listBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler := httphandler.NewListBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler(
		logger,
		listBlockDataFilteredBetweenBlockNumbersInChainIDService)
	signedTransactionSubmissionHTTPHandler := httphandler.NewSignedTransactionSubmissionHTTPHandler(
		logger,
		signedTransactionSubmissionService)
	httpMiddleware := httpmiddle.NewMiddleware(
		logger,
		blackp)
	httpServ := http.NewHTTPServer(
		cfg,
		logger,
		httpMiddleware,
		getVersionHTTPHandler,
		getGenesisBlockDataHTTPHandler,
		getBlockchainStateHTTPHandler,
		listAllBlockDataOrderedHashesHTTPHandler,
		listAllBlockDataUnorderedHashesHTTPHandler,
		getBlockDataHTTPHandler,
		listBlockDataFilteredInHashesHTTPHandler,
		listBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler,
		signedTransactionSubmissionHTTPHandler,
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

	// Run in background the peer to peer node which will synchronize our
	// blockchain with the network.
	// go peerNode.Run()
	go httpServ.Run()
	defer httpServ.Shutdown()

	logger.Info("Node running.")

	<-done
}
