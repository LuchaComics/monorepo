package cmd

import (
	"log"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/adapter/blockchain/eth"
	"github.com/spf13/cobra"
)

func init() {
	nftCmd.Flags().StringVarP(&smartContractAddress, "smartcontract", "i", "", "")
	nftCmd.MarkFlagRequired("smartcontract")
	nftCmd.Flags().Uint64VarP(&tokenID, "tokenid", "d", 0, "")
	nftCmd.MarkFlagRequired("subject_id")
	rootCmd.AddCommand(nftCmd)
}

var nftCmd = &cobra.Command{
	Use:   "nft",
	Short: "View the content of the non-fungible token",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		adapter := eth.NewAdapter(cliConfig)
		tokenID := big.NewInt(0)
		x, err := adapter.GetTokenURI(smartContractAddress, tokenID)
		if err != nil {
			log.Fatal(err)
		}
		log.Println(x)

	},
}
