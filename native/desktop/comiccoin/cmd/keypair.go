package cmd

import (
	"context"
	"log"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
	"github.com/davecgh/go-spew/spew"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(keypairCmd)
	keypairCmd.AddCommand(keypairNewCmd())
	keypairCmd.AddCommand(keypairPrintCmd())
}

var keypairCmd = &cobra.Command{
	Use:   "keypair",
	Short: "Interacts with keypair (new...).",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func keypairNewCmd() *cobra.Command {
	var keypairGetCmd = &cobra.Command{
		Use:   "new",
		Short: "Create new keypair",
		Run: func(cmd *cobra.Command, args []string) {
			// Load up our dependencies and configuration
			cfg := &config.Config{
				Blockchain: config.BlockchainConfig{
					ChainID:    constants.ChainIDMainNet,
					Difficulty: 1,
				},
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			keypairDS := keypair_ds.NewDatastore(cfg, logger, kvs)

			// Generate our new keypair
			logger.Info("Creating keypair...")
			ctx := context.Background()
			if err := keypairDS.GenerateNewKeyPairAndSetByName(ctx, flagKeypairName); err != nil {
				log.Fatalf("failed generating keypair: %v\n", err)
			}
			logger.Info("Keypair created!")
		},
	}
	keypairGetCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	keypairGetCmd.Flags().StringVar(&flagKeypairName, "name", "", "The name to assign the keypairs so you can reference in the future")
	keypairGetCmd.MarkFlagRequired("name")

	return keypairGetCmd
}

func keypairPrintCmd() *cobra.Command {
	var keypairGetCmd = &cobra.Command{
		Use:   "print",
		Short: "Prints keypair",
		Run: func(cmd *cobra.Command, args []string) {
			// Load up our dependencies and configuration
			cfg := &config.Config{
				Blockchain: config.BlockchainConfig{
					ChainID:    constants.ChainIDMainNet,
					Difficulty: 1,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewProvider()
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			keypairDS := keypair_ds.NewDatastore(cfg, logger, kvs)

			// Generate our new keypair
			logger.Info("Getting keypair...")
			ctx := context.Background()
			priv, pub, err := keypairDS.GetByName(ctx, flagKeypairName)
			if err != nil {
				log.Fatalf("failed getting keypair: %v\n", err)
			}

			spew.Dump(priv)
			spew.Dump(pub)
		},
	}
	keypairGetCmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// cmd.MarkFlagRequired("datadir")
	keypairGetCmd.Flags().StringVar(&flagKeypairName, "name", "", "The name to assign the keypairs so you can reference in the future")
	keypairGetCmd.MarkFlagRequired("name")

	return keypairGetCmd
}
