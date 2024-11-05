package cmd

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/kmutexutil"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/net/p2p"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	taskmnghandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin/interface/task/handler"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	ma "github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config/constants"
)

// Command line argument flags
var ()

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Commands used to run the ComicCoinc NFTStore service",
		Run: func(cmd *cobra.Command, args []string) {
			doDaemonCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")

	return cmd
}

func doDaemonCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	logger := logger.NewLogger()

	// Load up our operating system interaction handlers, more specifically
	// signals. The OS sends our application various signals based on the
	// OS's state, we want to listen into the termination signals.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

	// DEVELOPERS NOTE:
	// Every ComicCoin node must be connected to a peer whom coordinates
	// connecting all the other nodes in the network, therefore we get the
	// following node(s) that act in this role.
	bootstrapPeers, err := config.StringToAddres(constants.ComicCoinBootstrapPeers)
	if err != nil {
		logger.Error("Startup aborted: failed converting string to multi-addresses",
			slog.Any("error", err))
		log.Fatalf("Failed converting string to multi-addresses: %v\n", err)
	}

	cfg := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID:                        constants.ComicCoinChainID,
			TransPerBlock:                  constants.ComicCoinTransPerBlock,
			Difficulty:                     constants.ComicCoinDifficulty,
			ConsensusPollingDelayInMinutes: constants.ComicCoinConsensusPollingDelayInMinutes,
			ConsensusProtocol:              constants.ComicCoinConsensusProtocol,
		},
		App: config.AppConfig{
			DirPath: flagDataDir,
		},
		DB: config.DBConfig{
			DataDir: flagDataDir,
		},
		Peer: config.PeerConfig{
			ListenPort:     constants.ComicCoinPeerListenPort,
			KeyName:        constants.ComicCoinIdentityKeyID,
			BootstrapPeers: bootstrapPeers,
		},
		IPFS: config.IPFSConfig{
			RemoteIP:            constants.ComicCoinIPFSRemoteIP,
			RemotePort:          constants.ComicCoinIPFSRemotePort,
			PublicGatewayDomain: constants.ComicCoinIPFSPublicGatewayDomain,
		},
	}

	kmutex := kmutexutil.NewKMutexProvider()
	ikDB := disk.NewDiskStorage(cfg.DB.DataDir, "identity_key", logger)
	blockDataDB := disk.NewDiskStorage(cfg.DB.DataDir, "block_data", logger)
	tokDB := disk.NewDiskStorage(cfg.DB.DataDir, "token", logger)
	nftokDB := disk.NewDiskStorage(cfg.DB.DataDir, "non_fungible_token", logger)

	logger.Debug("Startup loading peer-to-peer client...")
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
	// a.libP2PNetwork = libP2PNetwork
	_ = libP2PNetwork

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

	genesisBlockDataRepo := repo.NewGenesisBlockDataRepo(
		cfg,
		logger,
		blockDataDB)
	tokRepo := repo.NewTokenRepo(
		cfg,
		logger,
		tokDB)
	issuedTokenDTO := repo.NewIssuedTokenDTORepo(
		cfg,
		logger,
		libP2PNetwork)
	nftokRepo := repo.NewNonFungibleTokenRepo(
		logger,
		nftokDB)
	ipfsRepo := repo.NewIPFSRepo(cfg, logger)

	//
	// Use-case.
	//

	logger.Debug("Startup loading usecases...")

	_ = nftokDB
	_ = tokDB
	_ = tokRepo
	_ = nftokRepo
	_ = ipfsRepo

	receiveIssuedTokenDTOUseCase := usecase.NewReceiveIssuedTokenDTOUseCase(
		cfg,
		logger,
		issuedTokenDTO)

	// Genesis Block Data
	loadGenesisBlockDataUseCase := usecase.NewLoadGenesisBlockDataUseCase(
		cfg,
		logger,
		genesisBlockDataRepo)

	// Token
	upsertTokenIfPreviousTokenNonceGTEUseCase := usecase.NewUpsertTokenIfPreviousTokenNonceGTEUseCase(
		cfg,
		logger,
		tokRepo)

	//
	// Services.
	//

	issuedTokenClientService := service.NewIssuedTokenClientService(
		cfg,
		logger,
		kmutex,
		receiveIssuedTokenDTOUseCase,
		loadGenesisBlockDataUseCase,
		upsertTokenIfPreviousTokenNonceGTEUseCase,
	)

	//
	// Interface.
	//

	tm10 := taskmnghandler.NewIssuedTokenClientServiceTaskHandler(
		cfg,
		logger,
		issuedTokenClientService)

	logger.Info("Node running.")

	go func(server *taskmnghandler.IssuedTokenClientServiceTaskHandler, loggerp *slog.Logger) {
		loggerp.Info("Running issued token dto server...")
		ctx := context.Background()
		for {
			if err := server.Execute(ctx); err != nil {
				loggerp.Error("issued token server error",
					slog.Any("error", err))
				time.Sleep(10 * time.Second)
				continue
			}
			// DEVELOPERS NOTE:
			// No need for delays, automatically start executing again.
			logger.Debug("issued token dto server executing again ...")
		}
	}(tm10, logger)

	<-done
}
