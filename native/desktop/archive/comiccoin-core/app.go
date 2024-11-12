package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/net/p2p"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/memory"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	task "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	service "github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	ma "github.com/multiformats/go-multiaddr"
)

// App struct
type App struct {
	ctx context.Context

	// Logger instance which provides detailed debugging information along
	// with the console log messages.
	logger *slog.Logger

	kmutex kmutexutil.KMutexProvider

	// Variable tells the application whether our app is connected to the
	// blockchain network as a node and executing successfully or not.
	isBlockchainNodeRunning bool

	// Variable controls the configuration of the blockchain node.
	config *config.Config

	// Variable holds access to our peer-to-peer network handler. We want to
	// have access so we can close our connection upon exit of this application.
	libP2PNetwork p2p.LibP2PNetwork

	getKeyService                                  *service.GetKeyService
	walletListService                              *service.WalletListService
	initAccountsFromBlockchainService              *service.InitAccountsFromBlockchainService
	createAccountService                           *service.CreateAccountService
	getAccountService                              *service.GetAccountService
	getAccountBalanceService                       *service.GetAccountBalanceService
	transferCoinService                            *service.TransferCoinService
	poaTokenMintService                            *service.ProofOfAuthorityTokenMintService
	transferTokenService                           *service.TransferTokenService
	burnTokenService                               *service.BurnTokenService
	getTokenService                                *service.GetTokenService
	listByOwnerTokenService                        *service.ListByOwnerTokenService
	countByOwnerTokenService                       *service.CountByOwnerTokenService
	mempoolReceiveService                          *service.MempoolReceiveService
	mempoolBatchSendService                        *service.MempoolBatchSendService
	proofOfWorkMiningService                       *service.ProofOfWorkMiningService
	proofOfAuthorityMiningService                  *service.ProofOfAuthorityMiningService
	proofOfWorkValidationService                   *service.ProofOfWorkValidationService
	proofOfAuthorityValidationService              *service.ProofOfAuthorityValidationService
	majorityVoteConsensusServerService             *service.MajorityVoteConsensusServerService
	majorityVoteConsensusClientService             *service.MajorityVoteConsensusClientService
	uploadServerService                            *service.BlockDataDTOServerService
	initBlockDataService                           *service.InitBlockDataService
	blockchainStartupService                       *service.BlockchainStartupService
	listRecentBlockTransactionService              *service.ListRecentBlockTransactionService
	listAllBlockTransactionService                 *service.ListAllBlockTransactionService
	getBlockDataByBlockTransactionTimestampService *service.GetBlockDataByBlockTransactionTimestampService
	getOrDownloadNonFungibleTokenService           *service.GetOrDownloadNonFungibleTokenService
	listNonFungibleTokensByOwnerService            *service.ListNonFungibleTokensByOwnerService
	signedIssuedTokenClientService                 *service.SignedIssuedTokenClientService
	mempoolReceiveTaskHandler                      *task.MempoolReceiveTaskHandler
	mempoolBatchSendTaskHandler                    *task.MempoolBatchSendTaskHandler
	proofOfWorkMiningTaskHandler                   *task.ProofOfWorkMiningTaskHandler
	proofOfAuthorityMiningTaskHandler              *task.ProofOfAuthorityMiningTaskHandler
	proofOfWorkValidationTaskHandler               *task.ProofOfWorkValidationTaskHandler
	proofOfAuthorityValidationTaskHandler          *task.ProofOfAuthorityValidationTaskHandler
	blockDataDTOServerTaskHandler                  *task.BlockDataDTOServerTaskHandler
	majorityVoteConsensusServerTaskHandler         *task.MajorityVoteConsensusServerTaskHandler
	majorityVoteConsensusClientTaskHandler         *task.MajorityVoteConsensusClientTaskHandler
	signedIssuedTokenClientServiceTaskHandler      *task.SignedIssuedTokenClientServiceTaskHandler
}

