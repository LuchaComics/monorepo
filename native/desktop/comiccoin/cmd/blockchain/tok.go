package blockchain

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	mintTokensURL     = "/v1/api/tokens/mint"
	transferTokensURL = "/v1/api/tokens/transfer"
	burnTokensURL     = "/v1/api/tokens/burn"
	tokenDetailURL    = "/v1/api/token/${TOKEN_ID}"
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
	cmd.AddCommand(httpJsonApiGetTokenCmd())
	cmd.AddCommand(httpJsonApiBurnTokenCmd())

	return cmd
}
