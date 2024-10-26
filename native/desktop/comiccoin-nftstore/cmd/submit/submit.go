package submit

import (
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/logger"
	disk "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/storage/disk/leveldb"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/usecase"
)

// Command line argument flags
var (
	flagMetadataURI string
	flagDataDir     string // Location of the database directory
)

func SubmitMetadataURICmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "submit",
		Short: "Commands used for submitting NFT metadata URI",
		Run: func(cmd *cobra.Command, args []string) {
			doSubmitCmd()
		},
	}

	cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your store's data dir where the assets will be/are stored")
	cmd.Flags().StringVar(&flagMetadataURI, "metadata-uri", "", "The URI of the NFT metadata")
	cmd.MarkFlagRequired("metadata-uri")

	return cmd
}

func doSubmitCmd() {
	//
	// STEP 1
	// Load up our dependencies and configuration
	//

	logger := logger.NewLogger()
	logger.Info("Submitting",
		slog.String("metadata-uri", flagMetadataURI))

	nftByTokenIDDB := disk.NewDiskStorage(flagDataDir, "nft_by_tokenid", logger)
	nftByMetadataURIDB := disk.NewDiskStorage(flagDataDir, "nft_by_metadatauri", logger)
	nftRepo := repo.NewNFTRepo(logger, nftByTokenIDDB, nftByMetadataURIDB)
	_ = nftRepo
	getNFTMetadataUseCase := usecase.NewGetNFTMetadataUseCase(logger)
	metadata, err := getNFTMetadataUseCase.Execute(flagMetadataURI)
	if err != nil {
		log.Fatalf("Failed downloading metadata via uri: %v\n", err)
	}
	logger.Debug("Downloaded", slog.Any("metadata", metadata))
	//TODO: Impl.
}
