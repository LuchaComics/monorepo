package coins

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	coinsTransferURL = "/v1/api/coins-transfer"
)

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
