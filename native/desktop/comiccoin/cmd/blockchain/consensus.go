package blockchain

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config/constants"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/interface/task/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/memory"
)

func ConsensusMechanismCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "consensus",
		Short: "Run a node on the peer-to-peer network for the purpose contributing to blockchain consensus and data hosting.",
		Run: func(cmd *cobra.Command, args []string) {
			doBlockchainConsensusMechanism()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagIdentityKeyID, "identitykey-id", "", "If you would like to use a custom identity then this is the identifier used to lookup a custom identity profile to assign for this blockchain node.")
	cmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", 26642, "The port to listen to for other peers")
	cmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")

	return cmd
}

func doBlockchainConsensusMechanism() {
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

	// ------------ Peer-to-Peer (P2P) ------------
	ikRepo := repo.NewIdentityKeyRepo(cfg, logger, db)
	ikGetUseCase := usecase.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
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

	// ------------ Repo ------------

	latestBlockDataHashRepo := repo.NewBlockchainLastestHashRepo(
		cfg,
		logger,
		db)
	blockDataRepo := repo.NewBlockDataRepo(
		cfg,
		logger,
		db)
	blockDataDTORepo := repo.NewBlockDataDTORepo(
		cfg,
		logger,
		libP2PNetwork)
	accountRepo := repo.NewAccountRepo(
		cfg,
		logger,
		memdb) // Do not store on disk, only in-memory.
	consensusRepo := repo.NewConsensusRepoImpl(
		cfg,
		logger,
		libP2PNetwork)

	// ------------ Use-case ------------

	// Account
	createAccountUseCase := usecase.NewCreateAccountUseCase(
		cfg,
		logger,
		accountRepo)
	getAccountUseCase := usecase.NewGetAccountUseCase(
		cfg,
		logger,
		accountRepo)
	upsertAccountUseCase := usecase.NewUpsertAccountUseCase(
		cfg,
		logger,
		accountRepo)

	// Block Data
	getBlockDataUseCase := usecase.NewGetBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo)
	createBlockDataUseCase := usecase.NewCreateBlockDataUseCase(
		cfg,
		logger,
		blockDataRepo)

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

	// Latest BlockData Hash
	getBlockchainLastestHashUseCase := usecase.NewGetBlockchainLastestHashUseCase(
		cfg,
		logger,
		latestBlockDataHashRepo)
	setBlockchainLastestHashUseCase := usecase.NewSetBlockchainLastestHashUseCase(
		cfg,
		logger,
		latestBlockDataHashRepo)

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

	initAccountsFromBlockchainService := service.NewInitAccountsFromBlockchainService(
		cfg,
		logger,
		getBlockchainLastestHashUseCase,
		getBlockDataUseCase,
		getAccountUseCase,
		createAccountUseCase,
		upsertAccountUseCase,
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

	// ------------ Interface ------------

	// TASK MANAGER
	// tm6 := taskmnghandler.NewConsensusTaskHandler(
	// 	cfg,
	// 	logger,
	// 	consensusService)
	tm5 := taskmnghandler.NewBlockDataDTOServerTaskHandler(
		cfg,
		logger,
		uploadServerService)
	tm6 := taskmnghandler.NewMajorityVoteConsensusServerTaskHandler(
		cfg,
		logger,
		majorityVoteConsensusServerService)
	tm7 := taskmnghandler.NewMajorityVoteConsensusClientTaskHandler(
		cfg,
		logger,
		majorityVoteConsensusClientService)

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

	// ------------ Execution ------------

	//
	// STEP 3
	// Run the main loop blocking code while other input ports run in
	// background.
	//

	logger.Info("Starting network services only...")

	go func(server *taskmnghandler.BlockDataDTOServerTaskHandler) {
		ctx := context.Background()
		for {
			if err := server.Execute(ctx); err != nil {
				logger.Error("blockdatabto consensus server error",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			time.Sleep(5 * time.Second)
			logger.Debug("shared local blockchain with network")
			break
		}
	}(tm5)

	go func(server *taskmnghandler.MajorityVoteConsensusServerTaskHandler) {
		ctx := context.Background()
		for {
			if err := server.Execute(ctx); err != nil {
				logger.Error("consensus server error", slog.Any("error", err))
			}
			time.Sleep(5 * time.Second)
			logger.Debug("blockchain consensus serving done, restarting again ...")
		}
	}(tm6)

	go func(client *taskmnghandler.MajorityVoteConsensusClientTaskHandler) {
		ctx := context.Background()
		for {
			if err := client.Execute(ctx); err != nil {
				logger.Error("consensus client error", slog.Any("error", err))
			}
			time.Sleep(5 * time.Second)
			logger.Debug("blockchain consensus polling done, restarting again...")
		}
	}(tm7)

	<-done
}
