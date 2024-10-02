package blockchain

import "github.com/spf13/cobra"

func SyncCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "sync",
		Short: "Fetch the latest blockchain from the peer-to-peer network.",
		Run: func(cmd *cobra.Command, args []string) {
			doBlockchainSync()
		},
	}

	return cmd
}

func doBlockchainSync() {
	//TODO: IMPL.
}
