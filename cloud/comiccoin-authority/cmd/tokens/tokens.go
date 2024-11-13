package tokens

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	accountsURL      = "/v1/api/accounts"
	accountDetailURL = "/v1/api/account/${ACCOUNT_ADDRESS}"
)

func TokensCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tokens",
		Short: "Execute commands related to tokens",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // // Attach our sub-commands for `account`
	cmd.AddCommand(MintTokenCmd())
	cmd.AddCommand(GetTokenCmd())
	cmd.AddCommand(TransferTokenCmd())
	cmd.AddCommand(BurnTokenCmd())

	return cmd
}
