package blockchain

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	mintTokensURL     = "/v1/api/tokens/mint"
	transferTokensURL = "/v1/api/tokens/transfer"
	tokenDetailURL    = "/v1/api/token/${ACCOUNT_ADDRESS}"
)

func TokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "token",
		Short: "Execute commands related to tokens (i.e. creating, transfering, etc.)",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `token`
	cmd.AddCommand(httpJsonApiMintTokenCmd())
	cmd.AddCommand(httpJsonApiTransferTokenCmd())

	return cmd
}
