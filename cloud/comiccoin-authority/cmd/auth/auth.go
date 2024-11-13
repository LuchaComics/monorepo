package auth

import "github.com/spf13/cobra"

func AuthCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "auth",
		Short: "Execute commands related to administrator authentication and authorization",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(GenerateAPIKeyCmd())

	return cmd
}
