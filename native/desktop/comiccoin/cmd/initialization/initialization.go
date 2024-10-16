package initialization

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config/constants"
	ik_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	ik_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	ik_use "github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/storage/disk/leveldb"
	ma "github.com/multiformats/go-multiaddr"
)

var (
	flagDataDir       string // Location of the database directory
	flagIdentityKeyID string // Custom profile id.
)

func InitCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes the peer-to-peer network configuration for this blockchain node",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies and configuration
			//

			// We can setup minimal settings as the systems affect don't have
			// much configuration to deal with.
			cfg := &config.Config{
				App: config.AppConfig{
					DirPath: flagDataDir,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
				Peer: config.PeerConfig{
					ListenPort:       26642, // Our application port
					KeyName:          "",
					RendezvousString: "",
					BootstrapPeers:   nil,
				},
			}
			logger := logger.NewLogger()
			db := disk.NewDiskStorage(cfg.DB.DataDir, "identity_key", logger)
			ikRepo := ik_repo.NewIdentityKeyRepo(cfg, logger, db)
			ikCreateUseCase := ik_use.NewCreateIdentityKeyUseCase(cfg, logger, ikRepo)
			ikGetUseCase := ik_use.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
			ikCreateService := ik_s.NewCreateIdentityKeyService(cfg, logger, ikCreateUseCase, ikGetUseCase)

			// If nothing was set then we use a default value. We do this to
			// simplify the user's experience.
			if flagIdentityKeyID == "" {
				flagIdentityKeyID = constants.DefaultIdentityKeyID
			}

			logger.Info("Blockchain node intitializing...")

			identityKey, err := ikCreateService.Execute(flagIdentityKeyID)
			if err != nil {
				log.Fatalf("Failed creating identity key: %v", err)
			}
			if identityKey == nil {
				log.Fatal("Failed creating identity key: d.n.e.")
			}

			privateKey, _ := identityKey.GetPrivateKey()
			publicKey, _ := identityKey.GetPublicKey()
			libP2PNetwork := p2p.NewLibP2PNetwork(cfg, logger, privateKey, publicKey)
			h := libP2PNetwork.GetHost()

			// Build host multiaddress
			hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", h.ID()))

			// Now we can build a full multiaddress to reach this host
			// by encapsulating both addresses:
			addr := h.Addrs()[0]
			fullAddr := addr.Encapsulate(hostAddr)

			logger.Info("Blockchain node intitialized and ready",
				slog.Any("peer identity", h.ID()),
				slog.Any("full address", fullAddr),
			)
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagIdentityKeyID, "id", "", "You can override the blockchain node's identity by setting a custom profile id on startup, you will need to reference it later when you run the daemon.")

	return cmd
}
