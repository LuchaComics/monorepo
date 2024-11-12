package blockchain

import "github.com/spf13/cobra"

// Command line argument flags
var ()

func BlockchainCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Commands related to blockchain operations (Create Account, Submit Payment, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands
	cmd.AddCommand(BlockchainSyncCmd())

	return cmd
}
