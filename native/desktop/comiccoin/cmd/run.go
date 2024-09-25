package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	dmqb "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/distributedmessagequeue"
	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	mqb "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/messagequeuebroker/simple"
	acc_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/controller"
	acc_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	acc_http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/httptransport"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	blockchain_http "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/httptransport"
	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport/http"
	httpmiddle "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport/http/middleware"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/uuid"
)

var ()

func init() {
	rootCmd.AddCommand(runCmd())
}

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Run a ComicCoin node instance",
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
				BlockchainDifficulty: 1,
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
			logger := logger.NewProvider()
			uuid := uuid.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			broker := mqb.NewMessageQueue(cfg, logger)
			keypairDS := keypair_ds.NewDatastore(cfg, logger, kvs)
			msgqueue := dmqb.NewDistributedMessageQueueAdapter(cfg, logger, keypairDS)
			accountDS := acc_s.NewDatastore(cfg, logger, kvs)
			accountController := acc_c.NewController(cfg, logger, accountDS)
			lastHashDS := lasthash_ds.NewDatastore(cfg, logger, kvs)
			signedTransactionDS := pt_ds.NewDatastore(cfg, logger, kvs)
			blockDS := block_ds.NewDatastore(cfg, logger, kvs)
			blockchainController := blockchain_c.NewController(cfg, logger, uuid, broker, msgqueue, accountDS, signedTransactionDS, lastHashDS, blockDS)
			accountHttp := acc_http.NewHandler(logger, accountController)
			blockchainHttp := blockchain_http.NewHandler(logger, blockchainController)
			// mempoolController := mempool_c.NewController(cfg, logger, uuid, broker, msgqueue, signedTransactionDS)
			// mempoolNode := mempool_p2p.NewNode(logger, mempoolController)
			// peerNode := p2p.NewInputPort(cfg, logger, keypairDS, mempoolNode)
			httpMiddleware := httpmiddle.NewMiddleware(cfg, logger)
			httpServ := http.NewInputPort(cfg, logger, httpMiddleware, accountHttp, blockchainHttp)

			//TODO: DELETE BELOW
			ctx := context.Background()
			priv, _, err := keypairDS.GetByName(ctx, cfg.Peer.KeyName)
			if err != nil {
				panic("test")
			}
			time.Sleep(10 * time.Second)
			msg := fmt.Sprintf("testing from %v", priv)
			go msgqueue.Publish(ctx, "mempool", []byte(msg))
			res := msgqueue.Subscribe(ctx, "mempool")
			log.Println(string(res))

			//
			// STEP 2
			// Execute our application.
			//

			// Run in background the peer to peer node which will synchronize our
			// blockchain with the network.
			// go peerNode.Run()
			go httpServ.Run()
			// defer peerNode.Shutdown()
			defer httpServ.Shutdown()

			//
			// STEP 3
			// Run the main loop blocking code while other input ports run in
			// background.
			//

			<-done

		},
	}
	runCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// runCmd.MarkFlagRequired("datadir")
	runCmd.Flags().IntVar(&flagListenPeerToPeerPort, "listen-p2p-port", 26642, "The port to listen to for other peers")
	// runCmd.MarkFlagRequired("listen-port")

	runCmd.Flags().IntVar(&flagListenHTTPPort, "listen-http-port", 8000, "The port to listen to for the HTTP JSON API server")
	runCmd.Flags().StringVar(&flagListenHTTPIP, "listen-http-ip", "127.0.0.1", "The IP address to attach our HTTP JSON API server")

	runCmd.Flags().StringVar(&flagKeypairName, "keypair-name", "", "The name of keypairs to apply to this server")
	runCmd.MarkFlagRequired("keypair-name")
	runCmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our blockchain with")
	// runCmd.MarkFlagRequired("bootstrap-peers")
	runCmd.Flags().StringVar(&flagRendezvousString, "rendezvous", "meet me here",
		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	// runCmd.Flags().StringVar(&flagBootstrapPeers, "peer", "", "Adds a peer multiaddress to the bootstrap list")
	// runCmd.Flags().StringVar(&flagListenAddresses, "listen", "", "Adds a multiaddress to the listen list")

	return runCmd
}
