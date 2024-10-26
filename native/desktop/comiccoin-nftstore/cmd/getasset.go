package cmd

import (
	"io"
	"log"
	"log/slog"
	"os"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/usecase"
)

// Command line argument flags
var (
	flagURI string
)

func GetAssetCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "getasset",
		Short: "Commands used to get the asset",
		Run: func(cmd *cobra.Command, args []string) {
			doGetAssetCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flagURI, "uri", "", "The URI of the NFT asset")
	cmd.MarkFlagRequired("uri")

	return cmd
}

func doGetAssetCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	logger := logger.NewLogger()
	logger.Info("Fetching",
		slog.String("uri", flagURI))

	// Lookup data from a site using either IPFS or HTTPS.
	getNFTAssetUseCase := usecase.NewGetNFTAssetUseCase(logger)
	resp, err := getNFTAssetUseCase.Execute(flagURI)
	if err != nil {
		log.Fatalf("Failed nft asset: %v", err)
	}

	logger.Debug("fetched asset",
		slog.String("filename", resp.Filename),
		slog.String("file_ext", resp.FileExtension),
		slog.String("content_type", resp.ContentType),
		slog.Int64("content_length", resp.ContentLength))

	// Save the data to file.
	f, err := os.Create(resp.Filename)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Content)
	if err != nil {
		log.Fatal(err)
	}
}
