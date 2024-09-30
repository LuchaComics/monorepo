package blockchain

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	accountsURL      = "/v1/api/accounts"
	accountDetailURL = "/v1/api/account/${ACCOUNT_ID}"
)

func accountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "account",
		Short: "Execute commands related to accounts",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(httpJsonApiNewAccountCmd())
	cmd.AddCommand(httpJsonApiGetAccountCmd())

	return cmd
}