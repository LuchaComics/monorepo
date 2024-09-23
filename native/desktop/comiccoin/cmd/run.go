package cmd

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	ledger_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/ledger/controller"
	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
)

var ()

func init() {
	rootCmd.AddCommand(runCmd())
}

func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run",
		Short: "Get balance of address",
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
				BlockchainDifficulty: 1,
				Peer: config.PeerConfig{
					ListenPort:       flagListenPort,
					KeyName:          flagKeypairName,
					RendezvousString: flagRendezvousString,
					BootstrapPeers:   bootstrapPeers,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			keypairDS := keypair_ds.NewDatastore(cfg, logger, kvs)
			lastHashDS := lasthash_ds.NewDatastore(cfg, logger, kvs)
			blockDS := block_ds.NewDatastore(cfg, logger, kvs)
			ledgerController := ledger_c.NewController(cfg, logger, lastHashDS, blockDS)
			node := p2p.NewInputPort(cfg, logger, keypairDS, ledgerController)

			//
			// STEP 2
			// Execute our application.
			//

			// Run in background the peer to peer node which will synchronize our
			// ledger with the network.
			go node.Run()
			defer node.Shutdown()

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
	runCmd.Flags().IntVar(&flagListenPort, "listen-port", 9000, "The port to listen to for other peers")
	runCmd.MarkFlagRequired("listen-port")
	runCmd.Flags().StringVar(&flagKeypairName, "keypair-name", "", "The name of keypairs to apply to this server.")
	runCmd.MarkFlagRequired("keypair-name")
	runCmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of peers used to synchronize our ledger with")
	// runCmd.MarkFlagRequired("bootstrap-peers")

	runCmd.Flags().StringVar(&flagRendezvousString, "rendezvous", "meet me here",
		"Unique string to identify group of nodes. Share this with your friends to let them connect with you")
	runCmd.Flags().StringVar(&flagBootstrapPeers, "peer", "", "Adds a peer multiaddress to the bootstrap list")
	runCmd.Flags().StringVar(&flagListenAddresses, "listen", "", "Adds a multiaddress to the listen list")

	return runCmd
}
