package cmd

import (
	"log"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/adapter/blockchain/eth"
	"github.com/spf13/cobra"
)

func init() {
	mintCmd.Flags().StringVarP(&smartContractAddress, "smartcontract", "i", "", "")
	mintCmd.MarkFlagRequired("smartcontract")
	mintCmd.Flags().StringVarP(&toAddress, "to", "o", "", "")
	mintCmd.MarkFlagRequired("to")
	rootCmd.AddCommand(mintCmd)
}

var mintCmd = &cobra.Command{
	Use:   "mint",
	Short: "Create a NFT for a specified address",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		adapter := eth.NewAdapter(cliConfig)
		err := adapter.Mint(cliConfig.EthServer.OwnerPrivateKey, smartContractAddress, toAddress)
		if err != nil {
			log.Fatal(err)
		}

	},
}
