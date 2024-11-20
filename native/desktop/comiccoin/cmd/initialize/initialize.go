package initialize

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"

	pref "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/preferences"
)

var (
	preferences *pref.Preferences
)

// Initialize function will be called when every command gets called.
func init() {
	preferences = pref.PreferencesInstance()
}

const (
	ChainIDMainNet = 1
)

// Command line argument flags
var (
	flagDataDirectory     string
	flagChainID           uint16
	flagAuthorityAddress  string
	flagNFTStorageAddress string
)

func InitializeCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "init",
		Short: "Initialize your local ComicCoin for the first time. Note: This action cannot be undone once executed.",
		Run: func(cmd *cobra.Command, args []string) {
			logger := logger.NewProvider()

			logger.Debug("Configured ComicCoin...",
				slog.Any("DataDirectory", preferences.DataDirectory),
				slog.Any("ChainID", preferences.ChainID),
				slog.Any("AuthorityAddress", preferences.AuthorityAddress),
				slog.Any("NFTStorageAddress", preferences.NFTStorageAddress),
				slog.Any("FilePathPreferences", preferences.GetFilePathOfPreferencesFile()))

			if preferences.DataDirectory != "" {
				log.Fatalf("You have already configured ComicCoin: DataDirectory was set with: %v\n", preferences.DataDirectory)
			}
			preferences.SetDataDirectory(flagDataDirectory)

			if preferences.ChainID > 0 {
				log.Fatalf("You have already configured ComicCoin: ChainID was set with: %v\n", preferences.ChainID)
			}
			preferences.SetChainID(flagChainID)

			if preferences.AuthorityAddress != "" {
				log.Fatalf("You have already configured ComicCoin: AuthorityAddress was set with: %v\n", preferences.AuthorityAddress)
			}
			preferences.SetAuthorityAddress(flagAuthorityAddress)

			if preferences.NFTStorageAddress != "" {
				log.Fatalf("You have already configured ComicCoin: NFTStorageAddress was set with: %v\n", preferences.NFTStorageAddress)
			}
			preferences.SetNFTStorageAddress(flagNFTStorageAddress)

			logger.Debug("Configued ComicCoin",
				slog.Any("DataDirectory", preferences.DataDirectory),
				slog.Any("ChainID", preferences.ChainID),
				slog.Any("AuthorityAddress", preferences.AuthorityAddress),
				slog.Any("NFTStorageAddress", preferences.NFTStorageAddress),
				slog.Any("FilePathPreferences", preferences.GetFilePathOfPreferencesFile()))
		},
	}
	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", pref.GetDefaultDataDirectory(), "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", ChainIDMainNet, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", "https://comiccoinauthority.com", "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", "https://comiccoinnftstorage.com", "The NFT storage service adress to connect to")

	return cmd
}
