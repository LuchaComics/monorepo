package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"os"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

func init() {
	rootCmd.AddCommand(transferCmd())
}

func transferCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "transfer",
		Short: "Transfer coins between addresses",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up a wallet which has coins in it.
			//

			senderKeyJson, err := ioutil.ReadFile(flagKeystoreFile) // TODO: CHANGE To coinbase key
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			senderKey, err := keystore.DecryptKey(senderKeyJson, flagPassword) // TODO: CHANGE To coinbase key
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			//
			// STEP 2
			// Load up our blockchain.
			//

			cfg := &config.Config{
				BlockchainDifficulty: 1,
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			kvs := kvs.NewKeyValueStorer(cfg)

			bc := blockchain.NewBlockchain(cfg, kvs, senderKey)
			defer bc.Close()

			//
			// STEP 3
			//

			recipientAddress := common.HexToAddress(flagRecipientAddress)

			transferAmount := new(big.Int).SetUint64(1 * 1e18) // 1 coin
			tx := blockchain.NewTransaction(senderKey.Address, recipientAddress, transferAmount, []byte("Transfer amount"), 0)
			err = tx.Sign(senderKey.PrivateKey)
			if err != nil {
				log.Fatalf("Failed to sign transaction: %v", err)
			}

			isOK := tx.Verify()
			if isOK == false {
				log.Fatalf("Failed to sign transaction: %v", err)
			}

			err = bc.AddBlock([]*blockchain.Transaction{tx})
			if err != nil {
				log.Fatalf("Failed to add block: %v", err)
			}
			fmt.Printf("Transferred %s coin(s):\n", transferAmount.String())

			senderBalance, err := bc.GetBalance(senderKey.Address)
			if err != nil {
				log.Fatalf("Failed to get balance: %v", err)
			}
			recipientBalance, err := bc.GetBalance(recipientAddress)
			if err != nil {
				log.Fatalf("Failed to get balance: %v", err)
			}
			fmt.Printf("senders balance: %s\n", senderBalance.String())
			fmt.Printf("recipient balance: %s\n", recipientBalance.String())
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagKeystoreFile, "sender-keystore", "", "Absolute path to the coin sender's wallet")
	cmd.MarkFlagRequired("keystore")
	cmd.Flags().StringVar(&flagPassword, "sender-password", "", "The password to decrypt the coin sender's wallet")
	cmd.MarkFlagRequired("password")
	cmd.Flags().StringVar(&flagCoinbaseAddress, "coinbase-address", "", "The address of the coinbase")
	cmd.MarkFlagRequired("coinbase-address")
	cmd.Flags().StringVar(&flagRecipientAddress, "recipient-address", "", "The address of the coin(s) recipient")
	cmd.MarkFlagRequired("recipient-address")

	return cmd
}
