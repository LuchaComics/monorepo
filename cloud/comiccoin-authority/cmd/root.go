package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/cmd/account"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/cmd/coins"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/cmd/daemon"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/cmd/genesis"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/cmd/version"
)

// Initialize function will be called when every command gets called.
func init() {

}

var rootCmd = &cobra.Command{
	Use:   "comiccoin-authority",
	Short: "ComicCoin Authority",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	// Attach sub-commands to our main root.
	rootCmd.AddCommand(account.AccountCmd())
	rootCmd.AddCommand(coins.CoinsCmd())
	rootCmd.AddCommand(genesis.GenesisCmd())
	rootCmd.AddCommand(daemon.DaemonCmd())
	rootCmd.AddCommand(version.VersionCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
