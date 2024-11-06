package cmd

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	pkg_config "github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	pkg_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/common/security/jwt"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/common/security/password"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http/middleware"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/usecase"
)

// Command line argument flags
var (
	flagListenHTTPAddress string
	flagAppSecretKey      string
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Commands used to run the ComicCoinc NFTStore service",
		Run: func(cmd *cobra.Command, args []string) {
			doDaemonCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flagListenHTTPAddress, "listen-http-address", "127.0.0.1:8080", "The IP and port to run our IPFS HTTP gateway on")

	return cmd
}

func doDaemonCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	appSecretKey := config.GetEnvString("COMICCOIN_NFTSTORE_APP_SECRET_KEY", true)
	hmacSecretKey := config.GetEnvBytes("COMICCOIN_NFTSTORE_HMAC_SECRET_KEY", true)

	logger := logger.NewLogger()
	logger.Info("Starting daemon...",
		slog.Any("flatHMACSecret", hmacSecretKey),
		slog.Any("appSecretKey", appSecretKey))

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

	comicCoinConfig := &pkg_config.Config{
		App: pkg_config.AppConfig{
			HTTPAddress: flagListenHTTPAddress,
		},
		IPFS: pkg_config.IPFSConfig{
			RemoteIP:            constants.ComicCoinIPFSRemoteIP,
			RemotePort:          constants.ComicCoinIPFSRemotePort,
			PublicGatewayDomain: constants.ComicCoinIPFSPublicGatewayDomain,
		},
	}
	nftStoreConfig := &config.Config{
		Blockchain: config.BlockchainConfig{
			ChainID:                        constants.ComicCoinChainID,
			TransPerBlock:                  constants.ComicCoinTransPerBlock,
			Difficulty:                     constants.ComicCoinDifficulty,
			ConsensusPollingDelayInMinutes: constants.ComicCoinConsensusPollingDelayInMinutes,
			ConsensusProtocol:              constants.ComicCoinConsensusProtocol,
		},
		App: config.AppConfig{
			DirPath:     flagDataDir,
			HTTPAddress: flagListenHTTPAddress,
			HMACSecret:  hmacSecretKey,
			AppSecret:   appSecretKey,
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

	passp := password.NewProvider()
	jwtp := jwt.NewProvider(nftStoreConfig)

	// --- Disk --- //

	pinObjsByCIDDB := disk.NewDiskStorage(nftStoreConfig.DB.DataDir, "pin_objects_by_cid", logger)
	pinObjsByRequestIDDB := disk.NewDiskStorage(nftStoreConfig.DB.DataDir, "pin_objects_by_request_id", logger)

	// --- Repository --- //

	ipfsRepo := pkg_repo.NewIPFSRepo(comicCoinConfig, logger)
	pinObjRepo := repo.NewPinObjectRepo(logger, pinObjsByCIDDB, pinObjsByRequestIDDB)

	// --- UseCase --- //

	ipfsGetNodeIDUseCase := usecase.NewIPFSGetNodeIDUseCase(logger, ipfsRepo)
	ipfsPinAddUsecase := usecase.NewIPFSPinAddUseCase(logger, ipfsRepo)
	upsertPinObjectUseCase := usecase.NewUpsertPinObjectUseCase(logger, pinObjRepo)

	// --- Service --- //

	ipfsPinAddService := service.NewIPFSPinAddService(
		nftStoreConfig,
		logger,
		jwtp,
		passp,
		ipfsGetNodeIDUseCase,
		ipfsPinAddUsecase,
		upsertPinObjectUseCase,
	)

	//
	// Interface.
	//

	// --- HTTP --- //

	ipfsGatewayHTTPHandler := httphandler.NewIPFSGatewayHTTPHandler(
		comicCoinConfig,
		logger)
	ipfsPinAddHTTPHandler := httphandler.NewIPFSPinAddHTTPHandler(
		logger,
		ipfsPinAddService)
	httpMiddleware := httpmiddle.NewMiddleware(
		comicCoinConfig,
		logger)
	httpServ := http.NewHTTPServer(
		comicCoinConfig,
		logger,
		httpMiddleware,
		ipfsGatewayHTTPHandler,
		ipfsPinAddHTTPHandler,
	)

	// Run in background the peer to peer node which will synchronize our
	// blockchain with the network.
	// go peerNode.Run()
	go httpServ.Run()
	defer httpServ.Shutdown()

	logger.Info("Node running.")

	<-done
}