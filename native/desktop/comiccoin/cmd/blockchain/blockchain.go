package blockchain

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

	flagListenHTTPAddress string
	flagListenRPCAddress  string

	flagIdentityKeyID string
)

// // HTTP endpoints
// const (
// 	accountsURL      = "/v1/api/accounts"
// 	accountDetailURL = "/v1/api/account/${ACCOUNT_ID}"
// )

func BlockchainCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Commands related to blockchain operations (Create Account, Submit Payment, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands
	cmd.AddCommand(accountCmd())
	cmd.AddCommand(txCmd())
	cmd.AddCommand(InitCmd())

	return cmd
}
