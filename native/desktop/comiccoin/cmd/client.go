package cmd

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientHelloCmd())
}

var (
	flagPeerAddress string
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to remote nodes",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func clientHelloCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "hello",
		Short: "Sends a hello message to peer, use this to test connection",
		Run: func(cmd *cobra.Command, args []string) {
			// // Parse the multiaddr string.
			// peerMA, err := multiaddr.NewMultiaddr(flagPeerAddress)
			// if err != nil {
			// 	panic(err)
			// }
			// peerAddrInfo, err := peer.AddrInfoFromP2pAddr(peerMA)
			// if err != nil {
			// 	panic(err)
			// }
			//
			// // Connect to the node at the given address.
			// if err := host.Connect(context.Background(), *peerAddrInfo); err != nil {
			// 	panic(err)
			// }
			// fmt.Println("Connected to", peerAddrInfo.String())

		},
	}

	// cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// // cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagPeerAddress, "peer-address", "", "The address of the remote node")
	cmd.MarkFlagRequired("peer-address")

	return cmd
}
