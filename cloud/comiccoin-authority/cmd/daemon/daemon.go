package daemon

import (
	"log"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/blockchain/keystore"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/database/mongodb"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
)

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Run the comiccoin authority",
		Run: func(cmd *cobra.Command, args []string) {
			log.Println("Running daemon......")
			doRunDaemon()
		},
	}
	return cmd
}

func doRunDaemon() {
	// Common
	logger := logger.NewProvider()
	cfg := config.NewProvider()
	dbClient := mongodb.NewProvider(cfg, logger)
	keystore := keystore.NewAdapter(cfg, logger)

	// Repository
	walletRepo := repo.NewWalletRepo(cfg, logger, dbClient)
	accountRepo := repo.NewAccountRepo(cfg, logger, dbClient)

	_ = keystore
	_ = walletRepo
	_ = accountRepo
}
