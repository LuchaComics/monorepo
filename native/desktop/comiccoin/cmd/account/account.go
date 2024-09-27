package account

import "github.com/spf13/cobra"

// Command line argument flags
var (
	flagKeystoreFile     string // Location of the wallet keystore
	flagDataDir          string // Location of the database directory
	flagPassword         string
	flagCoinbaseAddress  string
	flagRecipientAddress string
	flagAmount           uint64
	flagKeypairName      string
	flagAccountID        string

	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string
	flagProtocolID       string

	flagListenHTTPPort       int
	flagListenHTTPIP         string
	flagListenPeerToPeerPort int
)

// HTTP endpoints
const (
	accountsURL      = "/v1/api/accounts"
	accountDetailURL = "/v1/api/account/${ACCOUNT_ID}"
)

func AccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "account",
		Short: "Execute commands related to accounts",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands for `account`
	// Version Command
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(runCmd())
	cmd.AddCommand(httpJsonApiNewAccountCmd())
	cmd.AddCommand(httpJsonApiGetAccountCmd())

	return cmd
}
