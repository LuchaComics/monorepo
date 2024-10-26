package cmd

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/usecase"
)

// Command line argument flags
var (
	flagMetadataURI string
)

func GetMetadataCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "getmetadata",
		Short: "Commands used to get the nft metadata",
		Run: func(cmd *cobra.Command, args []string) {
			doGetMetadataCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flagMetadataURI, "uri", "", "The URI of the NFT asset")
	cmd.MarkFlagRequired("uri")

	return cmd
}

func doGetMetadataCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	logger := logger.NewLogger()
	logger.Info("Fetching metadata",
		slog.String("uri", flagMetadataURI))

	// Lookup data from a site using either IPFS or HTTPS.
	getNFTMetadataUseCase := usecase.NewGetNFTMetadataUseCase(logger)
	metadata, err := getNFTMetadataUseCase.Execute(flagMetadataURI)
	if err != nil {
		log.Fatalf("Failed nft metadata: %v", err)
	}

	logger.Debug("fetched",
		slog.Any("metadata", metadata))

}
