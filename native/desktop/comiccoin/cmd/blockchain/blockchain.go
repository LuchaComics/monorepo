package blockchain

import (
	"github.com/spf13/cobra"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/preferences"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

func BlockchainCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Commands related to blockchain operations (Create Account, Submit Payment, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands
	cmd.AddCommand(BlockchainSyncCmd())

	return cmd
}
