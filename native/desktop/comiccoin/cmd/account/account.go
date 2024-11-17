package account

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

func AccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "account",
		Short: "Execute commands related to accounts",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // // Attach our sub-commands for `account`
	cmd.AddCommand(NewAccountCmd())
	cmd.AddCommand(GetAccountCmd())
	cmd.AddCommand(ListAccountCmd())

	return cmd
}
