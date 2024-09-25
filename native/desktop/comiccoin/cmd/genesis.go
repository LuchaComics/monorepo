package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	mqb "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/messagequeuebroker/simple"
	acc_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	block_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/block/datastore"
	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/logger"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/provider/uuid"
)

func init() {
	rootCmd.AddCommand(genesisCmd())
}

func genesisCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "genesis",
		Short: "Initialize `ComicCoin` blockchain by creating the genesis block",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies
			//

			logger := logger.NewProvider()

			logger.Info("Creating genesis block...")

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
			uuid := uuid.NewProvider()
			broker := mqb.NewMessageQueue(cfg, logger)
			lastHashDS := lasthash_ds.NewDatastore(cfg, logger, kvs)
			accountStorer := acc_s.NewDatastore(cfg, logger, kvs)
			signedTransactionDS := pt_ds.NewDatastore(cfg, logger, kvs)
			blockDS := block_ds.NewDatastore(cfg, logger, kvs)
			blockchainController := blockchain_c.NewController(cfg, logger, uuid, broker, accountStorer, signedTransactionDS, lastHashDS, blockDS)

			//
			// STEP 2
			// Read the contents of the keystore.
			//

			coinbaseKey, err := accountStorer.GetKeyByNameAndPassword(context.Background(), flagAccountName, flagPassword)
			if err != nil {
				log.Fatalf("failed getting key by account by name: %v", err)
			}
			if coinbaseKey == nil {
				log.Fatalf("failed getting key by account by name because name d.n.e.: %s", flagAccountName)
			}

			//
			// STEP 3
			// Generate our genesis
			//

			ctx := context.Background()
			genesisBlock, err := blockchainController.NewGenesisBlock(ctx, coinbaseKey)
			if err != nil {
				log.Fatalf("failed to generate genesis block: %v", err)

			}

			logger.Info("Genesis block created successfully",
				slog.String("hash", genesisBlock.Hash))
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagAccountName, "coinbase-account-name", "", "The account name of the coinbase wallet")
	cmd.MarkFlagRequired("coinbase-account-name")
	cmd.Flags().StringVar(&flagPassword, "coinbase-password", "", "The password to decrypt the cointbase's account wallet")
	cmd.MarkFlagRequired("coinbase-password")

	return cmd
}
