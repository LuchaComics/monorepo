package peer

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
	flagPeerID           string

	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string
	flagProtocolID       string

	flagListenHTTPPort       int
	flagListenHTTPIP         string
	flagListenPeerToPeerPort int

	flagIdentityKeyID string
)

// HTTP endpoints
const (
	accountsURL      = "/v1/api/peers"
	accountDetailURL = "/v1/api/peer/${ACCOUNT_ID}"
)

func PeerCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "peer",
		Short: "Execute commands related to ComicCoin blockchain peer-to-peer networking functionality",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands for `account`
	// Version Command
	cmd.AddCommand(versionCmd)
	cmd.AddCommand(identityCmd())
	// cmd.AddCommand(httpJsonApiNewPeerCmd())
	// cmd.AddCommand(httpJsonApiGetPeerCmd())

	return cmd
}
