package daemon

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/emailer/mailgun"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/blacklist"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	cloudstorage "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/cloud/s3"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/templatedemailer"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http"
	httphandler "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/http/middleware"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/interface/task"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/repo"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
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
	kmutex := kmutexutil.NewKMutexProvider()
	cfg := config.NewProviderUsingEnvironmentVariables()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter()
	passp := password.NewProvider()
	blackp := blacklist.NewProvider()
	jwtp := jwt.NewProvider(cfg)
	cache := mongodbcache.NewCache(cfg, logger, dbClient)
	emailer := mailgun.NewEmailer(cfg, logger)
	templatedEmailer := templatedemailer.NewTemplatedEmailer(logger, emailer)
	cloudstore := cloudstorage.NewCloudStorage(cfg, logger)
	_ = cloudstore //TODO: Utilize in our app.

	//
	// Repository
	//

	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)
	tenantRepo := repo.NewTenantRepository(cfg, logger, dbClient)
	userRepo := repo.NewUserRepository(cfg, logger, dbClient)
	tokRepo := repo.NewTokenRepository(cfg, logger, dbClient)
	blockchainStateRepo := repo.NewBlockchainStateRepository(cfg, logger, dbClient)
	blockchainStateChangeEventDTOConfigurationProvider := repo.NewBlockchainStateChangeEventDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockchainStateChangeEventDTORepo := repo.NewBlockchainStateChangeEventDTORepo(
		blockchainStateChangeEventDTOConfigurationProvider,
		logger)
	blockchainStateDTOConfigurationProvider := repo.NewBlockchainStateDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockchainStateDTORepo := repo.NewBlockchainStateDTORepo(
		blockchainStateDTOConfigurationProvider,
		logger)
	blockDataRepo := repo.NewBlockDataRepository(
		cfg,
		logger,
		dbClient)
	blockDataDTOConfigurationProvider := repo.NewBlockDataDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	blockDataDTORepo := repo.NewBlockDataDTORepository(
		blockDataDTOConfigurationProvider,
		logger)
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepository(
		cfg,
		logger,
		dbClient)
	genesisBlockDataDTOConfigurationProvider := repo.NewGenesisBlockDataDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	genesisBlockDataDTORepo := repo.NewGenesisBlockDataDTORepository(
		genesisBlockDataDTOConfigurationProvider,
		logger)
	mempoolTransactionDTOConfigurationProvider := repo.NewMempoolTransactionDTOConfigurationProvider(cfg.App.AuthorityHTTPAddress)
	mempoolTxDTORepo := repo.NewMempoolTransactionDTORepo(mempoolTransactionDTOConfigurationProvider, logger)

	//
	// Use-case
	//

	// Email
	sendUserVerificationEmailUseCase := usecase.NewSendUserVerificationEmailUseCase(cfg, logger, templatedEmailer)

	// Wallet
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		keystore,
		walletRepo,
	)
	_ = walletDecryptKeyUseCase
	getWalletUseCase := usecase.NewGetWalletUseCase(
		cfg,
		logger,
		walletRepo,
	)
	_ = getWalletUseCase

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo)
	_ = createAccountUseCase
	getAccountUseCase := usecase.NewGetAccountUseCase(
		logger,
		accountRepo)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		logger,
		accountRepo)
	_ = getAccountsHashStateUseCase
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo)
	accountsFilterByAddressesUseCase := usecase.NewAccountsFilterByAddressesUseCase(
		logger,
		accountRepo,
	)
	_ = accountsFilterByAddressesUseCase

	// Genesis Block Data
	getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
		cfg,
		logger,
		genesisBlockDataRepo,
	)
	upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
		logger,
		genesisBlockDataRepo,
	)

	// Genesis Block Data DTO
	getGenesisBlockDataDTOFromBlockchainAuthorityUseCase := usecase.NewGetGenesisBlockDataDTOFromBlockchainAuthorityUseCase(
		logger,
		genesisBlockDataDTORepo)

	// Blockchain State
	getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
		cfg,
		logger,
		blockchainStateRepo,
	)
	upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
		cfg,
		logger,
		blockchainStateRepo,
	)
	_ = upsertBlockchainStateUseCase

	// Block Data
	getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo,
	)
	upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo,
	)
	listBlockTransactionsByAddressUseCase := usecase.NewListBlockTransactionsByAddressUseCase(
		cfg,
		logger,
		blockDataRepo,
	)
	_ = listBlockTransactionsByAddressUseCase

	// Block Data DTO
	getBlockDataDTOFromBlockchainAuthorityUseCase := usecase.NewGetBlockDataDTOFromBlockchainAuthorityUseCase(
		logger,
		blockDataDTORepo)

	// Token
	getTokenUseCase := usecase.NewGetTokenUseCase(
		logger,
		tokRepo,
	)
	_ = getTokenUseCase
	getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
		logger,
		tokRepo,
	)
	_ = getTokensHashStateUseCase
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokRepo,
	)
	listTokensByOwnerUseCase := usecase.NewListTokensByOwnerUseCase(
		logger,
		tokRepo,
	)
	_ = listTokensByOwnerUseCase

	// Blockchain State DTO
	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase := usecase.NewSubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase(
		logger,
		blockchainStateChangeEventDTORepo)

	_ = subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase

	// Blockchain State DTO
	getBlockchainStateDTOFromBlockchainAuthorityUseCase := usecase.NewGetBlockchainStateDTOFromBlockchainAuthorityUseCase(
		logger,
		blockchainStateDTORepo)

	// Mempooltx DTO
	submitMempoolTransactionDTOToBlockchainAuthorityUseCase := usecase.NewSubmitMempoolTransactionDTOToBlockchainAuthorityUseCase(
		logger,
		mempoolTxDTORepo,
	)

	// Tenant
	tenantGetByIDUseCase := usecase.NewTenantGetByIDUseCase(
		cfg,
		logger,
		tenantRepo)
	tenantUpdateUseCase := usecase.NewTenantUpdateUseCase(
		cfg,
		logger,
		tenantRepo)

	// User
	userGetByEmailUseCase := usecase.NewUserGetByEmailUseCase(
		cfg,
		logger,
		userRepo)
	userCreateUseCase := usecase.NewUserCreateUseCase(
		cfg,
		logger,
		userRepo)
	userUpdateUseCase := usecase.NewUserUpdateUseCase(
		cfg,
		logger,
		userRepo)
	userGetBySessionIDUseCase := usecase.NewUserGetBySessionIDUseCase(
		cfg,
		logger,
		cache)
	userGetByIDUseCase := usecase.NewUserGetByIDUseCase(
		cfg,
		logger,
		userRepo)
	userGetByVerificationCodeUseCase := usecase.NewUserGetByVerificationCodeUseCase(
		cfg,
		logger,
		userRepo)

	//
	// Service
	//

	faucetCoinTransferService := service.NewFaucetCoinTransferService(
		logger,
		kmutex,
		tenantGetByIDUseCase,
		tenantUpdateUseCase,
		getAccountUseCase,
		upsertAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		submitMempoolTransactionDTOToBlockchainAuthorityUseCase,
	)

	gatewayUserRegisterService := service.NewGatewayUserRegisterService(
		cfg,
		logger,
		passp,
		cache,
		jwtp,
		tenantGetByIDUseCase,
		userGetByEmailUseCase,
		userCreateUseCase,
		userUpdateUseCase,
		sendUserVerificationEmailUseCase,
	)

	gatewayLoginService := service.NewGatewayLoginService(
		logger,
		passp,
		cache,
		jwtp,
		tenantGetByIDUseCase,
		userGetByEmailUseCase,
		userUpdateUseCase,
	)

	gatewayLogoutService := service.NewGatewayLogoutService(
		logger,
		cache,
	)

	gatewayRefreshTokenService := service.NewGatewayRefreshTokenService(
		logger,
		cache,
		jwtp,
		userGetByEmailUseCase,
	)

	gatewayProfileGetService := service.NewGatewayProfileGetService(
		logger,
		userGetByIDUseCase,
	)
	gatewayProfileUpdateService := service.NewGatewayProfileUpdateService(
		logger,
		userGetByIDUseCase,
		userUpdateUseCase,
	)
	gatewayVerifyService := service.NewGatewayVerifyService(
		logger,
		kmutex,
		userGetByVerificationCodeUseCase,
		userUpdateUseCase,
	)
	gatewayChangePasswordService := service.NewGatewayChangePasswordService(
		logger,
		kmutex,
		passp,
		userGetByIDUseCase,
		userUpdateUseCase,
	)
	gatewayForgotPasswordService := service.NewGatewayForgotPasswordService(
		logger,
		templatedEmailer,
		userGetByEmailUseCase,
		userUpdateUseCase,
	)
	gatewayResetPasswordService := service.NewGatewayResetPasswordService(
		logger,
		kmutex,
		passp,
		userGetByVerificationCodeUseCase,
		userUpdateUseCase,
	)
	gatewayAddWalletAddressToFaucetService := service.NewGatewayAddWalletAddressToFaucetService(
		cfg,
		logger,
		tenantGetByIDUseCase,
		userGetByIDUseCase,
		userUpdateUseCase,
		faucetCoinTransferService,
	)

	blockchainSyncService := service.NewBlockchainSyncWithBlockchainAuthorityService(
		logger,
		getGenesisBlockDataUseCase,
		upsertGenesisBlockDataUseCase,
		getGenesisBlockDataDTOFromBlockchainAuthorityUseCase,
		getBlockchainStateUseCase,
		upsertBlockchainStateUseCase,
		getBlockchainStateDTOFromBlockchainAuthorityUseCase,
		getBlockDataUseCase,
		upsertBlockDataUseCase,
		getBlockDataDTOFromBlockchainAuthorityUseCase,
		getAccountUseCase,
		upsertAccountUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
		tenantGetByIDUseCase,
		tenantUpdateUseCase,
	)

	blockchainSyncManagerService := service.NewBlockchainSyncManagerService(
		logger,
		blockchainSyncService,
		subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase,
	)

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
	gatewayUserRegisterHTTPHandler := httphandler.NewGatewayUserRegisterHTTPHandler(
		logger,
		dbClient,
		gatewayUserRegisterService,
	)
	gatewayLoginHTTPHandler := httphandler.NewGatewayLoginHTTPHandler(
		logger,
		dbClient,
		gatewayLoginService,
	)
	gatewayLogoutHTTPHandler := httphandler.NewGatewayLogoutHTTPHandler(
		logger,
		dbClient,
		gatewayLogoutService,
	)
	gatewayRefreshTokenHTTPHandler := httphandler.NewGatewayRefreshTokenHTTPHandler(
		logger,
		dbClient,
		gatewayRefreshTokenService,
	)
	gatewayProfileDetailHTTPHandler := httphandler.NewGatewayProfileDetailHTTPHandler(
		logger,
		dbClient,
		gatewayProfileGetService,
	)
	gatewayProfileUpdateHTTPHandler := httphandler.NewGatewayProfileUpdateHTTPHandler(
		logger,
		dbClient,
		gatewayProfileUpdateService,
	)
	gatewayVerifyHTTPHandler := httphandler.NewGatewayVerifyHTTPHandler(
		logger,
		dbClient,
		gatewayVerifyService,
	)
	gatewayChangePasswordHTTPHandler := httphandler.NewGatewayChangePasswordHTTPHandler(
		logger,
		dbClient,
		gatewayChangePasswordService,
	)
	gatewayForgotPasswordHTTPHandler := httphandler.NewGatewayForgotPasswordHTTPHandler(
		logger,
		dbClient,
		gatewayForgotPasswordService,
	)
	gatewayResetPasswordHTTPHandler := httphandler.NewGatewayResetPasswordHTTPHandler(
		logger,
		dbClient,
		gatewayResetPasswordService,
	)
	gatewayProfileWalletAddressHTTPHandler := httphandler.NewGatewayProfileWalletAddressHTTPHandler(
		logger,
		dbClient,
		gatewayAddWalletAddressToFaucetService,
	)
	httpMiddleware := httpmiddle.NewMiddleware(
		logger,
		blackp,
		jwtp,
		userGetBySessionIDUseCase,
	)
	httpServ := http.NewHTTPServer(
		cfg,
		logger,
		httpMiddleware,
		getVersionHTTPHandler,
		getHealthCheckHTTPHandler,
		gatewayUserRegisterHTTPHandler,
		gatewayLoginHTTPHandler,
		gatewayLogoutHTTPHandler,
		gatewayRefreshTokenHTTPHandler,
		gatewayProfileDetailHTTPHandler,
		gatewayProfileUpdateHTTPHandler,
		gatewayVerifyHTTPHandler,
		gatewayChangePasswordHTTPHandler,
		gatewayForgotPasswordHTTPHandler,
		gatewayResetPasswordHTTPHandler,
		gatewayProfileWalletAddressHTTPHandler,
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
	go func() {
		for {
			ctx := context.Background()
			if err := blockchainSyncManagerService.Execute(ctx, cfg.Blockchain.ChainID, cfg.App.TenantID); err != nil {
				log.Fatalf("Failed to manage syncing: %v\n", err)
			}
		}
	}()
	go httpServ.Run()
	defer httpServ.Shutdown()
	go taskManager.Run()
	defer taskManager.Shutdown()

	logger.Info("ComicCoin Faucet is running.")

	<-done
}
