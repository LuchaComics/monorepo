package peer

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/config"
	ik_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/repo"
	ik_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/service"
	ik_use "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/usecase"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

func runCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "run",
		Short: "Starts the peer-to-peer service for ComicCoin",
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
				App: config.AppConfig{
					HTTPPort: flagListenHTTPPort,
					HTTPIP:   flagListenHTTPIP,
					DirPath:  flagDataDir,
				},
				Peer: config.PeerConfig{
					ListenPort:       flagListenPeerToPeerPort,
					KeyName:          flagKeypairName,
					RendezvousString: flagRendezvousString,
					BootstrapPeers:   bootstrapPeers,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewLogger()
			db := dbase.NewDatabase(cfg.DB.DataDir, logger)
			ikRepo := ik_repo.NewIdentityKeyRepo(cfg, logger, db)
			ikCreateUseCase := ik_use.NewCreateIdentityKeyUseCase(cfg, logger, ikRepo)
			ikGetUseCase := ik_use.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
			ikCreateService := ik_s.NewCreateIdentityKeyService(cfg, logger, ikCreateUseCase, ikGetUseCase)
			ikGetService := ik_s.NewGetIdentityKeyService(cfg, logger, ikGetUseCase)

			//TODO
			_ = ikCreateService
			_ = ikGetService

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
			libp2pnet := p2p.NewLibP2PNetwork(cfg, logger, privateKey, publicKey)

			_ = libp2pnet

			// USE CASES - NETWORK
			// - Share Signed Pending Transaction (Publisher)
			// - Receive Signed Pending Transaction (Subscriber)
			// - Share Purpose Block Data (Publisher)
			// - Receive Purpose Block Data (Subscriber)
			// - Share Block Data (Publisher)
			// - Receive Block Data (Subscriber)
			// - Ask Latest Block Hash (Req-Res)
			// - Receive Latest Block Hash (Req-Res)
			// - Ask Block Data (Req-Res)
			// - Receive Block Data (Req-Res)

			// USE CASES - HTTP
			// Send/Receive Signed Pending Transaction
			// Send/Receive Purpose Block Data
			// Send/Receive Block Data
			// Send/Receive Latest Block Hash
			// Send/Receive Block Data

			//
			// STEP 2
			// Execute our application.
			//

			//TODO:

			//
			// STEP 3
			// Run the main loop blocking code while other input ports run in
			// background.
			//

			<-done
		},
	}
	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagIdentityKeyID, "id", "", "The unique identifier  to use to lookup the identity key and assign to this peer")
	cmd.MarkFlagRequired("id")
	cmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", 26642, "The port to listen to for other peers")
	cmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8001, "The port to listen to for the HTTP JSON API server")
	cmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "The IP address to attach our HTTP JSON API server")
	cmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")

	return cmd
}
