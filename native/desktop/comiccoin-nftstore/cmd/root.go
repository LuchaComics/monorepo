package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/cmd/version"
)

var (
	flagDataDir string // Location of the database directory
)

// Initialize function will be called when every command gets called.
func init() {

}

var rootCmd = &cobra.Command{
	Use:   "comiccoin-nftstore",
	Short: "ComicCoin NFT Store CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	// Attach sub-commands to our main root.
	rootCmd.AddCommand(version.VersionCmd())
	rootCmd.AddCommand(GenerateAPIKeyCmd())
	rootCmd.AddCommand(PinAddCmd())
	rootCmd.AddCommand(DaemonCmd())

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
