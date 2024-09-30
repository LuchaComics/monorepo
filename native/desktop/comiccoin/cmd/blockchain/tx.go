package blockchain

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	transactionsURL = "/v1/api/txs"
)

func txCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "tx",
		Short: "Execute commands related to transactions",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `tx`
	cmd.AddCommand(submitTxCmd())

	return cmd
}
