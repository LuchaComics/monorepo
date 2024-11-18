package tokens

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"math/big"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/logger"
	"github.com/spf13/cobra"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/repo"
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

	logger := logger.NewProvider()
	comicCoincRPCClientRepoConfigurationProvider := repo.NewComicCoincRPCClientRepoConfigurationProvider("localhost", "2233")
	rpcClient := repo.NewComicCoincRPCClientRepo(comicCoincRPCClientRepoConfigurationProvider, logger)

	ctx := context.Background()

	tokID, ok := new(big.Int).SetString(flagTokenID, 10)
	if !ok {
		log.Fatal("Failed convert `token_id` to big.Int")
	}

	nftok, err := rpcClient.GetNonFungibleToken(ctx, tokID, flagDataDirectory)
	if err != nil {
		log.Fatalf("Failed downloading non-fungible tokens: %v\n", err)
	}
	logger.Debug(fmt.Sprintf("NFT contents found at path: %v", flagDataDirectory),
		slog.Any("token_id", nftok.TokenID),
		slog.Any("metadata_uri", nftok.MetadataURI),
		slog.Any("metadata", nftok.Metadata),
	)
}
