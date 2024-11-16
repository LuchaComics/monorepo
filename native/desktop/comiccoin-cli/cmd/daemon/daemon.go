package daemon

import (
	"context"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	auth_repo "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/repo"
	auth_usecase "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/common/preferences"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-cli/service"
)

var (
	preferences *pref.Preferences
)

// Command line argument flags
var (
	flagDataDirectory     string
	flagChainID           uint16
	flagAuthorityAddress  string
	flagNFTStorageAddress string
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

func DaemonCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "daemon",
		Short: "Runs a full node on your machine.",
		Run: func(cmd *cobra.Command, args []string) {
			doRunDaemon()
		},
	}
	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	return cmd
}

func doRunDaemon() {
	// ------ Common ------

	logger := logger.NewProvider()
	logger.Info("Syncing blockchain...",
		slog.Any("authority_address", flagAuthorityAddress))

	// ------------ Repo ------------

	blockchainStateChangeEventDTOConfigurationProvider := auth_repo.NewBlockchainStateChangeEventDTOConfigurationProvider(flagAuthorityAddress)
	blockchainStateChangeEventDTORepo := auth_repo.NewBlockchainStateChangeEventDTORepo(
		blockchainStateChangeEventDTOConfigurationProvider,
		logger)

	// ------------ Use-Case ------------

	// Blockchain State DTO
	subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase := auth_usecase.NewSubscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase(
		logger,
		blockchainStateChangeEventDTORepo)

	// ------------ Service ------------

	blockchainSyncManagerService := service.NewBlockchainSyncManagerService(
		logger,
		subscribeToBlockchainStateChangeEventsFromBlockchainAuthorityUseCase,
	)

	//
	// STEP X
	// Execute.
	//

	// Load up our operating system interaction handlers, more specifically
	// signals. The OS sends our application various signals based on the
	// OS's state, we want to listen into the termination signals.
	done := make(chan os.Signal, 1)
	signal.Notify(done, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL, syscall.SIGUSR1)

	go func() {
		ctx := context.Background()
		if err := blockchainSyncManagerService.Execute(ctx, flagChainID); err != nil {
			log.Fatalf("Failed to manage syncing: %v\n", err)
		}
	}()

	logger.Info("ComicCoin CLI daemon is running.")

	<-done
}
