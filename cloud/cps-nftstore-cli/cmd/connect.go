package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/adapter/blockchain/eth"
)

func init() {
	rootCmd.AddCommand(connectCmd)
}

var connectCmd = &cobra.Command{
	Use:   "connect",
	Short: "Test connecting to eth-node",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		eth.NewAdapter(cliConfig)
		fmt.Println("we have a connection with node at url:", cliConfig.EthServer.NodeURL)
	},
}
