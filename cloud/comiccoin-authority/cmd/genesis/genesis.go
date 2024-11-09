package genesis

import "github.com/spf13/cobra"

func GenesisCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "genesis",
		Short: "Execute commands related to genesis of the blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(NewGenesistCmd())

	return cmd
}