// NewApp creates a new App application struct
func NewApp() *App {
	logger := logger.NewProvider()
	kmutex := kmutexutil.NewKMutexProvider()

	cfg := &config.Config{}
	return &App{
		config:                  cfg,
		logger:                  logger,
		kmutex:                  kmutex,
		isBlockchainNodeRunning: false,
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
		a.logger.Debug("Startup halted: need to specify data directory")
		return
	}

	// DEVELOPERS NOTE:
	// Every ComicCoin node must be connected to a peer whom coordinates
	// connecting all the other nodes in the network, therefore we get the
	// following node(s) that act in this role.
	bootstrapPeers, err := config.StringToAddres(constants.ComicCoinBootstrapPeers)
	if err != nil {
		a.logger.Error("Startup aborted: failed converting string to multi-addresses",
			slog.Any("error", err))
		log.Fatalf("Failed converting string to multi-addresses: %v\n", err)
	}

	//
	// DEVELOPERS NOTE:
	// Load up our dependencies and configuration
	//

	// Initialize the configuration.
	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID:                        constants.ComicCoinChainID,
			TransPerBlock:                  constants.ComicCoinTransPerBlock,
			Difficulty:                     constants.ComicCoinDifficulty,
			ConsensusPollingDelayInMinutes: constants.ComicCoinConsensusPollingDelayInMinutes,
			ConsensusProtocol:              constants.ComicCoinConsensusProtocol,
		},
		App: config.AppConfig{
			DirPath: dataDir,
		},
		DB: config.DBConfig{
			DataDir: dataDir,
		},
		Peer: config.PeerConfig{
			ListenPort:     constants.ComicCoinPeerListenPort,
			KeyName:        constants.ComicCoinIdentityKeyID,
			BootstrapPeers: bootstrapPeers,
		},
		NFTAssetStore: config.NFTAssetStoreConfig{
			Address: constants.ComicCoinNFTAssetStoreAddress,
		},
	}
	a.config = cfg

	// For convinenence
	logger := a.logger
	kmutex := a.kmutex

	a.logger.Debug("Startup loading disk database...")
	walletDB := disk.NewDiskStorage(cfg.DB.DataDir, "wallet", logger)
	blockDataDB := disk.NewDiskStorage(cfg.DB.DataDir, "block_data", logger)
	latestHashDB := disk.NewDiskStorage(cfg.DB.DataDir, "latest_hash", logger)
	latestTokenIDDB := disk.NewDiskStorage(cfg.DB.DataDir, "latest_token_id", logger)
	ikDB := disk.NewDiskStorage(cfg.DB.DataDir, "identity_key", logger)
	pendingBlockDataDB := disk.NewDiskStorage(cfg.DB.DataDir, "pending_block_data", logger)
	mempoolTxDB := disk.NewDiskStorage(cfg.DB.DataDir, "mempool_tx", logger)
	tokDB := disk.NewDiskStorage(cfg.DB.DataDir, "token", logger)
	sitokenDB := disk.NewDiskStorage(cfg.DB.DataDir, "signed_issued_token", logger)
	nftokDB := disk.NewDiskStorage(cfg.DB.DataDir, "non_fungible_token", logger)
	memdb := memory.NewInMemoryStorage(logger)

	a.logger.Debug("Startup loading peer-to-peer client...")
	ikRepo := repo.NewIdentityKeyRepo(cfg, logger, ikDB)
	ikCreateUseCase := usecase.NewCreateIdentityKeyUseCase(cfg, logger, ikRepo)
	ikGetUseCase := usecase.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
	ikCreateService := service.NewCreateIdentityKeyService(cfg, logger, ikCreateUseCase, ikGetUseCase)
	ikGetService := service.NewGetIdentityKeyService(cfg, logger, ikGetUseCase)

	// Get our identity key.
	ik, err := ikGetService.Execute(constants.ComicCoinIdentityKeyID)
	if err != nil {
		log.Fatalf("Failed getting identity key: %v", err)
	}
	if ik == nil {
		ik, err = ikCreateService.Execute(constants.ComicCoinIdentityKeyID)
		if err != nil {
			log.Fatalf("Failed creating ComicCoin identity key: %v", err)
		}

		// This is anomously behaviour so crash if this happens.
		if ik == nil {
			log.Fatal("Failed creating ComicCoin identity key: d.n.e.")
		}
	}
	privateKey, _ := ik.GetPrivateKey()
	publicKey, _ := ik.GetPublicKey()
	libP2PNetwork := p2p.NewLibP2PNetwork(cfg, logger, privateKey, publicKey)
	h := libP2PNetwork.GetHost()

	// Save to our app.
	a.libP2PNetwork = libP2PNetwork

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := h.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	logger.Info("node ready",
		slog.Any("peer identity", h.ID()),
		slog.Any("full address", fullAddr),
	)

	//
	// Repositories
	//
	a.logger.Debug("Startup loading repositories...")

	walletRepo := repo.NewWalletRepo(
		cfg,
		logger,
		walletDB)
	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		cfg,
		logger,
		blockDataDB)
	accountRepo := repo.NewAccountRepo(
		cfg,
		logger,
		memdb) // Do not store on disk, only in-memory.
	tokRepo := repo.NewTokenRepo(
		cfg,
		logger,
		tokDB)
	sitokRepo := repo.NewSignedIssuedTokenRepo(
		cfg,
		logger,
		sitokenDB)
	sitokDTORepo := repo.NewSignedIssuedTokenDTORepo(
		cfg,
		logger,
		libP2PNetwork)
	nftokRepo := repo.NewNonFungibleTokenRepo(
		logger,
		nftokDB)
	nftAssetRepoCfg := repo.NewNFTAssetRepoConfigurationProvider(cfg.NFTAssetStore.Address, "")
	nftAssetRepo := repo.NewNFTAssetRepo(nftAssetRepoCfg, logger)
	mempoolTxRepo := repo.NewMempoolTransactionRepo(
		cfg,
		logger,
		mempoolTxDB)
	mempoolTransactionDTORepo := repo.NewMempoolTransactionDTORepo(
		cfg,
		logger,
		libP2PNetwork)
	pendingBlockTxRepo := repo.NewPendingBlockTransactionRepo(
		cfg,
		logger,
		pendingBlockDataDB)
	latestBlockDataHashRepo := repo.NewBlockchainLastestHashRepo(
		cfg,
		logger,
		latestHashDB)
	latestBlockDataTokenIDRepo := repo.NewBlockchainLastestTokenIDRepo(
		cfg,
		logger,
		latestTokenIDDB)
	blockDataRepo := repo.NewBlockDataRepo(
		cfg,
		logger,
		blockDataDB)
	proposedBlockDataDTORepo := repo.NewProposedBlockDataDTORepo(
		cfg,
		logger,
		libP2PNetwork)
	blockDataDTORepo := repo.NewBlockDataDTORepo(
		cfg,
		logger,
		libP2PNetwork)
	consensusRepo := repo.NewConsensusRepoImpl(
		cfg,
		logger,
		libP2PNetwork)

	//
	// USECASES
	//

	a.logger.Debug("Startup loading usecases...")

	// Storage Transcation
	storageTransactionOpenUseCase := usecase.NewStorageTransactionOpenUseCase(
		cfg,
		logger,
		accountRepo,
		tokRepo,
		latestBlockDataHashRepo,
		latestBlockDataTokenIDRepo,
		blockDataRepo,
		ikRepo,
		mempoolTxRepo,
		pendingBlockTxRepo,
		walletRepo,
		nftokRepo,
	)
	storageTransactionCommitUseCase := usecase.NewStorageTransactionCommitUseCase(
		cfg,
		logger,
		accountRepo,
		tokRepo,
		latestBlockDataHashRepo,
		latestBlockDataTokenIDRepo,
		blockDataRepo,
		ikRepo,
		mempoolTxRepo,
		pendingBlockTxRepo,
		walletRepo,
		nftokRepo,
	)
	storageTransactionDiscardUseCase := usecase.NewStorageTransactionDiscardUseCase(
		cfg,
		logger,
		accountRepo,
		tokRepo,
		latestBlockDataHashRepo,
		latestBlockDataTokenIDRepo,
		blockDataRepo,
		ikRepo,
		mempoolTxRepo,
		pendingBlockTxRepo,
		walletRepo,
		nftokRepo,
	)

	// Genesis Block Data
	loadGenesisBlockDataUseCase := usecase.NewLoadGenesisBlockDataUseCase(
		cfg,
		logger,
		genesisBlockDataRepo)

	// Wallet
	createWalletUseCase := usecase.NewCreateWalletUseCase(
		cfg,
		logger,
		walletRepo)
	walletDecryptKeyUseCase := usecase.NewWalletDecryptKeyUseCase(
		cfg,
		logger,
		walletRepo)
	walletEncryptKeyUseCase := usecase.NewWalletEncryptKeyUseCase(
		cfg,
		logger,
		walletRepo)
	getWalletUseCase := usecase.NewGetWalletUseCase(
		cfg,
		logger,
		walletRepo)
	listAllWalletUseCase := usecase.NewListAllWalletUseCase(
		cfg,
		logger,
		walletRepo)

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		cfg,
		logger,
		accountRepo)
	getAccountsHashStateUseCase := usecase.NewGetAccountsHashStateUseCase(
		cfg,
		logger,
		accountRepo)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo)

	// Token
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokRepo)
	getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
		cfg,
		logger,
		tokRepo)
	getTokenUseCase := usecase.NewGetTokenUseCase(
		cfg,
		logger,
		tokRepo)
	listTokensByOwnerUseCase := usecase.NewListTokensByOwnerUseCase(
		cfg,
		logger,
		tokRepo)
	countTokensByOwnerUseCase := usecase.NewCountTokensByOwnerUseCase(
		cfg,
		logger,
		tokRepo)

	// Non-Fungible Token
	getTokUseCase := usecase.NewGetTokenUseCase(
		cfg,
		logger,
		tokRepo)

	getNFTokUseCase := usecase.NewGetNonFungibleTokenUseCase(
		cfg,
		logger,
		nftokRepo)

	downloadNFTokMetadataUsecase := usecase.NewDownloadMetadataNonFungibleTokenUseCase(
		cfg,
		logger,
		nftAssetRepo)

	downloadNFTokAssetUsecase := usecase.NewDownloadNonFungibleTokenAssetUseCase(
		cfg,
		logger,
		nftAssetRepo)
	listNFTsWithFilterByTokenIDsyUseCase := usecase.NewListNonFungibleTokensWithFilterByTokenIDsyUseCase(
		cfg,
		logger,
		nftokRepo)
	upsertNFTokUseCase := usecase.NewUpsertNonFungibleTokenUseCase(
		cfg,
		logger,
		nftokRepo)

	// Signed Issued Token DTO
	createSignedIssuedTokenUseCase := usecase.NewCreateSignedIssuedTokenUseCase(
		cfg,
		logger,
		sitokRepo)
	broadcastSignedIssuedTokenDTOUseCase := usecase.NewBroadcastSignedIssuedTokenDTOUseCase(
		cfg,
		logger,
		sitokDTORepo)
	receiveSignedIssuedTokenDTOUseCase := usecase.NewReceiveSignedIssuedTokenDTOUseCase(
		cfg,
		logger,
		sitokDTORepo)

	// Mempool Transaction DTO
	broadcastMempoolTxDTOUseCase := usecase.NewBroadcastMempoolTransactionDTOUseCase(
		cfg,
		logger,
		mempoolTransactionDTORepo)
	receiveMempoolTxDTOUseCase := usecase.NewReceiveMempoolTransactionDTOUseCase(
		cfg,
		logger,
		mempoolTransactionDTORepo)

	// Mempool Transaction
	createMempoolTransactionUseCase := usecase.NewCreateMempoolTransactionUseCase(
		cfg,
		logger,
		mempoolTxRepo)
	listAllMempoolTransactionUseCase := usecase.NewListAllMempoolTransactionUseCase(
		cfg,
		logger,
		mempoolTxRepo)
	deleteAllMempoolTransactionUseCase := usecase.NewDeleteAllMempoolTransactionUseCase(
		cfg,
		logger,
		mempoolTxRepo)

	// Proposed Block Transaction
	createPendingBlockTxUseCase := usecase.NewCreatePendingBlockTransactionUseCase(
		cfg,
		logger,
		pendingBlockTxRepo)
	listAllPendingBlockTxUseCase := usecase.NewListAllPendingBlockTransactionUseCase(
		cfg,
		logger,
		pendingBlockTxRepo)
	deleteAllPendingBlockTxUseCase := usecase.NewDeleteAllPendingBlockTransactionUseCase(
		cfg,
		logger,
		pendingBlockTxRepo)

	// Latest BlockData Hash
	getBlockchainLastestHashUseCase := usecase.NewGetBlockchainLastestHashUseCase(
		cfg,
		logger,
		latestBlockDataHashRepo)
	setBlockchainLastestHashUseCase := usecase.NewSetBlockchainLastestHashUseCase(
		cfg,
		logger,
		latestBlockDataHashRepo)

	// Latest BlockData Token ID
	getBlockchainLastestTokenIDUseCase := usecase.NewGetBlockchainLastestTokenIDUseCase(
		cfg,
		logger,
		latestBlockDataTokenIDRepo)
	setBlockchainLastestTokenIDIfGTUseCase := usecase.NewSetBlockchainLastestTokenIDIfGTUseCase(
		cfg,
		logger,
		latestBlockDataTokenIDRepo)

	// Block Data
	getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo)
	createBlockDataUseCase := usecase.NewCreateBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo)

	// Block Transactions (via Block Data).
	listAllBlockTransactionByAddressUseCase := usecase.NewListAllBlockTransactionByAddressUseCase(
		cfg,
		logger,
		blockDataRepo)
	getBlockDataByBlockTransactionTimestampUseCase := usecase.NewGetBlockDataByBlockTransactionTimestampUseCase(
		cfg,
		logger,
		blockDataRepo)

	// Mining
	proofOfWorkUseCase := usecase.NewProofOfWorkUseCase(cfg, logger)

	// Proposed Block Data DTO
	broadcastProposedBlockDataDTOUseCase := usecase.NewBroadcastProposedBlockDataDTOUseCase(
		cfg,
		logger,
		proposedBlockDataDTORepo)
	receiveProposedBlockDataDTOUseCase := usecase.NewReceiveProposedBlockDataDTOUseCase(
		cfg,
		logger,
		proposedBlockDataDTORepo)

	// Block Data DTO
	blockDataDTOReceiveP2PResponseUseCase := usecase.NewBlockDataDTOReceiveP2PResponsetUseCase(
		cfg,
		logger,
		blockDataDTORepo)
	blockDataDTOReceiveP2PRequesttUseCase := usecase.NewBlockDataDTOReceiveP2PRequesttUseCase(
		cfg,
		logger,
		blockDataDTORepo)
	blockDataDTOSendP2PResponsetUseCase := usecase.NewBlockDataDTOSendP2PResponsetUseCase(
		cfg,
		logger,
		blockDataDTORepo)
	blockDataDTOSendP2PRequestUseCase := usecase.NewBlockDataDTOSendP2PRequestUseCase(
		cfg,
		logger,
		blockDataDTORepo)

	// Consensus Mechanism
	consensusMechanismBroadcastRequestToNetworkUseCase := usecase.NewConsensusMechanismBroadcastRequestToNetworkUseCase(
		cfg,
		logger,
		consensusRepo)
	consensusMechanismReceiveRequestFromNetworkUseCase := usecase.NewConsensusMechanismReceiveRequestFromNetworkUseCase(
		cfg,
		logger,
		consensusRepo)
	consensusMechanismSendResponseToPeerUseCase := usecase.NewConsensusMechanismSendResponseToPeerUseCase(
		cfg,
		logger,
		consensusRepo)
	consensusMechanismReceiveResponseFromNetworkUseCase := usecase.NewConsensusMechanismReceiveResponseFromNetworkUseCase(
		cfg,
		logger,
		consensusRepo)

	//
	// SERVICES
	//

	a.logger.Debug("Startup loading services...")

	// Wallet + Key service
	getKeyService := service.NewGetKeyService(
		cfg,
		logger,
		getWalletUseCase,
		walletDecryptKeyUseCase)
	walletListService := service.NewWalletListService(
		cfg,
		logger,
		listAllWalletUseCase)

	// Account
	initAccountsFromBlockchainService := service.NewInitAccountsFromBlockchainService(
		cfg,
		logger,
		loadGenesisBlockDataUseCase,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
		getAccountUseCase,
		getAccountsHashStateUseCase,
		createAccountUseCase,
		upsertAccountUseCase)
	createAccountService := service.NewCreateAccountService(
		cfg,
		logger,
		walletEncryptKeyUseCase,
		walletDecryptKeyUseCase,
		createWalletUseCase,
		createAccountUseCase,
		getAccountUseCase)
	getAccountService := service.NewGetAccountService(
		cfg,
		logger,
		getAccountUseCase,
		getWalletUseCase,
		createAccountUseCase,
	)
	getAccountBalanceService := service.NewGetAccountBalanceService(
		cfg,
		logger,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase)

	// Transaction
	transferCoinService := service.NewTransferCoinService(
		cfg,
		logger,
		getAccountUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		broadcastMempoolTxDTOUseCase)

	// Block Transaction
	listRecentBlockTransactionService := service.NewListRecentBlockTransactionService(
		cfg,
		logger,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
	)
	listAllBlockTransactionService := service.NewListAllBlockTransactionService(
		cfg,
		logger,
		listAllBlockTransactionByAddressUseCase,
	)
	getBlockDataByBlockTransactionTimestampService := service.NewGetBlockDataByBlockTransactionTimestampService(
		cfg,
		logger,
		getBlockDataByBlockTransactionTimestampUseCase,
	)

	// Tokens
	poaTokenMintService := service.NewProofOfAuthorityTokenMintService(
		cfg,
		logger,
		kmutex,
		loadGenesisBlockDataUseCase,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getBlockchainLastestTokenIDUseCase,
		broadcastMempoolTxDTOUseCase,
		broadcastSignedIssuedTokenDTOUseCase)

	transferTokenService := service.NewTransferTokenService(
		cfg,
		logger,
		kmutex,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getTokenUseCase,
		broadcastMempoolTxDTOUseCase)
	burnTokenService := service.NewBurnTokenService(
		cfg,
		logger,
		kmutex,
		getWalletUseCase,
		walletDecryptKeyUseCase,
		getTokenUseCase,
		broadcastMempoolTxDTOUseCase)
	getTokenService := service.NewGetTokenService(
		cfg,
		logger,
		getTokenUseCase)
	listByOwnerTokenService := service.NewListByOwnerTokenService(
		cfg,
		logger,
		listTokensByOwnerUseCase)
	countByOwnerTokenService := service.NewCountByOwnerTokenService(
		cfg,
		logger,
		countTokensByOwnerUseCase)

	// Non-Fungible Tokens
	getOrDownloadNonFungibleTokenService := service.NewGetOrDownloadNonFungibleTokenService(
		cfg,
		logger,
		getNFTokUseCase,
		getTokUseCase,
		downloadNFTokMetadataUsecase,
		downloadNFTokAssetUsecase,
		upsertNFTokUseCase)
	listNonFungibleTokensByOwnerService := service.NewListNonFungibleTokensByOwnerService(
		cfg,
		logger,
		listTokensByOwnerUseCase,
		listNFTsWithFilterByTokenIDsyUseCase,
		getOrDownloadNonFungibleTokenService)

	// Mempool
	mempoolReceiveService := service.NewMempoolReceiveService(
		cfg,
		logger,
		kmutex,
		receiveMempoolTxDTOUseCase,
		createMempoolTransactionUseCase)
	mempoolBatchSendService := service.NewMempoolBatchSendService(
		cfg,
		logger,
		kmutex,
		listAllMempoolTransactionUseCase,
		createPendingBlockTxUseCase,
		deleteAllMempoolTransactionUseCase)

	// Mining
	proofOfWorkMiningService := service.NewProofOfWorkMiningService(
		cfg,
		logger,
		kmutex,
		getAccountsHashStateUseCase,
		listAllPendingBlockTxUseCase,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
		proofOfWorkUseCase,
		broadcastProposedBlockDataDTOUseCase,
		deleteAllPendingBlockTxUseCase,
	)

	proofOfAuthorityMiningService := service.NewProofOfAuthorityMiningService(
		cfg,
		logger,
		kmutex,
		getKeyService,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		getAccountUseCase,
		getAccountsHashStateUseCase,
		getTokenUseCase,
		getTokensHashStateUseCase,
		listAllPendingBlockTxUseCase,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
		proofOfWorkUseCase,
		createBlockDataUseCase,
		broadcastProposedBlockDataDTOUseCase,
		deleteAllPendingBlockTxUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
		upsertAccountUseCase,
		setBlockchainLastestHashUseCase,
		getBlockchainLastestTokenIDUseCase,
		setBlockchainLastestTokenIDIfGTUseCase,
	)

	// Validation
	proofOfWorkValidationService := service.NewProofOfWorkValidationService(
		cfg,
		logger,
		kmutex,
		receiveProposedBlockDataDTOUseCase,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
		getAccountsHashStateUseCase,
		createBlockDataUseCase,
		setBlockchainLastestHashUseCase,
		getAccountUseCase,
		upsertAccountUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
		getBlockchainLastestTokenIDUseCase,
		setBlockchainLastestTokenIDIfGTUseCase,
	)
	proofOfAuthorityValidationService := service.NewProofOfAuthorityValidationService(
		cfg,
		logger,
		kmutex,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		receiveProposedBlockDataDTOUseCase,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
		getAccountsHashStateUseCase,
		getTokensHashStateUseCase,
		loadGenesisBlockDataUseCase,
		createBlockDataUseCase,
		setBlockchainLastestHashUseCase,
		setBlockchainLastestTokenIDIfGTUseCase,
		getAccountUseCase,
		upsertAccountUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
	)

	majorityVoteConsensusServerService := service.NewMajorityVoteConsensusServerService(
		cfg,
		logger,
		consensusMechanismReceiveRequestFromNetworkUseCase,
		getBlockchainLastestHashUseCase,
		consensusMechanismSendResponseToPeerUseCase,
	)
	majorityVoteConsensusClientService := service.NewMajorityVoteConsensusClientService(
		cfg,
		logger,
		storageTransactionOpenUseCase,
		storageTransactionCommitUseCase,
		storageTransactionDiscardUseCase,
		consensusMechanismBroadcastRequestToNetworkUseCase,
		consensusMechanismReceiveResponseFromNetworkUseCase,
		getBlockchainLastestHashUseCase,
		setBlockchainLastestHashUseCase,
		blockDataDTOSendP2PRequestUseCase,
		blockDataDTOReceiveP2PResponseUseCase,
		createBlockDataUseCase,
		getBlockDataUseCase,
		getAccountUseCase,
		upsertAccountUseCase,
		getAccountsHashStateUseCase,
		getTokensHashStateUseCase,
	)
	uploadServerService := service.NewBlockDataDTOServerService(
		cfg,
		logger,
		blockDataDTOReceiveP2PRequesttUseCase,
		getBlockDataUseCase,
		blockDataDTOSendP2PResponsetUseCase,
	)
	initBlockDataService := service.NewInitBlockDataService(
		cfg,
		logger,
		loadGenesisBlockDataUseCase,
		getBlockDataUseCase,
		createBlockDataUseCase,
		setBlockchainLastestHashUseCase,
	)
	blockchainStartupService := service.NewBlockchainStartupService(
		cfg,
		logger,
		initAccountsFromBlockchainService,
		initBlockDataService,
	)
	signedIssuedTokenClientService := service.NewSignedIssuedTokenClientService(
		cfg,
		logger,
		kmutex,
		receiveSignedIssuedTokenDTOUseCase,
		loadGenesisBlockDataUseCase,
		createSignedIssuedTokenUseCase,
	)

	// Save the services to our application.
	a.getKeyService = getKeyService
	a.walletListService = walletListService
	a.initAccountsFromBlockchainService = initAccountsFromBlockchainService
	a.createAccountService = createAccountService
	a.getAccountService = getAccountService
	a.getAccountBalanceService = getAccountBalanceService
	a.transferCoinService = transferCoinService
	a.listRecentBlockTransactionService = listRecentBlockTransactionService
	a.listAllBlockTransactionService = listAllBlockTransactionService
	a.getOrDownloadNonFungibleTokenService = getOrDownloadNonFungibleTokenService
	a.listNonFungibleTokensByOwnerService = listNonFungibleTokensByOwnerService
	a.getBlockDataByBlockTransactionTimestampService = getBlockDataByBlockTransactionTimestampService
	a.poaTokenMintService = poaTokenMintService
	a.transferTokenService = transferTokenService
	a.burnTokenService = burnTokenService
	a.getTokenService = getTokenService
	a.listByOwnerTokenService = listByOwnerTokenService
	a.countByOwnerTokenService = countByOwnerTokenService
	a.mempoolReceiveService = mempoolReceiveService
	a.mempoolBatchSendService = mempoolBatchSendService
	a.proofOfWorkMiningService = proofOfWorkMiningService
	a.proofOfAuthorityMiningService = proofOfAuthorityMiningService
	a.proofOfWorkValidationService = proofOfWorkValidationService
	a.proofOfAuthorityValidationService = proofOfAuthorityValidationService
	a.majorityVoteConsensusServerService = majorityVoteConsensusServerService
	a.majorityVoteConsensusClientService = majorityVoteConsensusClientService
	a.uploadServerService = uploadServerService
	a.initBlockDataService = initBlockDataService
	a.blockchainStartupService = blockchainStartupService
	a.signedIssuedTokenClientService = signedIssuedTokenClientService

	//
	// BACKGROUND TASKS
	//

	a.logger.Debug("Startup loading background tasks...")

	// TASK MANAGER
	tm1 := taskmnghandler.NewMempoolReceiveTaskHandler(
		cfg,
		logger,
		mempoolReceiveService)
	tm2 := taskmnghandler.NewMempoolBatchSendTaskHandler(
		cfg,
		logger,
		mempoolBatchSendService)
	tm3 := taskmnghandler.NewProofOfWorkMiningTaskHandler(
		cfg,
		logger,
		proofOfWorkMiningService)
	tm4 := taskmnghandler.NewProofOfAuthorityMiningTaskHandler(
		cfg,
		logger,
		proofOfAuthorityMiningService)
	tm5 := taskmnghandler.NewProofOfWorkValidationTaskHandler(
		cfg,
		logger,
		proofOfWorkValidationService)
	tm6 := taskmnghandler.NewProofOfAuthorityValidationTaskHandler(
		cfg,
		logger,
		proofOfAuthorityValidationService)
	tm7 := taskmnghandler.NewBlockDataDTOServerTaskHandler(
		cfg,
		logger,
		uploadServerService)
	tm8 := taskmnghandler.NewMajorityVoteConsensusServerTaskHandler(
		cfg,
		logger,
		majorityVoteConsensusServerService)
	tm9 := taskmnghandler.NewMajorityVoteConsensusClientTaskHandler(
		cfg,
		logger,
		majorityVoteConsensusClientService)
	tm10 := taskmnghandler.NewSignedIssuedTokenClientServiceTaskHandler(
		cfg,
		logger,
		signedIssuedTokenClientService)

	// Save the services to our application.
	a.mempoolReceiveTaskHandler = tm1
	a.mempoolBatchSendTaskHandler = tm2
	a.proofOfWorkMiningTaskHandler = tm3
	a.proofOfAuthorityMiningTaskHandler = tm4
	a.proofOfWorkValidationTaskHandler = tm5
	a.proofOfAuthorityValidationTaskHandler = tm6
	a.blockDataDTOServerTaskHandler = tm7
	a.majorityVoteConsensusServerTaskHandler = tm8
	a.majorityVoteConsensusClientTaskHandler = tm9
	a.signedIssuedTokenClientServiceTaskHandler = tm10

	//
	// STEP 2
	// Perform whatever startup proceedures necessary to get our
	// blockchain ready for execution in our app.
	//
	a.logger.Debug("Startup loading blockchain...")

	if err := blockchainStartupService.Execute(); err != nil {
		log.Fatalf("failed blockchain startup: %v\n", err)
	}

	a.logger.Debug("Startup finished")
	a.isBlockchainNodeRunning = true
	go a.startBackgroundTasks()
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

	a.isBlockchainNodeRunning = false
	go a.stopBackgroundTasks()

	a.logger.Debug("Peer-to-peer network shutting down...")
	a.libP2PNetwork.Close()
}