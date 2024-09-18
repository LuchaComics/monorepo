package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.AddCommand(balancesListCmd())
}

var balanceCmd = &cobra.Command{
	Use:   "balances",
	Short: "Interacts with balances (list...).",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return balancesListCmd
}
