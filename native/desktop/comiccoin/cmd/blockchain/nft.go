package blockchain

import "github.com/spf13/cobra"

// HTTP endpoints
const (
	nftsURL      = "/v1/api/nfts"
	nftDetailURL = "/v1/api/nft/${ACCOUNT_ADDRESS}"
)

func NFTCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "nft",
		Short: "Execute commands related to nfts (i.e. creating, transfering, etc.)",
		Run: func(cmd *cobra.Command, args []string) {
			// Do nothing...
		},
	}

	// // Attach our sub-commands for `nft`
	cmd.AddCommand(httpJsonApiMintNFTCmd())
	// cmd.AddCommand(httpJsonApiGetNFTCmd())

	return cmd
}
