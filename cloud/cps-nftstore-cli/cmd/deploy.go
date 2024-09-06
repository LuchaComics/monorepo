package cmd

import (
	"fmt"
	"log"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-cli/adapter/blockchain/eth"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(deployCmd)
}

var deployCmd = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy the `CollectibleProtectionServicesToken` smart contract to the blockchain for the first time.",
	Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("deploying...")
		fmt.Println("for private key:", cliConfig.EthServer.OwnerPrivateKey)
		adapter := eth.NewAdapter(cliConfig)
		deployedContractAddress, err := adapter.DeploySmartContractFromPrivateKey(cliConfig.EthServer.OwnerPrivateKey)
		if err != nil {
			log.Fatalf("failed deploying smart contract: %v", err)
		}
		fmt.Println("Smart Contract Address:", deployedContractAddress)
	},
}
