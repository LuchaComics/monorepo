package daemon

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ethereum/go-ethereum/common"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/net/p2p"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/memory"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/http"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/http/middleware"
	taskmng "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	service "github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	usecase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run a blockchain node",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1:
			// Command line parameters validation.
			//

			if flagConsensusProtocol == constants.ConsensusPoA && flagEnableMiner {
				if flagProofOfAuthorityAccountAddress == "" {
					log.Fatal("Error loading ComicCoin: You did not enter the authority's address when running `proof of authority` consensus protocol miner, please use `--poa-address` flag next time.")
				}
				if flagProofOfAuthorityWalletPassword == "" {
					log.Fatal("Error loading ComicCoin: You did not enter the authority's password when running `proof of authority` consensus protocol miner, please use `--poa-password` flag next time.")
				}
			}

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
					ChainID:                        constants.ChainIDMainNet,
					TransPerBlock:                  1,
					Difficulty:                     2,
					ConsensusPollingDelayInMinutes: flagConsensusPollingDelayInMinutes,
					EnableMiner:                    flagEnableMiner,
					ConsensusProtocol:              flagConsensusProtocol,
				},
				App: config.AppConfig{
					DirPath:     flagDataDir,
					HTTPAddress: flagListenHTTPAddress,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
				Peer: config.PeerConfig{
					ListenPort:     flagListenPeerToPeerPort,
					KeyName:        flagKeypairName,
					BootstrapPeers: bootstrapPeers,
				},
			}
			if flagConsensusProtocol == constants.ConsensusPoA && flagEnableMiner {
				coinbaseAddr := common.HexToAddress(flagProofOfAuthorityAccountAddress)
				cfg.Blockchain.ProofOfAuthorityAccountAddress = &coinbaseAddr
				cfg.Blockchain.ProofOfAuthorityWalletPassword = flagProofOfAuthorityWalletPassword
			}

			logger := logger.NewLogger()
			walletDB := disk.NewDiskStorage(cfg.DB.DataDir, "wallet", logger)
			blockDataDB := disk.NewDiskStorage(cfg.DB.DataDir, "block_data", logger)
			latestHashDB := disk.NewDiskStorage(cfg.DB.DataDir, "latest_hash", logger)
			latestTokenIDDB := disk.NewDiskStorage(cfg.DB.DataDir, "latest_token_id", logger)
			ikDB := disk.NewDiskStorage(cfg.DB.DataDir, "identity_key", logger)
			pendingBlockDataDB := disk.NewDiskStorage(cfg.DB.DataDir, "pending_block_data", logger)
			mempoolTxDB := disk.NewDiskStorage(cfg.DB.DataDir, "mempool_tx", logger)
			tokenDB := disk.NewDiskStorage(cfg.DB.DataDir, "token", logger)
			memdb := memory.NewInMemoryStorage(logger)
			kmutex := kmutexutil.NewKMutexProvider()

			// ------------ Peer-to-Peer (P2P) ------------
			ikRepo := repo.NewIdentityKeyRepo(cfg, logger, ikDB)
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
				// CASE 1 OF 2: If custom identity key D.N.E.
				if flagIdentityKeyID != constants.DefaultIdentityKeyID {
					log.Fatal("Failed getting identity key: d.n.e.")
				}
				// CASE 2 OF 2: Create `DefaultIdentityKeyID`
				ik, err = ikCreateService.Execute(flagIdentityKeyID)
				if err != nil {
					log.Fatalf("Failed creating default identity key: %v", err)
				}

				// This is anomously behaviour so crash if this happens.
				if ik == nil {
					log.Fatal("Failed creating default identity key: d.n.e.")
				}
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

			// ------------ Repo ------------
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
			tokenRepo := repo.NewTokenRepo(
				cfg,
				logger,
				tokenDB)
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

			// ------------ Use-case ------------
			// Genesis Block Data
			loadGenesisBlockDataAccountUseCase := usecase.NewLoadGenesisBlockDataUseCase(
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
				tokenRepo)
			getTokensHashStateUseCase := usecase.NewGetTokensHashStateUseCase(
				cfg,
				logger,
				tokenRepo)
			getTokenUseCase := usecase.NewGetTokenUseCase(
				cfg,
				logger,
				tokenRepo)

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
			setBlockchainLastestTokenIDIfGreatestUseCase := usecase.NewSetBlockchainLastestTokenIDIfGreatestUseCase(
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

			// ------------ Service ------------

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
			_ = walletListService // We can use this elsewere...

			// Account
			initAccountsFromBlockchainService := service.NewInitAccountsFromBlockchainService(
				cfg,
				logger,
				loadGenesisBlockDataAccountUseCase,
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
			_ = getAccountBalanceService // TODO

			// Transaction
			createTxService := service.NewCreateTransactionService(
				cfg,
				logger,
				getAccountUseCase,
				getWalletUseCase,
				walletDecryptKeyUseCase,
				broadcastMempoolTxDTOUseCase)

			// Tokens
			poaTokenMintService := service.NewProofOfAuthorityTokenMintService(
				cfg,
				logger,
				kmutex,
				loadGenesisBlockDataAccountUseCase,
				getWalletUseCase,
				walletDecryptKeyUseCase,
				getBlockchainLastestTokenIDUseCase,
				broadcastMempoolTxDTOUseCase)

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
				setBlockchainLastestTokenIDIfGreatestUseCase,
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
				setBlockchainLastestTokenIDIfGreatestUseCase,
			)
			proofOfAuthorityValidationService := service.NewProofOfAuthorityValidationService(
				cfg,
				logger,
				kmutex,
				receiveProposedBlockDataDTOUseCase,
				getBlockchainLastestHashUseCase,
				getBlockDataUseCase,
				getAccountsHashStateUseCase,
				getTokensHashStateUseCase,
				createBlockDataUseCase,
				setBlockchainLastestHashUseCase,
				setBlockchainLastestTokenIDIfGreatestUseCase,
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
				loadGenesisBlockDataAccountUseCase,
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
			poaTokenMintHTTPHandler := httphandler.NewProofOfAuthorityTokenMintHTTPHandler(
				cfg,
				logger,
				poaTokenMintService)
			transferTokenHTTPHandler := httphandler.NewTransferTokenHTTPHandler(
				cfg,
				logger,
				transferTokenService)
			burnTokenHTTPHandler := httphandler.NewBurnTokenHTTPHandler(
				cfg,
				logger,
				burnTokenService)
			getTokenHTTPHandler := httphandler.NewGetTokenHTTPHandler(
				cfg,
				logger,
				getTokenService)
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
				poaTokenMintHTTPHandler,
				transferTokenHTTPHandler,
				burnTokenHTTPHandler,
				getTokenHTTPHandler,
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

			taskManager := taskmng.NewTaskManager(
				cfg,
				logger,
				tm1,
				tm2,
				tm3,
				tm4,
				tm5,
				tm6,
				tm7,
				tm8,
				tm9,
			)

			//
			// STEP 2
			// Perform whatever startup proceedures necessary to get our
			// blockchain ready for execution in our app.
			//

			if err := blockchainStartupService.Execute(); err != nil {
				log.Fatalf("failed blockchain startup: %v\n", err)
			}

			//
			// STEP 3
			// Run the main loop blocking code while other input ports run in
			// background.
			//

			logger.Info("Starting node...")

			// Run in background the peer to peer node which will synchronize our
			// blockchain with the network.
			// go peerNode.Run()

			if flagEnableHTTPJSONAPIServer {
				go httpServ.Run()
				defer httpServ.Shutdown()
			}
			go taskManager.Run()
			defer taskManager.Shutdown()

			logger.Info("Node running.")

			<-done
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", config.GetDefaultDataDirectory(), "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagListenHTTPAddress, "listen-http-address", "127.0.0.1:8000", "The IP and port to attach for our HTTP JSON API server")
	cmd.Flags().StringVar(&flagIdentityKeyID, "identitykey-id", "", "If you would like to use a custom identity then this is the identifier used to lookup a custom identity profile to assign for this blockchain node")
	cmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", config.GetDefaultPeerToPeerListenPort(), "The port to listen to for other peers")
	cmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")
	cmd.Flags().Int64Var(&flagConsensusPollingDelayInMinutes, "consensus-polling-delay-in-minutes", 1, "The delay interval between your node polling the network on the latest consensus")
	cmd.Flags().BoolVar(&flagEnableMiner, "enable-miner", false, "Controls whether you want your node to have a miner running in the background")
	cmd.Flags().StringVar(&flagConsensusProtocol, "consensus-protocol", "None", "Controls what consensus protocol to execute for the miner, choices are: PoW, PoA, or None.")
	cmd.Flags().StringVar(&flagProofOfAuthorityAccountAddress, "poa-address", "", "(Required for `PoA` consensus protocol) The address of the authority's account")
	cmd.Flags().StringVar(&flagProofOfAuthorityWalletPassword, "poa-password", "", "(Required for `PoA` consensus protocol) The password in the authority's wallet")
	cmd.Flags().BoolVar(&flagEnableHTTPJSONAPIServer, "enable-http-api", true, "Controls whether you want your node to have a HTTP JSON API server acessible by a GUI, or you are rnning only a node then this must be set to `false`")

	return cmd
}
