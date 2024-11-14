package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/cmd/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/cmd/daemon"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/cmd/initialization"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/cmd/version"
)

var (
	flagKeystoreFile     string // Location of the wallet keystore
	flagDataDir          string // Location of the database directory
	flagPassword         string
	flagCoinbaseAddress  string
	flagRecipientAddress string
	flagAmount           uint64
	flagKeypairName      string
	flagAccountName      string
)

// Initialize function will be called when every command gets called.
func init() {

}

var rootCmd = &cobra.Command{
	Use:   "comiccoin",
	Short: "ComicCoin CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	// Attach sub-commands to our main root.
	rootCmd.AddCommand(initialization.InitCmd())
	rootCmd.AddCommand(blockchain.BlockchainCmd())
	rootCmd.AddCommand(daemon.DaemonCmd())
	rootCmd.AddCommand(version.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
