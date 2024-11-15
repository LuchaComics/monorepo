package coins

import (
	"github.com/spf13/cobra"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/preferences"
)

// HTTP endpoints
const (
	coinsTransferURL = "/v1/api/coins-transfer"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

func CoinsCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "coins",
		Short: "Execute commands related to coins",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(TransferCoinsCmd())

	return cmd
}
