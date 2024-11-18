package tokens

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/big"

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
	nftokDB := disk.NewDiskStorage(flagDataDirectory, "non_fungible_token", logger)

	// ------ Repo ------

	tokenRepo := repo.NewTokenRepo(
		logger,
		tokenDB)
	nftokenRepo := repo.NewNonFungibleTokenRepo(logger, nftokDB)
	nftAssetRepoConfig := repo.NewNFTAssetRepoConfigurationProvider(flagNFTStorageAddress, "")
	nftAssetRepo := repo.NewNFTAssetRepo(nftAssetRepoConfig, logger)

	// ------ Use-case ------

	getTokUseCase := usecase.NewGetTokenUseCase(
		logger,
		tokenRepo)
	getNFTokUseCase := usecase.NewGetNonFungibleTokenUseCase(
		logger,
		nftokenRepo)
	downloadNFTokMetadataUsecase := usecase.NewDownloadMetadataNonFungibleTokenUseCase(
		logger,
		nftAssetRepo)
	downloadNFTokAssetUsecase := usecase.NewDownloadNonFungibleTokenAssetUseCase(
		logger,
		nftAssetRepo)
	upsertNFTokUseCase := usecase.NewUpsertNonFungibleTokenUseCase(
		logger,
		nftokenRepo)

	// ------ Service ------

	_ = getNFTokUseCase

	getOrDownloadNonFungibleTokenService := service.NewGetOrDownloadNonFungibleTokenService(
		logger,
		getNFTokUseCase,
		getTokUseCase,
		downloadNFTokMetadataUsecase,
		downloadNFTokAssetUsecase,
		upsertNFTokUseCase)

	// ------ Execute ------

	ctx := context.Background()

	//
	// STEP 2
	// Check if we can connect with IPFS node.
	//

	version, err := nftAssetRepo.Version(ctx)
	if err != nil {
		log.Fatalf("Failed connecting to NFT assets store, you are not connected.")
	}
	fmt.Printf("NFT Assets Store Version: %s\n", version)

	//
	// STEP 3
	// Lookup our `token id` in our NFT db and if it exists we can
	// exit this command as we've already downloaded the data.
	//

	tokID, ok := new(big.Int).SetString(flagTokenID, 10)
	if !ok {
		log.Fatal("Failed convert `token_id` to big.Int")
	}

	nftok, err := getOrDownloadNonFungibleTokenService.Execute(ctx, tokID, flagDataDirectory)
	if err != nil {
		log.Fatalf("Failed downloading non-fungible tokens: %v\n", err)
	}

	logger.Debug("Downloaded NFT successfully.",
		slog.Any("token_id", nftok.TokenID),
		slog.Any("metadata_uri", nftok.MetadataURI),
		slog.Any("metadata", nftok.Metadata),
	)
}
