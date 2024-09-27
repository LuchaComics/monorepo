package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"log/slog"

	"github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	acc_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
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
			// STEP 1
			// Load up our dependencies.
			//

			cfg := &config.Config{
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			accountStorer := acc_s.NewDatastore(cfg, logger, kvs)

			//
			// STEP 2
			// Create our wallet in our filesystem.
			//

			account, insertErr := accountStorer.Insert(context.Background(), flagAccountName, flagPassword)
			if insertErr != nil {
				logger.Error("failed inserting new account into database",
					slog.Any("name", flagAccountName),
					slog.Any("error", insertErr))
				log.Fatalf("failed inserting into database: %s", insertErr)
			}
			if account == nil {
				log.Fatal("account does not exist")
			}
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagAccountName, "account-name", "", "The name to assign this account")
	cmd.MarkFlagRequired("account-name")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to encrypt the new wallet")
	cmd.MarkFlagRequired("password")

	return cmd
}

func walletPrintPrivKeyCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "pk-print",
		Short: "Unlocks keystore file and prints the Private + Public keys.",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies.
			//

			cfg := &config.Config{
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			accountStorer := acc_s.NewDatastore(cfg, logger, kvs)

			//
			// STEP 2
			// Read the contents of the keystore.
			//

			account, err := accountStorer.GetByName(context.Background(), flagAccountName)
			if err != nil {
				log.Fatalf("failed getting account by name: %v", err)
			}
			if account == nil {
				log.Fatalf("failed getting account by name because it d.n.e for name: %s", flagAccountName)
			}

			keyJson, err := ioutil.ReadFile(account.WalletFilepath)
			if err != nil {
				log.Fatalf("failed reading file: %v", err)
			}

			key, err := keystore.DecryptKey(keyJson, flagPassword)
			if err != nil {
				log.Fatalf("failed decrypting file: %v", err)
			}

			spew.Dump(key)
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagAccountName, "account-name", "", "The name to assign this account")
	cmd.MarkFlagRequired("account-name")
	cmd.Flags().StringVar(&flagPassword, "password", "", "The password to encrypt the new wallet")
	cmd.MarkFlagRequired("password")

	return cmd
}
