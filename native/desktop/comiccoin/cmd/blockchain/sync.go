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
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

func SyncCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "Run a node on the peer-to-peer network for the purpose only sharing blockchain.",
		Run: func(cmd *cobra.Command, args []string) {
			doBlockchainSync()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagIdentityKeyID, "identitykey-id", "", "If you would like to use a custom identity then this is the identifier used to lookup a custom identity profile to assign for this blockchain node.")
	cmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", 26642, "The port to listen to for other peers")
	cmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")

	return cmd
}

func doBlockchainSync() {
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
	db := dbase.NewDatabase(cfg.DB.DataDir, logger)

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
	blockDataRepo := repo.NewBlockDataRepo(
		cfg,
		logger,
		db)
	_ = blockDataRepo

	blockchainSyncRepo := repo.NewBlockchainSyncRepo(
		cfg,
		logger,
		libP2PNetwork)

	go func() {
		ctx := context.Background()
		for {
			log.Println("sending...")
			if err := blockchainSyncRepo.SendRequestToRandomPeer(ctx, "test 1 2 3"); err != nil {
				log.Println("err:", err)
			}
			time.Sleep(5 * time.Second)
		}
	}()
	go func() {
		ctx := context.Background()
		for {
			peerID, hash, err := blockchainSyncRepo.ReceiveRequestFromNetwork(ctx)
			if err != nil {
				log.Println("err:", err)
			}
			log.Println("peerID:", peerID)
			log.Println("hash:", hash)
			time.Sleep(5 * time.Second)
		}
	}()

	// ------------ Use-case ------------

	// // Blockchain Sync'ing
	// uc1 := usecase.NewBlockchainSyncSendRequestUseCase(
	// 	cfg,
	// 	logger,
	// 	libP2PNetwork)
	// uc2 := usecase.NewBlockchainSyncReceiveRequestUseCase(
	// 	cfg,
	// 	logger,
	// 	libP2PNetwork)
	// uc3 := usecase.NewBlockchainSyncSendResponseUseCase(
	// 	cfg,
	// 	logger,
	// 	libP2PNetwork,
	// 	blockDataRepo)
	// uc4 := usecase.NewBlockchainSyncReceiveResponseUseCase(
	// 	cfg,
	// 	logger,
	// 	libP2PNetwork,
	// 	blockDataRepo)
	//
	// // ------------ Service ------------
	//
	// // Blockchain Sync'ing
	// s1 := service.NewBlockchainSyncServerService(
	// 	cfg,
	// 	logger,
	// 	uc2,
	// 	uc3)
	//
	// s2 := service.NewBlockchainSyncClientService(
	// 	cfg,
	// 	logger,
	// 	uc1,
	// 	uc4)
	//
	// // ------------ Interface ------------
	//
	// go func() {
	// 	ctx := context.Background()
	// 	if err := s1.Execute(ctx); err != nil {
	// 		logger.Error("failed executing syncing server, restarting task in 10 seconds...", slog.Any("error", err))
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }()
	// go func() {
	// 	ctx := context.Background()
	// 	if err := s2.Execute(ctx); err != nil {
	// 		logger.Error("failed executing syncing client, restarting task in 10 seconds...", slog.Any("error", err))
	// 	}
	// 	time.Sleep(10 * time.Second)
	// }()

	// tm5 := taskmnghandler.NewSyncBlockDataDTOTaskHandler(
	// 	cfg,
	// 	logger,
	// 	syncBlockDataDTOService)

	go func() {

		// ctx := context.Background()
		// for {
		// 	taskErr := tm5.Execute(ctx)
		// 	if taskErr != nil {
		// 		logger.Error("failed executing syncing task, restarting task in 10 seconds...", slog.Any("error", taskErr))
		// 		time.Sleep(10 * time.Second)
		// 	}
		// 	time.Sleep(10 * time.Second)
		// }
	}()
	<-done
}
