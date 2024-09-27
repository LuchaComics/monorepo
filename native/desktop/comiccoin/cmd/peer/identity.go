package peer

import (
	"log"
	"log/slog"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/config"
	ik_repo "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/repo"
	ik_s "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/service"
	ik_use "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/peer/usecase"
	dbase "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/db/leveldb"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/logger"
)

func identityCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "identity",
		Short: "Peer identity related commands",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	cmd.AddCommand(newIdentityCmd())
	cmd.AddCommand(getIdentityCmd())

	return cmd
}

func newIdentityCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "new",
		Short: "Creates new identity key",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies and configuration
			//

			cfg := &config.Config{
				App: config.AppConfig{
					HTTPPort: flagListenHTTPPort,
					HTTPIP:   flagListenHTTPIP,
					DirPath:  flagDataDir,
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

			ik, err := ikCreateService.Execute(flagIdentityKeyID)
			if err != nil {
				log.Fatalf("Failed creating identity key: %v", err)
			}
			if ik == nil {
				log.Fatal("Failed creating identity key: d.n.e.")
			}
			logger.Debug("Created identity key", slog.Any("result", ik))

		},
	}
	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagIdentityKeyID, "id", "", "The id to assign the identity key so you can reference throughout the app")
	cmd.MarkFlagRequired("id")

	return cmd
}

func getIdentityCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "get",
		Short: "Gets exisiting identity key or errors if d.n.e.",
		Run: func(cmd *cobra.Command, args []string) {
			//
			// STEP 1
			// Load up our dependencies and configuration
			//

			cfg := &config.Config{
				App: config.AppConfig{
					HTTPPort: flagListenHTTPPort,
					HTTPIP:   flagListenHTTPIP,
					DirPath:  flagDataDir,
				},
				DB: config.DBConfig{
					DataDir: flagDataDir,
				},
			}
			logger := logger.NewLogger()
			db := dbase.NewDatabase(cfg.DB.DataDir, logger)
			ikRepo := ik_repo.NewIdentityKeyRepo(cfg, logger, db)
			ikGetUseCase := ik_use.NewGetIdentityKeyUseCase(cfg, logger, ikRepo)
			ikGetService := ik_s.NewGetIdentityKeyService(cfg, logger, ikGetUseCase)

			ik, err := ikGetService.Execute(flagIdentityKeyID)
			if err != nil {
				log.Fatalf("Failed getting identity key: %v", err)
			}
			if ik == nil {
				log.Fatal("Failed getting identity key: d.n.e.")
			}
			logger.Debug("Identity key found", slog.Any("result", ik))

		},
	}
	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagIdentityKeyID, "id", "", "The id to assign the identity key so you can reference throughout the app")
	cmd.MarkFlagRequired("id")

	return cmd
}
