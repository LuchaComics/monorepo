package cmd

import (
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.AddCommand(balancesListCmd())
	balanceCmd.AddCommand(balanceGetCmd())
}

var balanceCmd = &cobra.Command{
	Use:   "balances",
	Short: "Interacts with balances (list...).",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func balancesListCmd() *cobra.Command {
	var balancesListCmd = &cobra.Command{
		Use:   "list",
		Short: "Lists all balances.",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}

	return balancesListCmd
}

func balanceGetCmd() *cobra.Command {
	var balanceGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get balance of account.",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our blockchain.
			//

			cfg := &config.Config{
				BlockchainDifficulty: 1,
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			kvs := kvs.NewKeyValueStorer(cfg)

			bc := blockchain.NewBlockchain(cfg, kvs)
			defer bc.Close()

			//
			// STEP 2
			// Lookup balance.
			//

			address := common.HexToAddress(flagRecipientAddress)
			balance, err := bc.GetBalance(address)
			if err != nil {
				log.Fatalf("Failed to get balance: %v", err)
			}

			fmt.Printf("Balance: %s\n", balance.String())

		},
	}
	balanceGetCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	balanceGetCmd.Flags().StringVar(&flagRecipientAddress, "address", "", "The address of the coin(s) receipient")
	balanceGetCmd.MarkFlagRequired("address")

	return balanceGetCmd
}
