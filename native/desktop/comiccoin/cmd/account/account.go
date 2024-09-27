package account

import "github.com/spf13/cobra"

var (
	flagKeystoreFile     string // Location of the wallet keystore
	flagDataDir          string // Location of the database directory
	flagPassword         string
	flagCoinbaseAddress  string
	flagRecipientAddress string
	flagAmount           uint64
	flagKeypairName      string
	flagAccountName      string

	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string
	flagProtocolID       string

	flagListenHTTPPort       int
	flagListenHTTPIP         string
	flagListenPeerToPeerPort int
)

func AccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "account",
		Short: "Execute commands related user accounts",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `account`
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(runCmd)

	return cmd
}
