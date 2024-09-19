package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport/http"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport/peer"
)

var (
	flagBootstrapPeers string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// runCmd.MarkFlagRequired("datadir")
	// runCmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of IPFS peerIDs used to synchronize our blockchain with")
	// runCmd.MarkFlagRequired("bootstrap-peers")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launches the Comic Coin node and its HTTP API.",
	Run: func(cmd *cobra.Command, args []string) {
		//
		// STEP 1
		// Load up our application and system dependencies.
		//

		// Load up our operating system interaction handlers, more specifically
		// signals. The OS sends our application various signals based on the
		// OS's state, we want to listen into the termination signals.
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

		// Load up our configuration.
		cfg := &config.Config{
			BlockchainDifficulty: 1,
			AppPort:              9000,
			DB: config.DBConfig{
				DataDir: flagDataDir,
			},
		}

		// Load up our database.
		kvs := kvs.NewKeyValueStorer(cfg)

		// Load up our blockchain.
		bc := blockchain.NewBlockchain(cfg, kvs)
		defer bc.Close()

		// Load up our peer node.
		pserv := peer.NewInputPort(cfg, kvs, bc)

		// Load up our local http server to host our local API.
		httpserv := http.NewInputPort(cfg, kvs, bc)

		//
		// STEP 2
		// Execute our application.
		//

		// Run in background the peer to peer node which will synchronize our
		// blockchain with the network.
		go pserv.Run()
		defer pserv.Shutdown()

		// Run in background the HTTP server.
		go httpserv.Run()
		defer httpserv.Shutdown()

		//
		// STEP 3
		// Run the main loop blocking code while other input ports run in
		// background.
		//

		<-done
	},
}
