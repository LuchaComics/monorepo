package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/config"
)

var (
	// Global state variables for command line interfrace
	cliConfig *config.Conf

	// Command line interface input arguments.
	smartContractAddress string
	toAddress            string
	tokenID              uint64
)

// Initialize function will be called when every command gets called.
func init() {
	cliConfig = config.New()

}

var rootCmd = &cobra.Command{
	Use:   "cps-nftstore-cli",
	Short: "",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
