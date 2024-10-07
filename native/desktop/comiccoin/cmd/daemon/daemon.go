package daemon

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config/constants"
	http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/http/middleware"
	taskmng "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	service "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/memory"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run a blockchain node",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies and configuration
			//

			// Load up our operating system interaction handlers, more specifically
			// signals. The OS sends our application various signals based on the
			// OS's state, we want to listen into the termination signals.
			done := make(chan os.Signal, 1)
			signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

			bootstrapPeers, err := StringToAddres(flagBootstrapPeers)
			if err != nil {
				log.Fatalf("Failed converting string to multi-addresses: %v\n", err)
			}

			cfg := &config.Config{
				Blockchain: config.BlockchainConfig{
					ChainID:       constants.ChainIDMainNet,
					TransPerBlock: 1,
					Difficulty:    2,
				},
				App: config.AppConfig{
					DirPath:     flagDataDir,
					HTTPAddress: flagListenHTTPAddress,
					RPCAddress:  flagListenRPCAddress,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
				Peer: config.PeerConfig{
					ListenPort:       flagListenPeerToPeerPort,
					KeyName:          flagKeypairName,
					RendezvousString: flagRendezvousString,
					BootstrapPeers:   bootstrapPeers,
				},
			}
			logger := logger.NewLogger()
			db := disk.NewDiskStorage(cfg.DB.DataDir, logger)
			memdb := memory.NewInMemoryStorage(logger)
			kmutex := kmutexutil.NewKMutexProvider()

			// ------------ Peer-to-Peer (P2P) ------------
			ikRepo := repo.NewIdentityKeyRepo(cfg, logger, db)
			ikCreateUseCase := usecase.NewCreateIdentityKeyUseCase(cfg, logger, ikRepo)
			ikGetUseCase := usecase.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
			ikCreateService := service.NewCreateIdentityKeyService(cfg, logger, ikCreateUseCase, ikGetUseCase)
			ikGetService := service.NewGetIdentityKeyService(cfg, logger, ikGetUseCase)

			// If nothing was set then we use a default value. We do this to
			// simplify the user's experience.
			if flagIdentityKeyID == "" {
				flagIdentityKeyID = constants.DefaultIdentityKeyID
			}

			// Get our identity key.
			ik, err := ikGetService.Execute(flagIdentityKeyID)
			if err != nil {
				log.Fatalf("Failed getting identity key: %v", err)
			}
			if ik == nil {
				log.Fatal("Failed getting identity key: d.n.e.")
			}
			logger.Debug("Identity key found")

			privateKey, _ := ik.GetPrivateKey()
			publicKey, _ := ik.GetPublicKey()
			libP2PNetwork := p2p.NewLibP2PNetwork(cfg, logger, privateKey, publicKey)
			defer libP2PNetwork.Close()
			h := libP2PNetwork.GetHost()

			// Build host multiaddress
			hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID()))

			// Now we can build a full multiaddress to reach this host
			// by encapsulating both addresses:
			addr := h.Addrs()[0]
			fullAddr := addr.Encapsulate(hostAddr)

			logger.Info("Blockchain node ready",
				slog.Any("peer identity", h.ID()),
				slog.Any("full address", fullAddr),
			)

			//TODO
			_ = ikCreateService
			_ = ikGetService

			// ------------ Repo ------------
			walletRepo := repo.NewWalletRepo(
				cfg,
				logger,
				db)
			accountRepo := repo.NewAccountRepo(
				cfg,
				logger,
				memdb) // Do not store on disk, only in-memory.
			mempoolTxRepo := repo.NewMempoolTransactionRepo(
				cfg,
				logger,
				db)
			mempoolTransactionDTORepo := repo.NewMempoolTransactionDTORepo(
				cfg,
				logger,
				libP2PNetwork)
			pendingBlockTxRepo := repo.NewPendingBlockTransactionRepo(
				cfg,
				logger,
				db)
			latestBlockDataHashRepo := repo.NewBlockchainLastestHashRepo(
				cfg,
				logger,
				db)
			blockDataRepo := repo.NewBlockDataRepo(
				cfg,
				logger,
				db)
			proposedBlockDataDTORepo := repo.NewProposedBlockDataDTORepo(
				cfg,
				logger,
				libP2PNetwork)
			blockDataDTORepo := repo.NewBlockDataDTORepo(
				cfg,
				logger,
				libP2PNetwork)
			consensusRepo := repo.NewMajorityVoteConsensusRepoImpl(
				cfg,
				logger,
				libP2PNetwork)

			// ------------ Use-case ------------
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

			// Block Data
			getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
				cfg,
				logger,
				blockDataRepo)
			createBlockDataUseCase := usecase.NewCreateBlockDataUseCase(
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
			castVoteForLatestHashInMajorityVotingConsensusUseCase := usecase.NewCastVoteForLatestHashInMajorityVotingConsensusUseCase(
				cfg,
				logger,
				consensusRepo)
			queryLatestHashByMajorityVotingConsensusUseCase := usecase.NewQueryLatestHashByMajorityVotingConsensusUseCase(
				cfg,
				logger,
				consensusRepo)

			// ------------ Service ------------
			// Account
			initAccountsFromBlockchainService := service.NewInitAccountsFromBlockchainService(
				cfg,
				logger,
				getBlockchainLastestHashUseCase,
				getBlockDataUseCase,
				getAccountUseCase,
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
				getAccountUseCase)
			getAccountBalanceService := service.NewGetAccountBalanceService(
				cfg,
				logger,
				getBlockchainLastestHashUseCase,
				getBlockDataUseCase)
			_ = getAccountBalanceService // TODO

			// Key
			getKeyService := service.NewGetKeyService(
				cfg,
				logger,
				getWalletUseCase,
				walletDecryptKeyUseCase)
			_ = getKeyService // TODO: USE IN FUTURE

			// Transaction
			createTxService := service.NewCreateTransactionService(
				cfg,
				logger,
				getAccountUseCase,
				getWalletUseCase,
				walletDecryptKeyUseCase,
				broadcastMempoolTxDTOUseCase)

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

			// Miner
			miningService := service.NewMiningService(
				cfg,
				logger,
				kmutex,
				getAccountsHashStateUseCase,
				listAllPendingBlockTxUseCase,
				getBlockchainLastestHashUseCase,
				setBlockchainLastestHashUseCase,
				getBlockDataUseCase,
				createBlockDataUseCase,
				proofOfWorkUseCase,
				broadcastProposedBlockDataDTOUseCase,
				deleteAllPendingBlockTxUseCase,
				getAccountUseCase,
				upsertAccountUseCase,
			)

			// Validation
			validationService := service.NewValidationService(
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
			)

			consensusService := service.NewConsensusService(
				cfg,
				logger,
				queryLatestHashByMajorityVotingConsensusUseCase,
				castVoteForLatestHashInMajorityVotingConsensusUseCase,
				getBlockchainLastestHashUseCase,
				setBlockchainLastestHashUseCase,
				blockDataDTOSendP2PRequestUseCase,
				blockDataDTOReceiveP2PResponseUseCase,
				createBlockDataUseCase,
				getBlockDataUseCase,
				getAccountUseCase,
				upsertAccountUseCase,
			)

			uploadServerService := service.NewBlockDataDTOServerService(
				cfg,
				logger,
				blockDataDTOReceiveP2PRequesttUseCase,
				getBlockDataUseCase,
				blockDataDTOSendP2PResponsetUseCase,
			)

			// ------------ Interface ------------
			// HTTP
			createAccountHTTPHandler := httphandler.NewCreateAccountHTTPHandler(
				cfg,
				logger,
				createAccountService)
			getAccountHTTPHandler := httphandler.NewGetAccountHTTPHandler(
				cfg,
				logger,
				getAccountService)
			createTransactionHTTPHandler := httphandler.NewCreateTransactionHTTPHandler(
				cfg,
				logger,
				createTxService)
			httpMiddleware := httpmiddle.NewMiddleware(
				cfg,
				logger)
			httpServ := http.NewHTTPServer(
				cfg,
				logger,
				httpMiddleware,
				createAccountHTTPHandler,
				getAccountHTTPHandler,
				createTransactionHTTPHandler,
			)

			// TASK MANAGER
			tm1 := taskmnghandler.NewMempoolReceiveTaskHandler(
				cfg,
				logger,
				mempoolReceiveService)
			tm2 := taskmnghandler.NewMempoolBatchSendTaskHandler(
				cfg,
				logger,
				mempoolBatchSendService)
			tm3 := taskmnghandler.NewMiningTaskHandler(
				cfg,
				logger,
				miningService)
			tm4 := taskmnghandler.NewValidationTaskHandler(
				cfg,
				logger,
				validationService)
			tm5 := taskmnghandler.NewConsensusTaskHandler(
				cfg,
				logger,
				consensusService)
			tm6 := taskmnghandler.NewBlockDataDTOServerTaskHandler(
				cfg,
				logger,
				uploadServerService)

			taskManager := taskmng.NewTaskManager(
				cfg,
				logger,
				tm1,
				tm2,
				tm3,
				tm4,
				tm5,
				tm6,
			)

			//
			// STEP 2
			// Load up the accounts into the in-memory storage before loading
			// the application because our accounts are only stored in memory
			// and not on disk.
			//

			logger.Info("Loading accounts into memory...")
			if err := initAccountsFromBlockchainService.Execute(); err != nil {
				log.Fatalf("failed executing accounts initialization: %v\n", err)
			}
			logger.Info("Accounts ready.")

			//
			// STEP 3
			// Run the main loop blocking code while other input ports run in
			// background.
			//

			logger.Info("Starting node...")

			// Run in background the peer to peer node which will synchronize our
			// blockchain with the network.
			// go peerNode.Run()

			go httpServ.Run()
			go taskManager.Run()
			defer httpServ.Shutdown()
			defer taskManager.Shutdown()

			logger.Info("Node running.")

			<-done
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagListenHTTPAddress, "listen-http-address", "127.0.0.1:8000", "The IP and port to attach for our HTTP JSON API server")
	cmd.Flags().StringVar(&flagIdentityKeyID, "identitykey-id", "", "If you would like to use a custom identity then this is the identifier used to lookup a custom identity profile to assign for this blockchain node.")
	cmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", 26642, "The port to listen to for other peers")
	cmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")

	return cmd
}
