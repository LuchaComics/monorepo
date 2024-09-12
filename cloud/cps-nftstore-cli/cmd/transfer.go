package cmd

import (
	"log"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/adapter/blockchain/eth"
	"github.com/spf13/cobra"
)

var (
	privateKey string
)

func init() {
	transferCmd.Flags().StringVarP(&privateKey, "privatekey", "i", "", "")
	transferCmd.MarkFlagRequired("privatekey")
	transferCmd.Flags().StringVarP(&toAddress, "to", "o", "", "")
	transferCmd.MarkFlagRequired("to")
	rootCmd.AddCommand(transferCmd)
}

var transferCmd = &cobra.Command{
	Use:   "transfer",
	Short: "Transfer eth from account with eth to wallet",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {

		adapter := eth.NewAdapter(cliConfig)
		err := adapter.Transfer(privateKey, toAddress)
		if err != nil {
			log.Fatal(err)
		}

	},
}
