package daemon

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	auth_repo "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	"github.com/spf13/cobra"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/preferences"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/rpc"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

var (
	preferences *pref.Preferences
)

// Command line argument flags
var (
	flagDataDirectory     string
	flagChainID           uint16
	flagAuthorityAddress  string
	flagNFTStorageAddress string
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Runs a full node on your machine which will automatically synchronize the local blockchain with the Global Blockchain Network on any new changes.",
		Run: func(cmd *cobra.Command, args []string) {
			// Developers Note:
			// Before executing this command, check to ensure the user has
			// configured our app before proceeding.
			preferences.RunFatalIfHasAnyMissingFields()

			// Load up our operating system interaction handlers, more specifically
			// signals. The OS sends our application various signals based on the
			// OS's state, we want to listen into the termination signals.
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

			go doRunDaemonCmd()

			<-done
		},
	}
	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	return cmd
}

func doRunDaemonCmd() {
	// ------ Common ------

	logger := logger.NewProvider()
	keystore := keystore.NewAdapter()
	walletDB := disk.NewDiskStorage(flagDataDirectory, "wallet", logger)
	accountDB := disk.NewDiskStorage(flagDataDirectory, "account", logger)
	genesisBlockDataDB := disk.NewDiskStorage(flagDataDirectory, "genesis_block_data", logger)
	blockchainStateDB := disk.NewDiskStorage(flagDataDirectory, "blockchain_state", logger)
	blockDataDB := disk.NewDiskStorage(flagDataDirectory, "block_data", logger)
	tokDB := disk.NewDiskStorage(flagDataDirectory, "token", logger)
	nftokDB := disk.NewDiskStorage(flagDataDirectory, "non_fungible_token", logger)
	pstxDB := disk.NewDiskStorage(flagDataDirectory, "pending_signed_transaction", logger)

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
	blockchainStateDTORepoConfig := auth_repo.NewBlockchainStateDTOConfigurationProvider(flagAuthorityAddress)
	blockchainStateDTORepo := auth_repo.NewBlockchainStateDTORepo(
		blockchainStateDTORepoConfig,
		logger)
	genesisBlockDataDTORepoConfig := auth_repo.NewGenesisBlockDataDTOConfigurationProvider(flagAuthorityAddress)
	genesisBlockDataDTORepo := auth_repo.NewGenesisBlockDataDTORepo(
		genesisBlockDataDTORepoConfig,
		logger)
	blockDataRepo := repo.NewBlockDataRepo(
		logger,
		blockDataDB)
	blockDataDTORepoConfig := auth_repo.NewBlockDataDTOConfigurationProvider(flagAuthorityAddress)
	blockDataDTORepo := auth_repo.NewBlockDataDTORepo(
		blockDataDTORepoConfig,
		logger)
	tokRepo := repo.NewTokenRepo(
		logger,
		tokDB)
	blockchainStateChangeEventDTOConfigurationProvider := auth_repo.NewBlockchainStateChangeEventDTOConfigurationProvider(flagAuthorityAddress)
	blockchainStateChangeEventDTORepo := auth_repo.NewBlockchainStateChangeEventDTORepo(
		blockchainStateChangeEventDTOConfigurationProvider,
		logger)
	nftokenRepo := repo.NewNonFungibleTokenRepo(logger, nftokDB)
	nftAssetRepoConfig := repo.NewNFTAssetRepoConfigurationProvider(flagNFTStorageAddress, "")
	nftAssetRepo := repo.NewNFTAssetRepo(nftAssetRepoConfig, logger)
	mempoolTxDTORepoConfig := auth_repo.NewMempoolTransactionDTOConfigurationProvider(flagAuthorityAddress)
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
	_ = listAllWalletUseCase

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

	rpcServerConfigurationProvider := rpc.NewRPCServerConfigurationProvider("localhost", "2233")
	rpcServer := rpc.NewRPCServer(
		rpcServerConfigurationProvider,
		logger,
		getAccountService,
		createAccountService,
		accountListingByLocalWalletsService,
		coinTransferService,
		tokenGetService,
		tokenTransferService,
		tokenBurnService,
		getOrDownloadNonFungibleTokenService,
		listBlockTransactionsByAddressService,
		getByBlockTransactionTimestampService,
		blockDataGetByHashService,
		tokenListByOwnerService,
		exportWalletService,
		importWalletService,
	)

	//
	// Execute.
	//

	go func() {
		for {
			logger.Debug("Starting sync-manager...")
			if err := blockchainSyncManagerService.Execute(context.Background(), flagChainID); err != nil {
				logger.Error("Failed to manage syncing", slog.Any("error", err))
			}
			logger.Debug("Sync-manager will restart again.")
		}
	}()

	go func() {
		rpcServer.Run("localhost", "2233")
		defer rpcServer.Shutdown()
	}()

	logger.Info("ComicCoin CLI daemon is running.")
}
