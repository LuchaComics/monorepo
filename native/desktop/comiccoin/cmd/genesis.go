package cmd

import (
	"context"
	"io/ioutil"
	"log"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	ledger_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/ledger/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
)

func init() {
	rootCmd.AddCommand(genesisCmd())
}

func genesisCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "genesis",
		Short: "Initialize `ComicCoin` ledger by creating the genesis block",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()

			logger.Info("Creating genesis block...")

			//
			// STEP 1
			// Load up a wallet which has coins in it.
			//

			coinbaseKeyJson, err := ioutil.ReadFile(flagKeystoreFile)
			if err != nil {
				log.Fatalf("failed reading file: %v", err)
			}

			coinbaseKey, err := keystore.DecryptKey(coinbaseKeyJson, flagPassword)
			if err != nil {
				log.Fatalf("failed decrypting file: %v", err)
			}

			logger.Info("Coinbase wallet was successfully opened")

			//
			// STEP 2
			// Load up our ledger.
			//

			// Load up the configuration.
			cfg := &config.Config{
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				BlockchainDifficulty: 1,
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}

			// Load up our dependencies
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			lastHashDS := lasthash_ds.NewDatastore(cfg, logger, kvs)
			blockDS := block_ds.NewDatastore(cfg, logger, kvs)
			ledgerController := ledger_c.NewController(cfg, logger, lastHashDS, blockDS)

			// Generate our genesis
			ctx := context.Background()
			genesisBlock, err := ledgerController.NewGenesisBlock(ctx, coinbaseKey)
			if err != nil {
				log.Fatalf("failed to generate genesis block: %v", err)

			}

			logger.Info("Genesis block created successfully",
				slog.String("hash", genesisBlock.Hash))

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
