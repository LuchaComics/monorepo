package tokens

import (
	"github.com/spf13/cobra"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/preferences"
)

// HTTP endpoints
const (
	tokensTransferURL = "/v1/api/tokens-transfer"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

func TokensCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tokens",
		Short: "Execute commands related to tokens",
		Run: func(cmd *cobra.Command, args []string) {
			// Developers Note:
			// Before executing this command, check to ensure the user has
			// configured our app before proceeding.
			preferences.RunFatalIfHasAnyMissingFields()
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(TransferTokensCmd())
	cmd.AddCommand(GetTokenCmd())
	cmd.AddCommand(DownloadTokenCmd())
	cmd.AddCommand(TransferTokensCmd())
	cmd.AddCommand(BurnTokensCmd())
	cmd.AddCommand(ListTokensCmd())

	return cmd
}
