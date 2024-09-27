package cmd

import (
	"context"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	kvs "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore/leveldb"
	local_pubsub "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/pubsub/local"
	p2p_pubsub "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/pubsub/p2p"
	acc_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/account/datastore"
	blockchain_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain/controller"
	blockdata_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockdata/datastore"
	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	lasthash_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/lasthash/datastore"
	pt_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/signedtransaction/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config/constants"
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
				Blockchain: config.BlockchainConfig{
					ChainID:       constants.ChainIDMainNet,
					TransPerBlock: 1,
					Difficulty:    1,
				},
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				Peer: config.PeerConfig{
					ListenPort:       flagListenPeerToPeerPort,
					KeyName:          flagKeypairName,
					RendezvousString: flagRendezvousString,
					// BootstrapPeers:   bootstrapPeers,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}

			// Load up our dependencies
			kvs := kvs.NewKeyValueStorer(cfg, logger)
			uuid := uuid.NewProvider()
			keypairDS := keypair_ds.NewDatastore(cfg, logger, kvs)
			localPubSubBroker := local_pubsub.NewAdapter(cfg, logger)
			p2pPubSubBroker := p2p_pubsub.NewAdapter(cfg, logger, keypairDS)
			lastHashDS := lasthash_ds.NewDatastore(cfg, logger, kvs)
			accountStorer := acc_s.NewDatastore(cfg, logger, kvs)
			signedTransactionDS := pt_ds.NewDatastore(cfg, logger, kvs)
			blockDataDS := blockdata_ds.NewDatastore(cfg, logger, kvs)
			blockchainController := blockchain_c.NewController(cfg, logger, uuid, localPubSubBroker, p2pPubSubBroker, accountStorer, signedTransactionDS, lastHashDS, blockDataDS)

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
				slog.String("zero_hash", genesisBlock.Hash)) // Not a bug, genesis is always zero.
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.Flags().StringVar(&flagKeypairName, "keypair-name", "", "The name of keypairs to apply to this server")
	cmd.MarkFlagRequired("keypair-name")
	cmd.Flags().StringVar(&flagAccountName, "coinbase-account-name", "", "The account name of the coinbase wallet")
	cmd.MarkFlagRequired("coinbase-account-name")
	cmd.Flags().StringVar(&flagPassword, "coinbase-password", "", "The password to decrypt the cointbase's account wallet")
	cmd.MarkFlagRequired("coinbase-password")

	return cmd
}
