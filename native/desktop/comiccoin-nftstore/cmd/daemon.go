package cmd

import (
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/config/constants"
	http "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http"
	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http/handler"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http/middleware"
)

// Command line argument flags
var (
	flagListenHTTPAddress string
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
			DirPath:     flagDataDir,
			HTTPAddress: flagListenHTTPAddress,
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

	//
	// Interface.
	//

	// --- HTTP --- //

	ipfsGatewayHTTPHandler := httphandler.NewIPFSGatewayHTTPHandler(
		cfg,
		logger)
	httpMiddleware := httpmiddle.NewMiddleware(
		cfg,
		logger)
	httpServ := http.NewHTTPServer(
		cfg,
		logger,
		httpMiddleware,
		ipfsGatewayHTTPHandler,
	)

	// Run in background the peer to peer node which will synchronize our
	// blockchain with the network.
	// go peerNode.Run()
	go httpServ.Run()
	defer httpServ.Shutdown()

	logger.Info("Node running.")

	<-done
}
