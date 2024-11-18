package daemon

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/blacklist"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/security/password"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http"
	httphandler "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/task"
	taskhandler "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/interface/task/handler"
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
	kmutex := kmutexutil.NewKMutexProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter()
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
	mempoolTxRepo := repo.NewMempoolTransactionRepo(cfg, logger, dbClient)
	tokenRepo := repo.NewTokenRepo(cfg, logger, dbClient)

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
	upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
		cfg,
		logger,
		bcStateRepo,
	)
	blockchainStateUpdateDetectorUseCase := usecase.NewBlockchainStateUpdateDetectorUseCase(
		cfg,
		logger,
		bcStateRepo,
	)
	defer func() {
		// When we are done, we will need to terminate our access to this resource.
		blockchainStateUpdateDetectorUseCase.Terminate()
	}()

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
	upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
		cfg,
		logger,
		bdRepo,
	)

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
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		cfg,
		logger,
		accountRepo,
	)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo,
	)

	// Token
	getTokenUseCase := usecase.NewGetTokenUseCase(
		cfg,
		logger,
		tokenRepo,
	)
	getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
		cfg,
		logger,
		tokenRepo,
	)
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokenRepo,
	)

	// Mempool Transaction
	mempoolTransactionCreateUseCase := usecase.NewMempoolTransactionCreateUseCase(
		cfg,
		logger,
		mempoolTxRepo,
	)
	mempoolTransactionListByChainIDUseCase := usecase.NewMempoolTransactionListByChainIDUseCase(
		cfg,
		logger,
		mempoolTxRepo,
	)
	_ = mempoolTransactionListByChainIDUseCase
	mempoolTransactionDeleteByIDUseCase := usecase.NewMempoolTransactionDeleteByIDUseCase(
		cfg,
		logger,
		mempoolTxRepo,
	)
	mempoolTransactionInsertionDetectorUseCase := usecase.NewMempoolTransactionInsertionDetectorUseCase(
		cfg,
		logger,
		mempoolTxRepo,
	)
	defer func() {
		// When we are done, we will need to terminate our access to this resource.
		mempoolTransactionInsertionDetectorUseCase.Terminate()
	}()

	// Proof of Work
	proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(
		cfg,
		logger,
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

	// MempoolTransaction
	mempoolTransactionReceiveDTOFromNetworkService := service.NewMempoolTransactionReceiveDTOFromNetworkService(
		cfg,
		logger,
		mempoolTransactionCreateUseCase,
	)

	// Proof of Authority Consensus Mechanism
	getProofOfAuthorityPrivateKeyService := service.NewGetProofOfAuthorityPrivateKeyService(
		cfg,
		logger,
		getWalletUseCase,
		walletDecryptKeyUseCase,
	)
	proofOfAuthorityConsensusMechanismService := service.NewProofOfAuthorityConsensusMechanismService(
		cfg,
		logger,
		kmutex,
		dbClient, // We do this so we can use MongoDB's "transactions"
		getProofOfAuthorityPrivateKeyService,
		mempoolTransactionInsertionDetectorUseCase,
		mempoolTransactionDeleteByIDUseCase,
		getBlockchainStateUseCase,
		upsertBlockchainStateUseCase,
		getGenesisBlockDataUseCase,
		getBlockDataUseCase,
		getAccountUseCase,
		getAccountsHashStateUseCase,
		upsertAccountUseCase,
		getTokenUseCase,
		getTokensHashStateUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
		proofOfWorkUseCase,
		upsertBlockDataUseCase,
	)

	//
	// Interface.
	//

	// --- Task Manager --- //
	poaConsensusMechanismTask := taskhandler.NewProofOfAuthorityConsensusMechanismTaskHandler(
		cfg,
		logger,
		proofOfAuthorityConsensusMechanismService,
	)
	taskManager := task.NewTaskManager(
		cfg,
		logger,
		poaConsensusMechanismTask,
	)

	// --- HTTP --- //
	getVersionHTTPHandler := httphandler.NewGetVersionHTTPHandler(
		logger)
	getHealthCheckHTTPHandler := httphandler.NewGetHealthCheckHTTPHandler(
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
	blockchainStateChangeEventsHTTPHandler := httphandler.NewBlockchainStateChangeEventDTOHTTPHandler(
		logger,
		blockchainStateUpdateDetectorUseCase)
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
	mempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler := httphandler.NewMempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler(
		logger,
		mempoolTransactionReceiveDTOFromNetworkService)
	httpMiddleware := httpmiddle.NewMiddleware(
		logger,
		blackp)
	httpServ := http.NewHTTPServer(
		cfg,
		logger,
		httpMiddleware,
		getVersionHTTPHandler,
		getHealthCheckHTTPHandler,
		getGenesisBlockDataHTTPHandler,
		getBlockchainStateHTTPHandler,
		blockchainStateChangeEventsHTTPHandler,
		listAllBlockDataOrderedHashesHTTPHandler,
		listAllBlockDataUnorderedHashesHTTPHandler,
		getBlockDataHTTPHandler,
		listBlockDataFilteredInHashesHTTPHandler,
		listBlockDataFilteredBetweenBlockNumbersInChainIDHTTPHandler,
		signedTransactionSubmissionHTTPHandler,
		mempoolTransactionReceiveDTOFromNetworkServiceHTTPHandler,
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

	logger.Info("ComicCoin Authority is running.")

	<-done
}
