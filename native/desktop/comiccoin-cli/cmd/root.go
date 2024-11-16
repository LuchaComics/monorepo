package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/cmd/account"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/cmd/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/cmd/coins"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/cmd/daemon"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/cmd/tokens"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/cmd/version"
	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/preferences"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

var rootCmd = &cobra.Command{
	Use:   "comiccoin-cli",
	Short: "ComicCoin CLI",
	Long:  `ComicCoin Command Line Interface`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	// Attach sub-commands to our main root.
	rootCmd.AddCommand(version.VersionCmd())
	rootCmd.AddCommand(account.AccountCmd())
	rootCmd.AddCommand(blockchain.BlockchainCmd())
	rootCmd.AddCommand(coins.CoinsCmd())
	rootCmd.AddCommand(daemon.DaemonCmd())
	rootCmd.AddCommand(tokens.TokensCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
