package cmd

import (
	"context"
	"log"
	"log/slog"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	ledger_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/ledger/controller"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(balanceCmd)
	balanceCmd.AddCommand(balanceGetCmd())
	// balanceCmd.AddCommand(balancePrintCmd())
}

var balanceCmd = &cobra.Command{
	Use:   "balance",
	Short: "Interacts with balance (new...).",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func balanceGetCmd() *cobra.Command {
	var balanceGetCmd = &cobra.Command{
		Use:   "get",
		Short: "Get balance of address",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies and configuration
			//

			cfg := &config.Config{
				BlockchainDifficulty: 1,
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			lastHashDS := lasthash_ds.NewDatastore(cfg, logger, kvs)
			blockDS := block_ds.NewDatastore(cfg, logger, kvs)
			ledgerController := ledger_c.NewController(cfg, logger, lastHashDS, blockDS)

			// // Generate our new balance
			// logger.Info("Creating balance...")
			// if err := balanceDS.GenerateNewKeyPairAndSetByName(ctx, flagKeypairName); err != nil {
			// 	log.Fatalf("failed generating balance: %v\n", err)
			// }
			// logger.Info("Keypair created!")

			//
			// STEP 2
			// Lookup balance.
			//

			ctx := context.Background()
			address := common.HexToAddress(flagRecipientAddress)
			balance, err := ledgerController.GetBalanceByAddress(ctx, address)
			if err != nil {
				log.Fatalf("Failed to get balance: %v", err)
			}
			logger.Info("Fetched balance", slog.Any("amount", balance))

		},
	}
	balanceGetCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	balanceGetCmd.Flags().StringVar(&flagRecipientAddress, "address", "", "The address of the coin(s) receipient")
	balanceGetCmd.MarkFlagRequired("address")

	return balanceGetCmd
}

// func balancePrintCmd() *cobra.Command {
// 	var balanceGetCmd = &cobra.Command{
// 		Use:   "print",
// 		Short: "Prints balance",
// 		Run: func(cmd *cobra.Command, args []string) {
// 			// Load up our dependencies and configuration
// 			cfg := &config.Config{
// 				BlockchainDifficulty: 1,
// 				DB: config.DBConfig{
// 					DataDir: flagDataDir,
// 				},
// 			}
// 			logger := logger.NewProvider()
// 			kvs := kvs.NewKeyValueStorer(cfg, logger)
// 			balanceDS := balance_ds.NewDatastore(cfg, logger, kvs)
//
// 			// Generate our new balance
// 			logger.Info("Getting balance...")
// 			ctx := context.Background()
// 			priv, pub, err := balanceDS.GetByName(ctx, flagKeypairName)
// 			if err != nil {
// 				log.Fatalf("failed getting balance: %v\n", err)
// 			}
//
// 			spew.Dump(priv)
// 			spew.Dump(pub)
// 		},
// 	}
// 	balanceGetCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
// 	// cmd.MarkFlagRequired("datadir")
// 	balanceGetCmd.Flags().StringVar(&flagKeypairName, "name", "", "The name to assign the balances so you can reference in the future")
// 	balanceGetCmd.MarkFlagRequired("name")
//
// 	return balanceGetCmd
// }
