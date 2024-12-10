package cmd

import (
	"log"
	"log/slog"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/blacklist"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/jwt"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/password"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/storage/disk/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/interface/http"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/interface/http/middleware"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/usecase"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Commands used to run the ComicCoinc NFTStore service",
		Run: func(cmd *cobra.Command, args []string) {
			doDaemonCmd()
		},
	}

	return cmd
}

func doDaemonCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	// Load up the environment variables.
	dataDir := config.GetEnvString("COMICCOIN_NFTSTORAGE_APP_DATA_DIRECTORY", true)
	listenHTTPAddress := config.GetEnvString("COMICCOIN_NFTSTORAGE_ADDRESS", true)
	appSecretKey := config.GetSecureStringEnv("COMICCOIN_NFTSTORAGE_APP_SECRET_KEY", true)
	hmacSecretKey := config.GetSecureBytesEnv("COMICCOIN_NFTSTORAGE_APP_HMAC_SECRET", true)
	ipfsIP := config.GetEnvString("COMICCOIN_NFTSTORAGE_IPFS_IP", true)
	ipfsPort := config.GetEnvString("COMICCOIN_NFTSTORAGE_IPFS_PORT", true)
	ipfsPublicGatewayAddress := config.GetEnvString("COMICCOIN_NFTSTORAGE_IPFS_PUBLIC_GATEWAY", true)

	// The following block of code will be used to resolve the dns of our
	// other docker container to get the `ipfs-node` ip address.
	ips, err := net.LookupIP(ipfsIP)
	if err != nil {
		log.Fatalf("failed to lookup dns record: %v", err)
	}
	for _, ip := range ips {
		ipfsIP = ip.String()
		break
	}

	logger := logger.NewProvider()
	logger.Info("Starting daemon...")

	// Load up our operating system interaction handlers, more specifically
	// signals. The OS sends our application various signals based on the
	// OS's state, we want to listen into the termination signals.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

	config := &config.Config{
		App: config.AppConfig{
			DirPath:     dataDir,
			HTTPAddress: listenHTTPAddress,
			HMACSecret:  hmacSecretKey,
			AppSecret:   appSecretKey,
		},
		DB: config.DBConfig{
			DataDir: dataDir,
		},
	}

	passp := password.NewProvider()
	jwtp := jwt.NewProvider(config)
	blackp := blacklist.NewProvider()

	// --- Disk --- //

	pinObjsByCIDDB := disk.NewDiskStorage(config.DB.DataDir, "pin_objects_by_cid", logger)
	pinObjsByRequestIDDB := disk.NewDiskStorage(config.DB.DataDir, "pin_objects_by_request_id", logger)

	// --- Repository --- //
	ipfsRepoConfig := repo.NewIPFSRepoConfigurationProvider(
		ipfsIP,
		ipfsPort,
		ipfsPublicGatewayAddress,
	)
	ipfsRepo := repo.NewIPFSRepo(ipfsRepoConfig, logger)
	pinObjRepo := repo.NewPinObjectRepo(logger, pinObjsByCIDDB, pinObjsByRequestIDDB)

	// --- UseCase --- //

	ipfsGetNodeIDUseCase := usecase.NewIPFSGetNodeIDUseCase(logger, ipfsRepo)
	ipfsPinAddUsecase := usecase.NewIPFSPinAddUseCase(logger, ipfsRepo)
	ipfsGetUseCase := usecase.NewIPFSGetUseCase(logger, ipfsRepo)
	upsertPinObjectUseCase := usecase.NewUpsertPinObjectUseCase(logger, pinObjRepo)
	pinObjectGetByCIDUseCase := usecase.NewPinObjectGetByCIDUseCase(logger, pinObjRepo)

	// --- Service --- //

	ipfsPinAddService := service.NewIPFSPinAddService(
		config,
		logger,
		jwtp,
		passp,
		ipfsGetNodeIDUseCase,
		ipfsPinAddUsecase,
		upsertPinObjectUseCase,
	)
	pinObjectGetByCIDService := service.NewPinObjectGetByCIDService(
		logger,
		pinObjectGetByCIDUseCase,
		ipfsGetUseCase,
	)

	//
	// Interface.
	//

	// --- HTTP --- //
	getVersionHTTPHandler := httphandler.NewGetVersionHTTPHandler(
		logger)
	getHealthCheckHTTPHandler := httphandler.NewGetHealthCheckHTTPHandler(
		logger)
	ipfsGatewayHTTPHandler := httphandler.NewIPFSGatewayHTTPHandler(
		logger,
		pinObjectGetByCIDService)
	ipfsPinAddHTTPHandler := httphandler.NewIPFSPinAddHTTPHandler(
		logger,
		ipfsPinAddService)
	httpMiddleware := httpmiddle.NewMiddleware(
		logger,
		blackp)
	httpServ := http.NewHTTPServer(
		config,
		logger,
		httpMiddleware,
		getVersionHTTPHandler,
		getHealthCheckHTTPHandler,
		ipfsGatewayHTTPHandler,
		ipfsPinAddHTTPHandler,
	)

	// Run in background the peer to peer node which will synchronize our
	// blockchain with the network.
	// go peerNode.Run()
	go httpServ.Run()
	defer httpServ.Shutdown()

	logger.Info("Node running.",
		slog.Any("dataDir", dataDir),
		slog.Any("listenHTTPAddress", listenHTTPAddress),
		// slog.Any("hmacSecretKey", hmacSecretKey), // For debugging purposes only.
		// slog.Any("appSecretKey", appSecretKey),   // For debugging purposes only.
		slog.Any("ipfsIP", ipfsIP),
		slog.Any("ipfsPort", ipfsPort),
		slog.Any("ipfsPublicGatewayAddress", ipfsPublicGatewayAddress),
	)

	<-done
}
