package cmd

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	flagBootstrapPeers string
)

func init() {
	rootCmd.AddCommand(runCmd)
	runCmd.Flags().StringVar(&flagBootstrapPeers, "bootstrap-peers", "", "The list of IPFS peerIDs used to synchronize our blockchain with")
	runCmd.MarkFlagRequired("bootstrap-peers")
}

var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Launches the Comic Coin node and its HTTP API.",
	Run: func(cmd *cobra.Command, args []string) {
		done := make(chan os.Signal, 1)
		signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

		// Run in background the HTTP server.
		go doExecuteHTTPServer()

		// Run in background the event scheduler server.
		go doExecuteIPFSServer()

		// Run the main loop blocking code while other input ports run in background.
		<-done

		doShutdown()
	},
}

func doExecuteHTTPServer() {

}

func doExecuteIPFSServer() {

}

func doShutdown() {
	// a.HttpServer.Shutdown()
	// a.Logger.Info("Application shutdown")
}
