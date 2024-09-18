package cmd

import (
	"fmt"

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
		fmt.Println("running...")
		fmt.Println("connecting to peers:", flagBootstrapPeers)
	},
}
