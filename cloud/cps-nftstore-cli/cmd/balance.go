package cmd

import (
	"fmt"
	"log"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/adapter/blockchain/eth"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Deploy the `CollectibleProtectionServiceSubmissions` smart contract to the blockchain for the first time.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		adapter := eth.NewAdapter(cliConfig)
		balance, err := adapter.Balance(cliConfig.EthServer.OwnerAddress)
		if err != nil {
			log.Fatalf("failed getting account authorization: %v", err)
		}
		fmt.Println(balance)
	},
}
