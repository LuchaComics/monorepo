package service

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

// NonFungibleTokenAssetsService struct represents the service which
// handles polling the local database for any new issued (and signed) tokens
// recently saved and then proceed to download that token's assets to our
// computer.
type NonFungibleTokenAssetsService struct {
	config                          *config.Config
	logger                          *slog.Logger
	listAllSignedIssuedTokenUseCase *usecase.ListAllSignedIssuedTokenUseCase
	getNFTokUseCase                 *usecase.GetNonFungibleTokenUseCase
	downloadNFTokMetadataUsecase    *usecase.DownloadMetadataNonFungibleTokenUseCase
	downloadNFTokAssetUsecase       *usecase.DownloadNonFungibleTokenAssetUseCase
	upsertNFTokUseCase              *usecase.UpsertNonFungibleTokenUseCase
}

func NewNonFungibleTokenAssetsService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ListAllSignedIssuedTokenUseCase,
	uc2 *usecase.GetNonFungibleTokenUseCase,
	uc3 *usecase.DownloadMetadataNonFungibleTokenUseCase,
	uc4 *usecase.DownloadNonFungibleTokenAssetUseCase,
	uc5 *usecase.UpsertNonFungibleTokenUseCase,
) *NonFungibleTokenAssetsService {
	return &NonFungibleTokenAssetsService{cfg, logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *NonFungibleTokenAssetsService) Execute() error {
	//
	// STEP 1
	// Iterate through all the signed tokens we have locally saved on
	// this machine.
	//

	sitoks, err := s.listAllSignedIssuedTokenUseCase.Execute()
	if err != nil {
		s.logger.Error("Failed listing all signed issued tokens.",
			slog.Any("error", err))
		return err
	}

	for _, sitok := range sitoks {
		//
		// STEP 2:
		// Check to see if we already have a non-fungible token saved locally
		// for the given signed issued token: If already saved then skip this
		// loop iteration because we already have the NFT assets, else we need
		// to go and get the NFT assets.
		//

		nftok, err := s.getNFTokUseCase.Execute(sitok.ID)
		if err != nil {
			s.logger.Error("Failed getting non-fungible token.",
				slog.Any("tokenID", sitok.ID),
				slog.Any("error", err))
			return err
		}
		if nftok != nil {
			// Remember that if an NFT already exists, then we already have
			// the digital assets saved locally on this machine, and the
			// purpose of this service is to download NFT assets which we don't
			// have so therefore we can skip the rest of this function.
			return nil
		}

		//
		// STEP 3: CASE A:
		// Confirm URI is using protocol our app supports and proceed to
		// download the assets based on that protocol.
		//

		if strings.Contains(sitok.MetadataURI, "ipfs://") == true {
			return s.downloadNFTAssets(sitok)
		} else if strings.Contains(sitok.MetadataURI, "https://") == true {
			return s.downloadNFTAssets(sitok)
		}

		//
		// STEP 3: CASE B:
		// If our URI does not have a supported protocol then abort with error.
		//

		s.logger.Error("Token metadata URI contains protocol we do not support:",
			slog.Any("tokenID", sitok.ID),
			slog.Any("metadataURI", sitok.MetadataURI))

		return fmt.Errorf("Token metadata URI contains protocol we do not support: %v\n", sitok.MetadataURI)
	}

	return nil
}

func (s *NonFungibleTokenAssetsService) downloadNFTAssets(sitok *domain.SignedIssuedToken) error {
	// Variable holds the path to the NFT's metadata file on the network.
	// Later in this function we will replace this URI with the filepath
	// of were the metadata is saved locally on our machine.
	metadataURI := sitok.MetadataURI

	//
	// STEP 1:
	// Download the metadata from the "network".
	//

	metadata, metadataFilepath, err := s.downloadNFTokMetadataUsecase.Execute(sitok.ID, metadataURI)
	if err != nil {
		s.logger.Error("Failed downloading nft metadata.",
			slog.Any("tokenID", sitok.ID))
		return fmt.Errorf("Failed downloading nft metadata: %v\n", err)
	}

	// Replace the IPFS path with our local systems filepath.
	metadataURI = metadataFilepath

	//
	// STEP 2
	// Create our NFT token to be referenced in future.
	//

	nftok := &domain.NonFungibleToken{
		TokenID:     sitok.ID,
		MetadataURI: metadataURI,
		Metadata:    metadata,
		State:       domain.NonFungibleTokenStateNotReady,
	}

	if err := s.upsertNFTokUseCase.Execute(nftok); err != nil {
		s.logger.Error("Failed creating nft token.",
			slog.Any("tokenID", sitok.ID))
		return fmt.Errorf("Failed creating nft token: %v\n", err)
	}

	//
	// STEP 3
	// Download the image file from the "network" and save locally.
	//

	imageFilepath, err := s.downloadNFTokAssetUsecase.Execute(sitok.ID, metadata.Image)
	if err != nil {
		s.logger.Error("Failed downloading nft image asset.",
			slog.Any("tokenID", sitok.ID),
			slog.Any("Image", metadata.Image),
			slog.Any("err", err))
		return fmt.Errorf("Failed downloading nft image asset: %v\n", err)
	}

	// Replace the IPFS path with our local systems filepath.
	metadata.Image = imageFilepath

	//
	// STEP 4
	// Download the animation file from the "network" and save locally.
	//

	animationFilepath, err := s.downloadNFTokAssetUsecase.Execute(sitok.ID, metadata.AnimationURL)
	if err != nil {
		s.logger.Error("Failed downloading nft animation asset.",
			slog.Any("tokenID", sitok.ID),
			slog.Any("AnimationURL", metadata.AnimationURL),
			slog.Any("err", err))
		return fmt.Errorf("Failed downloading nft animation asset: %v\n", err)
	}

	// Replace the IPFS path with our local systems filepath.
	metadata.AnimationURL = animationFilepath

	//
	// STEP 5
	// Update our NFT token to be referenced in future.
	//

	nftok = &domain.NonFungibleToken{
		TokenID:     sitok.ID,
		MetadataURI: metadataURI,
		Metadata:    metadata,
		State:       domain.NonFungibleTokenStateReady,
	}

	if err := s.upsertNFTokUseCase.Execute(nftok); err != nil {
		s.logger.Error("Failed creating nft token.",
			slog.Any("tokenID", sitok.ID))
		return fmt.Errorf("Failed creating nft token: %v\n", err)
	}

	s.logger.Debug("Downloaded NFT successfully.",
		slog.Any("token_id", nftok.TokenID),
		slog.Any("metadata_uri", nftok.MetadataURI),
		slog.Any("metadata", nftok.Metadata),
	)
	return nil
}
