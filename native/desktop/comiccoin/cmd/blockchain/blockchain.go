package blockchain

import "github.com/spf13/cobra"

// Command line argument flags
var (
	flagKeystoreFile     string // Location of the wallet keystore
	flagDataDir          string // Location of the database directory
	flagLabel            string
	flagPassword         string
	flagPasswordRepeated string
	flagCoinbaseAddress  string
	flagRecipientAddress string
	flagAmount           uint64
	flagKeypairName      string
	flagAccountAddress   string
	flagData             string

	flagRendezvousString string
	flagBootstrapPeers   string
	flagListenAddresses  string

	flagListenHTTPPort       int
	flagListenHTTPIP         string
	flagListenPeerToPeerPort int

	flagListenHTTPAddress string

	flagIdentityKeyID string
)

func BlockchainCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "blockchain",
		Short: "Commands related to blockchain operations (Create Account, Submit Payment, etc)",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// Attach our sub-commands
	cmd.AddCommand(AccountCmd())
	cmd.AddCommand(CoinCmd())
	cmd.AddCommand(InitCmd())
	cmd.AddCommand(TokenCmd())

	return cmd
}
