package blockchain

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	transactionsURL = "/v1/api/coins/transfer"
)

func CoinCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "coin",
		Short: "Execute commands related to your coins",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `tx`
	cmd.AddCommand(transferCoinCmd())

	return cmd
}
