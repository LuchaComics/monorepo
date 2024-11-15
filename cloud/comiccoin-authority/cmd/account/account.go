package account

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	accountsURL      = "/v1/api/accounts"
	accountDetailURL = "/v1/api/account/${ACCOUNT_ADDRESS}"
)

func AccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "account",
		Short: "Execute commands related to accounts",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(NewAccountCmd())
	cmd.AddCommand(GetAccountCmd())
	cmd.AddCommand(ListAccountCmd())

	return cmd
}
