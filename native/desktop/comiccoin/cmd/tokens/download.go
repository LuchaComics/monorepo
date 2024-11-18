package tokens

import (
	"context"
	"log"
	"log/slog"
	"math/big"
	"strings"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	disk "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/storage/disk/leveldb"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

// Command line argument flags
var ()

func DownloadTokenCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "download",
		Short: "Download the token contents to your local machine from the blockchain network.",
		Run: func(cmd *cobra.Command, args []string) {
			doRunDownloadTokenCommand()
		},
	}

	cmd.Flags().StringVar(&flagDataDirectory, "data-directory", preferences.DataDirectory, "The data directory to save to")
	cmd.Flags().Uint16Var(&flagChainID, "chain-id", preferences.ChainID, "The blockchain to sync with")
	cmd.Flags().StringVar(&flagAuthorityAddress, "authority-address", preferences.AuthorityAddress, "The BlockChain authority address to connect to")
	cmd.Flags().StringVar(&flagNFTStorageAddress, "nftstorage-address", preferences.NFTStorageAddress, "The NFT storage service adress to connect to")

	cmd.Flags().StringVar(&flagTokenID, "token-id", "", "The unique token identification to use to lookup the token")
	cmd.MarkFlagRequired("token-id")

	return cmd
}

func doRunDownloadTokenCommand() {
	// ------ Common ------

	logger := logger.NewProvider()
	tokenDB := disk.NewDiskStorage(flagDataDirectory, "token", logger)

	// ------ Repo ------

	tokenRepo := repo.NewTokenRepo(
		logger,
		tokenDB)

	nftAssetRepoConfig := repo.NewNFTAssetRepoConfigurationProvider(flagNFTStorageAddress, "")
	nftAssetRepo := repo.NewNFTAssetRepo(nftAssetRepoConfig, logger)

	// ------ Use-case ------

	getTokenUseCase := usecase.NewGetTokenUseCase(
		logger,
		tokenRepo,
	)

	// ------ Service ------

	ctx := context.Background()
	tokenGetService := service.NewTokenGetService(
		logger,
		getTokenUseCase,
	)

	// ------ Execute ------

	tokenID, ok := new(big.Int).SetString(flagTokenID, 10)
	if !ok {
		log.Fatal("Failed convert `token_id` to big.Int")
	}

	logger.Debug("Token retrieving...",
		slog.Any("token_id", tokenID))

	tok, retrieveErr := tokenGetService.Execute(
		ctx,
		tokenID,
	)
	if retrieveErr != nil {
		log.Fatalf("Failed execute get token: %v", retrieveErr)
	}
	if tok == nil {
		logger.Warn("Token does not exist",
			slog.Any("token_id", tokenID))
	} else {
		logger.Debug("Retrieved token",
			slog.Any("token_id", tokenID),
			slog.Any("owner", tok.Owner),
			slog.Any("metadata_uri", tok.MetadataURI),
			slog.Any("nonce", tok.GetNonce()))

		ctx := context.Background()
		metadataCid := strings.Replace(tok.MetadataURI, "ipfs://", "", -1)
		metadata, err := nftAssetRepo.Get(ctx, metadataCid)
		if err != nil {
			log.Fatalf("Failed getting nft metadata: %v", err)
		}
		//TODO: IMPLEMENT
	}
}
