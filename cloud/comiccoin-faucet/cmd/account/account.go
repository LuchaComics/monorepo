package account

import "github.com/spf13/cobra"

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

	return cmd
}
