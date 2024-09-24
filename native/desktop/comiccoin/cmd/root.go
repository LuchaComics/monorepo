package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	flagKeystoreFile     string // Location of the wallet keystore
	flagDataDir          string // Location of the database directory
	flagPassword         string
	flagCoinbaseAddress  string
	flagRecipientAddress string
	flagAmount           uint64
	flagKeypairName      string
	flagAccountName      string
)

// Initialize function will be called when every command gets called.
func init() {
	// Attach `API` sub-commands to our main root.
	rootCmd.AddCommand(httpJsonApiCmd)
	httpJsonApiCmd.AddCommand(httpJsonApiAccountCmd())
	httpJsonApiCmd.AddCommand(httpJsonApiBlockchainCmd())
}

var rootCmd = &cobra.Command{
	Use:   "comiccoin",
	Short: "Comic Coin CLI",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing.
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var httpJsonApiCmd = &cobra.Command{
	Use:   "api",
	Short: "Execute commands for local running ComicCoin node instance via HTTP JSON API",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func httpJsonApiAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "account",
		Short: "Execute commands related user accounts",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands for `account`
	cmd.AddCommand(httpJsonApiNewAccountCmd())
	cmd.AddCommand(httpJsonApiGetAccountCmd())
	cmd.AddCommand(httpJsonApiListAccountsCmd())

	return cmd
}

func httpJsonApiBlockchainCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Execute commands related to the blockchain blockchain",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands for `blockchain`
	cmd.AddCommand(httpJsonApiBlockchainGetBalanceCmd())

	return cmd
}
