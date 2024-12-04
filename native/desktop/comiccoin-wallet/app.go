package main

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	auth_repo "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

// App struct
type App struct {
	ctx context.Context

	// Logger instance which provides detailed debugging information along
	// with the console log messages.
	logger *slog.Logger

	kmutex kmutexutil.KMutexProvider

	getAccountService                              *service.GetAccountService
	createAccountService                           *service.CreateAccountService
	accountListingByLocalWalletsService            *service.AccountListingByLocalWalletsService
	coinTransferService                            *service.CoinTransferService
	tokenGetService                                *service.TokenGetService
	tokenTransferService                           *service.TokenTransferService
	tokenBurnService                               *service.TokenBurnService
	getOrDownloadNonFungibleTokenService           *service.GetOrDownloadNonFungibleTokenService
	listBlockTransactionsByAddressService          *service.ListBlockTransactionsByAddressService
	listWithLimitBlockTransactionsByAddressService *service.ListWithLimitBlockTransactionsByAddressService
	getByBlockTransactionTimestampService          *service.GetByBlockTransactionTimestampService
	blockDataGetByHashService                      *service.BlockDataGetByHashService
	tokenListByOwnerService                        *service.TokenListByOwnerService
	tokenCountByOwnerService                       *service.TokenCountByOwnerService
	blockchainSyncService                          *service.BlockchainSyncWithBlockchainAuthorityService
	blockchainSyncManagerService                   *service.BlockchainSyncManagerService
	walletsFilterByLocalService                    *service.WalletsFilterByLocalService
	listNonFungibleTokensByOwnerService            *service.ListNonFungibleTokensByOwnerService
	pendingSignedTransactionListService            *service.PendingSignedTransactionListService
	exportWalletService                            *service.ExportWalletService
	importWalletService                            *service.ImportWalletService
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := logger.NewProvider()
	kmutex := kmutexutil.NewKMutexProvider()
	return &App{
		logger: logger,
		kmutex: kmutex,
	}
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	// Ensure that this function executes only one time and never concurrently.
	a.kmutex.Acquire("startup")
	defer a.kmutex.Release("startup")

	a.ctx = ctx
	a.logger.Debug("Startup beginning...")

	// DEVELOPERS NOTE:
	// Before we startup our app, we need to make sure the `data directory` is
	// set for this application by the user, else stop the app startup
	// proceedure. This is done on purpose because we need the user to specify
	// the location they want to store instead of having one automatically set.
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory
	if dataDir == "" {
		return
	}

	// DEVELOPERS NOTE:
	// Defensive code for programmer in case all the required environment
	// variables are not set then abort this program.
	// preferences.RunFatalIfHasAnyMissingFields() // ONLY USE IN CLI, NOT GUI WALLET!

	nftStorageAddress := preferences.NFTStorageAddress
	chainID := preferences.ChainID
	authorityAddress := preferences.AuthorityAddress

	//
	// DEVELOPERS NOTE:
	// Load up our dependencies and configuration
	//

	// ------ Common ------

	logger := logger.NewProvider()
	keystore := keystore.NewAdapter()
	walletDB := disk.NewDiskStorage(dataDir, "wallet", logger)
	accountDB := disk.NewDiskStorage(dataDir, "account", logger)
	genesisBlockDataDB := disk.NewDiskStorage(dataDir, "genesis_block_data", logger)
	blockchainStateDB := disk.NewDiskStorage(dataDir, "blockchain_state", logger)
	blockDataDB := disk.NewDiskStorage(dataDir, "block_data", logger)
	tokDB := disk.NewDiskStorage(dataDir, "token", logger)
	nftokDB := disk.NewDiskStorage(dataDir, "non_fungible_token", logger)
	pstxDB := disk.NewDiskStorage(dataDir, "pending_signed_transaction", logger)

	// ------------ Repo ------------

	accountRepo := repo.NewAccountRepo(
		logger,
		accountDB)
	walletRepo := repo.NewWalletRepo(
		logger,
		walletDB)
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		logger,
		genesisBlockDataDB)
	blockchainStateRepo := repo.NewBlockchainStateRepo(
		logger,
		blockchainStateDB)
	blockchainStateDTORepoConfig := auth_repo.NewBlockchainStateDTOConfigurationProvider(authorityAddress)
	blockchainStateDTORepo := auth_repo.NewBlockchainStateDTORepo(
		blockchainStateDTORepoConfig,
		logger)
	genesisBlockDataDTORepoConfig := auth_repo.NewGenesisBlockDataDTOConfigurationProvider(authorityAddress)
	genesisBlockDataDTORepo := auth_repo.NewGenesisBlockDataDTORepo(
		genesisBlockDataDTORepoConfig,
		logger)
	blockDataRepo := repo.NewBlockDataRepo(
		logger,
		blockDataDB)
	blockDataDTORepoConfig := auth_repo.NewBlockDataDTOConfigurationProvider(authorityAddress)
	blockDataDTORepo := auth_repo.NewBlockDataDTORepo(
		blockDataDTORepoConfig,
		logger)
	tokRepo := repo.NewTokenRepo(
		logger,
		tokDB)
	blockchainStateChangeEventDTOConfigurationProvider := auth_repo.NewBlockchainStateChangeEventDTOConfigurationProvider(authorityAddress)
	blockchainStateChangeEventDTORepo := auth_repo.NewBlockchainStateChangeEventDTORepo(
		blockchainStateChangeEventDTOConfigurationProvider,
		logger)
	nftokenRepo := repo.NewNonFungibleTokenRepo(logger, nftokDB)
	nftAssetRepoConfig := repo.NewNFTAssetRepoConfigurationProvider(nftStorageAddress, "")
	nftAssetRepo := repo.NewNFTAssetRepo(nftAssetRepoConfig, logger)
	mempoolTxDTORepoConfig := auth_repo.NewMempoolTransactionDTOConfigurationProvider(authorityAddress)
	mempoolTxDTORepo := auth_repo.NewMempoolTransactionDTORepo(mempoolTxDTORepoConfig, logger)
	pstxRepo := repo.NewPendingSignedTransactionRepo(logger, pstxDB)

	// ------------ Use-Case ------------

	// Storage Transaction
	storageTransactionOpenUseCase := usecase.NewStorageTransactionOpenUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo,
		tokRepo,
		pstxRepo)
	storageTransactionCommitUseCase := usecase.NewStorageTransactionCommitUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo,
		tokRepo,
		pstxRepo)
	storageTransactionDiscardUseCase := usecase.NewStorageTransactionDiscardUseCase(
		logger,
		walletRepo,
		accountRepo,
		genesisBlockDataRepo,
		blockchainStateRepo,
		blockDataRepo,
		tokRepo,
		pstxRepo)

	// Wallet
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		logger,
		keystore,
		walletRepo)
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		logger,
		keystore,
		walletRepo)
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		logger,
		walletRepo)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		logger,
		walletRepo)
	listAllWalletUseCase := usecase.NewListAllWalletUseCase(
		logger,
		walletRepo)
	listAllAddressesWalletUseCase := usecase.NewListAllAddressesWalletUseCase(
		logger,
		walletRepo,
	)

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		logger,
		accountRepo)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		logger,
		accountRepo)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		logger,
		accountRepo)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		logger,
		accountRepo)
	accountsFilterByAddressesUseCase := usecase.NewAccountsFilterByAddressesUseCase(
		logger,
		accountRepo,
	)
	_ = getAccountsHashStateUseCase

	// Blockchain State
	upsertBlockchainStateUseCase := usecase.NewUpsertBlockchainStateUseCase(
		logger,
		blockchainStateRepo)
	getBlockchainStateUseCase := usecase.NewGetBlockchainStateUseCase(
		logger,
		blockchainStateRepo)

	// Blockchain State DTO
	getBlockchainStateDTOFromBlockchainAuthorityUseCase := auth_usecase.NewGetBlockchainStateDTOFromBlockchainAuthorityUseCase(
		logger,
		blockchainStateDTORepo)

	// Genesis Block Data
	upsertGenesisBlockDataUseCase := usecase.NewUpsertGenesisBlockDataUseCase(
		logger,
		genesisBlockDataRepo)
	getGenesisBlockDataUseCase := usecase.NewGetGenesisBlockDataUseCase(
		logger,
		genesisBlockDataRepo)

	// Genesis Block Data DTO
	getGenesisBlockDataDTOFromBlockchainAuthorityUseCase := auth_usecase.NewGetGenesisBlockDataDTOFromBlockchainAuthorityUseCase(
		logger,
		genesisBlockDataDTORepo)

	// Block Data
	upsertBlockDataUseCase := usecase.NewUpsertBlockDataUseCase(
		logger,
		blockDataRepo)
	getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
		logger,
		blockDataRepo)
	getByBlockTransactionTimestampUseCase := usecase.NewGetByBlockTransactionTimestampUseCase(
		logger,
		blockDataRepo)

	// Block Transactions
	listBlockTransactionsByAddressUseCase := usecase.NewListBlockTransactionsByAddressUseCase(
		logger,
		blockDataRepo)
	listWithLimitBlockTransactionsByAddressUseCase := usecase.NewListWithLimitBlockTransactionsByAddressUseCase(
		logger,
		blockDataRepo)

	// Block Data DTO
	getBlockDataDTOFromBlockchainAuthorityUseCase := auth_usecase.NewGetBlockDataDTOFromBlockchainAuthorityUseCase(
		logger,
		blockDataDTORepo)

	// Token
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		logger,
		tokRepo,
	)
	listTokensByOwnerUseCase := usecase.NewListTokensByOwnerUseCase(
		logger,
		tokRepo,
	)
	countTokensByOwnerUseCase := usecase.NewCountTokensByOwnerUseCase(
		logger,
		tokRepo,
	)

	// Blockchain State DTO
	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase := auth_usecase.NewSubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase(
		logger,
		blockchainStateChangeEventDTORepo)

	// Token
	getTokUseCase := usecase.NewGetTokenUseCase(
		logger,
		tokRepo)

	// Non-Fungible Token
	getNFTokUseCase := usecase.NewGetNonFungibleTokenUseCase(
		logger,
		nftokenRepo)
	downloadNFTokMetadataUsecase := usecase.NewDownloadMetadataNonFungibleTokenUseCase(
		logger,
		nftAssetRepo)
	downloadNFTokAssetUsecase := usecase.NewDownloadNonFungibleTokenAssetUseCase(
		logger,
		nftAssetRepo)
	upsertNFTokUseCase := usecase.NewUpsertNonFungibleTokenUseCase(
		logger,
		nftokenRepo)
	listNonFungibleTokensWithFilterByTokenIDsyUseCase := usecase.NewListNonFungibleTokensWithFilterByTokenIDsyUseCase(
		logger,
		nftokenRepo)

	// Mempool Transaction DTO
	submitMempoolTransactionDTOToBlockchainAuthorityUseCase := auth_usecase.NewSubmitMempoolTransactionDTOToBlockchainAuthorityUseCase(
		logger,
		mempoolTxDTORepo,
	)

	// Pending Signed Transaction
	upsertPendingSignedTransactionUseCase := usecase.NewUpsertPendingSignedTransactionUseCase(
		logger,
		pstxRepo)
	listPendingSignedTransactionUseCase := usecase.NewListPendingSignedTransactionUseCase(
		logger,
		pstxRepo)
	deletePendingSignedTransactionUseCase := usecase.NewDeletePendingSignedTransactionUseCase(
		logger,
		pstxRepo)

	// ------------ Service ------------

	getAccountService := service.NewGetAccountService(
		logger,
		getAccountUseCase,
	)
	createAccountService := service.NewCreateAccountService(
		logger,
		walletEncryptKeyUseCase,
		walletDecryptKeyUseCase,
		createWalletUseCase,
		createAccountUseCase,
		getAccountUseCase,
	)
	accountListingByLocalWalletsService := service.NewAccountListingByLocalWalletsService(
		logger,
		listAllAddressesWalletUseCase,
		accountsFilterByAddressesUseCase,
	)
	coinTransferService := service.NewCoinTransferService(
		logger,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		listPendingSignedTransactionUseCase,
		getGenesisBlockDataUseCase,
		upsertPendingSignedTransactionUseCase,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		submitMempoolTransactionDTOToBlockchainAuthorityUseCase,
	)
	tokenGetService := service.NewTokenGetService(
		logger,
		getTokUseCase,
	)
	tokenTransferService := service.NewTokenTransferService(
		logger,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		listPendingSignedTransactionUseCase,
		getGenesisBlockDataUseCase,
		upsertPendingSignedTransactionUseCase,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getTokUseCase,
		submitMempoolTransactionDTOToBlockchainAuthorityUseCase,
	)
	tokenBurnService := service.NewTokenBurnService(
		logger,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		listPendingSignedTransactionUseCase,
		getGenesisBlockDataUseCase,
		upsertPendingSignedTransactionUseCase,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getTokUseCase,
		submitMempoolTransactionDTOToBlockchainAuthorityUseCase,
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
		deletePendingSignedTransactionUseCase,
	)

	blockchainSyncManagerService := service.NewBlockchainSyncManagerService(
		logger,
		blockchainSyncService,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase,
	)

	getOrDownloadNonFungibleTokenService := service.NewGetOrDownloadNonFungibleTokenService(
		logger,
		getNFTokUseCase,
		getTokUseCase,
		downloadNFTokMetadataUsecase,
		downloadNFTokAssetUsecase,
		upsertNFTokUseCase)

	listBlockTransactionsByAddressService := service.NewListBlockTransactionsByAddressService(
		logger,
		listBlockTransactionsByAddressUseCase,
	)
	listWithLimitBlockTransactionsByAddressService := service.NewListWithLimitBlockTransactionsByAddressService(
		logger,
		listWithLimitBlockTransactionsByAddressUseCase,
	)
	getByBlockTransactionTimestampService := service.NewGetByBlockTransactionTimestampService(
		logger,
		getByBlockTransactionTimestampUseCase,
	)
	blockDataGetByHashService := service.NewBlockDataGetByHashService(
		logger,
		getBlockDataUseCase,
	)
	tokenListByOwnerService := service.NewTokenListByOwnerService(
		logger,
		listTokensByOwnerUseCase,
	)
	tokenCountByOwnerService := service.NewTokenCountByOwnerService(
		logger,
		countTokensByOwnerUseCase,
	)
	walletsFilterByLocalService := service.NewWalletsFilterByLocalService(
		logger,
		listAllWalletUseCase,
	)
	listNonFungibleTokensByOwnerService := service.NewListNonFungibleTokensByOwnerService(
		logger,
		listTokensByOwnerUseCase,
		listNonFungibleTokensWithFilterByTokenIDsyUseCase,
		getOrDownloadNonFungibleTokenService,
	)
	pendingSignedTransactionListService := service.NewPendingSignedTransactionListService(
		logger,
		listPendingSignedTransactionUseCase,
	)
	exportWalletService := service.NewExportWalletService(
		logger,
		getAccountUseCase,
		getWalletUseCase,
	)
	importWalletService := service.NewImportWalletService(
		logger,
		getAccountUseCase,
		getWalletUseCase,
		upsertAccountUseCase,
		createWalletUseCase,
	)

	// ------------ Interfaces ------------

	a.getAccountService = getAccountService
	a.createAccountService = createAccountService
	a.accountListingByLocalWalletsService = accountListingByLocalWalletsService
	a.coinTransferService = coinTransferService
	a.tokenGetService = tokenGetService
	a.tokenTransferService = tokenTransferService
	a.tokenBurnService = tokenBurnService
	a.blockchainSyncService = blockchainSyncService
	a.blockchainSyncManagerService = blockchainSyncManagerService
	a.getOrDownloadNonFungibleTokenService = getOrDownloadNonFungibleTokenService
	a.listBlockTransactionsByAddressService = listBlockTransactionsByAddressService
	a.listWithLimitBlockTransactionsByAddressService = listWithLimitBlockTransactionsByAddressService
	a.getByBlockTransactionTimestampService = getByBlockTransactionTimestampService
	a.blockDataGetByHashService = blockDataGetByHashService
	a.tokenListByOwnerService = tokenListByOwnerService
	a.tokenCountByOwnerService = tokenCountByOwnerService
	a.walletsFilterByLocalService = walletsFilterByLocalService
	a.listNonFungibleTokensByOwnerService = listNonFungibleTokensByOwnerService
	a.pendingSignedTransactionListService = pendingSignedTransactionListService
	a.exportWalletService = exportWalletService
	a.importWalletService = importWalletService

	//
	// Execute.
	//

	go func(ctx context.Context, cid uint16) {
		for {
			a.logger.Debug("Starting sync-manager...")
			if err := blockchainSyncManagerService.Execute(ctx, cid); err != nil {
				a.logger.Error("Failed to manage syncing", slog.Any("error", err))
			}
			a.logger.Debug("Sync-manager will restart again.")
		}
	}(a.ctx, chainID)

	logger.Info("ComicCoin Wallet is running.")
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) shutdown(ctx context.Context) {
	a.logger.Debug("Shutting down now...")
	defer a.logger.Debug("Shutting down finished")

	// DEVELOPERS NOTE:
	// Before we startup our app, we need to make sure the `data directory` is
	// set for this application by the user, else stop the app startup
	// proceedure. This is done on purpose because we need the user to specify
	// the location they want to store instead of having one automatically set.
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory
	if dataDir == "" {
		return
	}

}
