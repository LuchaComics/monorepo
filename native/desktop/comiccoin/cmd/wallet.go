package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/wallet"
)

func init() {
	rootCmd.AddCommand(walletCmd)
	walletCmd.AddCommand(walletNewAccountCmd())
	walletCmd.AddCommand(walletPrintPrivKeyCmd())
}

var walletCmd = &cobra.Command{
	Use:   "wallet",
	Short: "Manages blockchain accounts and keys.",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func walletNewAccountCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new-account",
		Short: "Creates a new account with a new set of a elliptic-curve Private + Public keys.",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Creating new wallet...")
			acc, err := wallet.NewKeystoreAccount(flagDataDir, flagPassword)
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
			fmt.Printf("New wallet created - your address: %s\n", acc.Hex())
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to encrypt the new wallet")
	cmd.MarkFlagRequired("password")

	return cmd
}

func walletPrintPrivKeyCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pk-print",
		Short: "Unlocks keystore file and prints the Private + Public keys.",
		Run: func(cmd *cobra.Command, args []string) {
			keyJson, err := ioutil.ReadFile(flagKeystoreFile)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			key, err := keystore.DecryptKey(keyJson, flagPassword)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			spew.Dump(key)
		},
	}

	cmd.Flags().StringVar(&flagKeystoreFile, "keystore", "", "Absolute path to the encrypted keystore file")
	cmd.MarkFlagRequired("keystore")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to decrypt the wallet with")
	cmd.MarkFlagRequired("password")

	return cmd
}
