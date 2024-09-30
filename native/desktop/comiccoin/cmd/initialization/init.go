package initialization

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config/constants"
	ik_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/repo"
	ik_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/service"
	ik_use "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
)

var (
	flagDataDir       string // Location of the database directory
	flagIdentityKeyID string // Custom profile id.
)

func InitCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initializes the blockchain node",
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
			}
			logger := logger.NewLogger()
			db := dbase.NewDatabase(cfg.DB.DataDir, logger)
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

			ik, err := ikCreateService.Execute(flagIdentityKeyID)
			if err != nil {
				log.Fatalf("Failed creating identity key: %v", err)
			}
			if ik == nil {
				log.Fatal("Failed creating identity key: d.n.e.")
			}
			logger.Info("Blockchain node intitialized and ready")
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagIdentityKeyID, "id", "", "You can override the blockchain node's identity by setting a custom profile id on startup, you will need to reference it later when you run the daemon.")

	return cmd
}
