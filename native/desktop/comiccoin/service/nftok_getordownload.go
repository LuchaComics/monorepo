package service

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/usecase"
)

type GetOrDownloadNonFungibleTokenService struct {
	config                       *config.Config
	logger                       *slog.Logger
	getNFTokUseCase              *usecase.GetNonFungibleTokenUseCase
	getTokUseCase                *usecase.GetTokenUseCase
	downloadNFTokMetadataUsecase *usecase.DownloadMetadataNonFungibleTokenUseCase
	downloadNFTokAssetUsecase    *usecase.DownloadNonFungibleTokenAssetUseCase
	upsertNFTokUseCase           *usecase.UpsertNonFungibleTokenUseCase
}

func NewGetOrDownloadNonFungibleTokenService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetNonFungibleTokenUseCase,
	uc2 *usecase.GetTokenUseCase,
	uc3 *usecase.DownloadMetadataNonFungibleTokenUseCase,
	uc4 *usecase.DownloadNonFungibleTokenAssetUseCase,
	uc5 *usecase.UpsertNonFungibleTokenUseCase,
) *GetOrDownloadNonFungibleTokenService {
	return &GetOrDownloadNonFungibleTokenService{cfg, logger, uc1, uc2, uc3, uc4, uc5}
}

func (s *GetOrDownloadNonFungibleTokenService) Execute(tokenID uint64) (*domain.NonFungibleToken, error) {
	//
	// STEP 1
	// Lookup our `token id` in our NFT db and if it exists we can return value.
	//

	nftok, err := s.getNFTokUseCase.Execute(tokenID)
	if err != nil {
		s.logger.Error("Failed getting non-fungible token.",
			slog.Any("tokenID", tokenID),
			slog.Any("error", err))
		return nil, err
	}
	if nftok != nil {
		return nftok, nil
	}

	//
	// STEP 2
	// Lookup our `token` in our db and retrieve the record so we can
	// extract the `Metadata URI` value necessary to lookup later in
	// the decentralized storage service (IPFS).
	//

	tok, err := s.getTokUseCase.Execute(tokenID)
	if err != nil {
		s.logger.Error("Failed getting token",
			slog.Any("tokenID", tokenID),
			slog.Any("error", err))
		return nil, fmt.Errorf("Failed getting token due to err: %v\n", err)
	}
	if tok == nil {
		s.logger.Error("Token does not exist.",
			slog.Any("tokenID", tokenID))
		return nil, fmt.Errorf("Token does not exist for: %v\n", tokenID)
	}

	metadataURI := tok.MetadataURI
	// Confirm URI is using protocol our app supports.
	if strings.Contains(metadataURI, "ipfs://") == false {
		s.logger.Error("Token metadata URI contains protocol we do not support:",
			slog.Any("tokenID", tokenID),
			slog.Any("metadataURI", metadataURI))
		return nil, fmt.Errorf("Token metadata URI contains protocol we do not support: %v\n", metadataURI)
	}

	metadata, metadataFilepath, err := s.downloadNFTokMetadataUsecase.Execute(tok.ID, metadataURI)
	if err != nil {
		s.logger.Error("Failed downloading nft metadata.",
			slog.Any("tokenID", tokenID))
		return nil, fmt.Errorf("Failed downloading nft metadata: %v\n", err)
	}

	// Replace the IPFS path with our local systems filepath.
	metadataURI = metadataFilepath

	//
	// STEP 3
	// Create our NFT token to be referenced in future.
	//

	nftok = &domain.NonFungibleToken{
		TokenID:     tokenID,
		MetadataURI: metadataURI,
		Metadata:    metadata,
		State:       domain.NonFungibleTokenStateNotReady,
	}

	if err := s.upsertNFTokUseCase.Execute(nftok); err != nil {
		s.logger.Error("Failed creating nft token.",
			slog.Any("tokenID", tokenID))
		return nil, fmt.Errorf("Failed creating nft token: %v\n", err)
	}

	//
	// STEP 4
	// Download the image file from IPFS and save locally.
	//

	imageCID := strings.Replace(metadata.Image, "ipfs://", "", -1)
	imageFilepath, err := s.downloadNFTokAssetUsecase.Execute(tok.ID, imageCID)
	if err != nil {
		s.logger.Error("Failed downloading nft image asset.",
			slog.Any("tokenID", tokenID))
		return nil, fmt.Errorf("Failed downloading nft image asset: %v\n", err)
	}

	// Replace the IPFS path with our local systems filepath.
	metadata.Image = imageFilepath

	//
	// STEP 5
	// Download the animation file from IPFS and save locally.
	//

	animationCID := strings.Replace(metadata.AnimationURL, "ipfs://", "", -1)
	animationFilepath, err := s.downloadNFTokAssetUsecase.Execute(tok.ID, animationCID)
	if err != nil {
		s.logger.Error("Failed downloading nft image asset.",
			slog.Any("tokenID", tokenID))
		return nil, fmt.Errorf("Failed downloading nft image asset: %v\n", err)
	}

	// Replace the IPFS path with our local systems filepath.
	metadata.AnimationURL = animationFilepath

	//
	// STEP 6
	// Update our NFT token to be referenced in future.
	//

	nftok = &domain.NonFungibleToken{
		TokenID:     tokenID,
		MetadataURI: metadataURI,
		Metadata:    metadata,
		State:       domain.NonFungibleTokenStateReady,
	}

	if err := s.upsertNFTokUseCase.Execute(nftok); err != nil {
		s.logger.Error("Failed creating nft token.",
			slog.Any("tokenID", tokenID))
		return nil, fmt.Errorf("Failed creating nft token: %v\n", err)
	}

	s.logger.Debug("Downloaded NFT successfully.",
		slog.Any("token_id", nftok.TokenID),
		slog.Any("metadata_uri", nftok.MetadataURI),
		slog.Any("metadata", nftok.Metadata),
	)
	return nftok, nil
}
