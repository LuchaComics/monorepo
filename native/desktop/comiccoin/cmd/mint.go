package cmd

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
)

func init() {
	rootCmd.AddCommand(mintCmd())
}

func mintCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "mint",
		Short: "Creates a new coin for the address.",
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

			blochchain, err := blockchain.NewPoABlockchain(key)
			if err != nil {
				fmt.Println(err.Error())
				os.Exit(1)
			}

			spew.Dump(blochchain)
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagKeystoreFile, "keystore", "", "Absolute path to the encrypted keystore file")
	cmd.MarkFlagRequired("keystore")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to decrypt the wallet with")
	cmd.MarkFlagRequired("password")

	return cmd
}
