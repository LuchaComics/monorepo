package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

func init() {
	rootCmd.AddCommand(genesisCmd())
}

func genesisCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "genesis",
		Short: "Initialize ComicCoin blockchain by creating the genesis block",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up a wallet which has coins in it.
			//

			coinbaseKeyJson, err := ioutil.ReadFile(flagKeystoreFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			coinbaseKey, err := keystore.DecryptKey(coinbaseKeyJson, flagPassword)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			//
			// STEP 2
			// Load up our blockchain.
			//

			// Load up the configuration.
			cfg := &config.Config{
				BlockchainDifficulty: 1,
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}

			// Load up our database.
			kvs := kvs.NewKeyValueStorer(cfg)

			// Finally load up the blockchain we will be initializing.
			bc := blockchain.NewBlockchainWithCoinbaseKey(cfg, kvs, coinbaseKey)
			defer bc.Close()

			//
			// STEP 3
			//

			coinbaseBalance, err := bc.GetBalance(coinbaseKey.Address)
			if err != nil {
				log.Fatalf("Failed to get balance: %v", err)
			}
			fmt.Printf("coinbases balance: %s\n", coinbaseBalance.String())
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagKeystoreFile, "coinbase-keystore", "", "Absolute path to the coinbase's wallet")
	cmd.MarkFlagRequired("coinbase-keystore")
	cmd.Flags().StringVar(&flagPassword, "coinbase-password", "", "The password to decrypt the cointbase's wallet")
	cmd.MarkFlagRequired("coinbase-password")

	return cmd
}
