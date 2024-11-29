package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

// GetTokens returns a list of all Tokens stored in the repository.
func (a *App) GetTokens() ([]*domain.Token, error) {
	// Retrieve all Tokens from the repository.
	res, err := a.tokenRepo.ListAll()
	if err != nil {
		// If an error occurs, return an empty list and the error.
		return make([]*domain.Token, 0), err
	}
	// If no Tokens are found, return an empty list.
	if res == nil {
		res = make([]*domain.Token, 0)
	}
	return res, nil
}

// GetTokens returns the Token stored in the repository for the particular `tokenID`.
func (a *App) GetToken(tokenID uint64) (*domain.Token, error) {
	tok, err := a.tokenRepo.GetByTokenID(tokenID)
	if err != nil {
		a.logger.Error("Failed getting by token ID.",
			slog.Any("error", err),
		)
		return nil, err
	}

	a.logger.Debug("fetched token",
		slog.Any("token_id", tok.TokenID),
		slog.Any("metadata_uri", tok.MetadataURI),
		slog.Any("image", tok.Metadata.Image),
		slog.Any("animation", tok.Metadata.AnimationURL),
	)
	return tok, nil
}

// GetImageFilePathFromDialog opens a file dialog for the user to select an image file.
func (a *App) GetImageFilePathFromDialog() string {
	// Initialize Wails runtime to interact with the desktop application.
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		// Set the title of the file dialog.
		Title: "Please select the image for this Token",
		// Set the file filters to only show image files.
		Filters: []runtime.FileFilter{
			{
				// Set the display name of the filter.
				DisplayName: "Images (*.png;*.jpg)",
				// Set the file pattern to match.
				Pattern: "*.png;*.jpg",
			},
		},
	})
	if err != nil {
		// If an error occurs, log the error and exit the application.
		a.logger.Error("Failed opening file dialog",
			slog.Any("error", err))
		log.Fatalf("%v", err)
	}
	return result
}

// GetVideoFilePathFromDialog opens a file dialog for the user to select a video file.
func (a *App) GetVideoFilePathFromDialog() string {
	// Initialize Wails runtime to interact with the desktop application.
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		// Set the title of the file dialog.
		Title: "Please select the video for this Token",
		// Set the file filters to only show video files.
		Filters: []runtime.FileFilter{
			{
				// Set the display name of the filter.
				DisplayName: "Videos (*.mov;*.mp4;*.webm)",
				// Set the file pattern to match.
				Pattern: "*.mov;*.mp4;*.webm",
			},
		},
	})
	if err != nil {
		// If an error occurs, log the error and exit the application.
		a.logger.Error("Failed opening file dialog",
			slog.Any("error", err))
		log.Fatalf("%v", err)
	}
	return result
}

func isValidHexColor(color string) bool {
	// Regular expression pattern for HTML hex colors
	pattern := `^#([A-Fa-f0-9]{6}|[A-Fa-f0-9]{3})$`
	match, _ := regexp.MatchString(pattern, color)
	return match
}

// CreateToken creates a new Token with the given metadata and uploads it to IPFS.
func (a *App) CreateToken(
	name string,
	description string,
	image string,
	animation string,
	youtubeURL string,
	externalURL string,
	attributes string,
	backgroundColor string,
) (*domain.Token, error) {
	//
	// STEP 1: Validation.
	//

	// Log the input values for debugging purposes.
	a.logger.Debug("received",
		slog.String("name", name),
		slog.String("image", image),
		slog.String("animation", animation),
		slog.String("youtubeURL", youtubeURL),
		slog.String("externalURL", externalURL),
		slog.Any("attributes", attributes),
		slog.String("backgroundColor", backgroundColor),
	)

	// Check if any of the required fields are missing.
	e := make(map[string]string)
	if name == "" {
		e["name"] = "missing value"
	}
	if description == "" {
		e["description"] = "missing value"
	}
	if image == "" {
		e["image"] = "missing value"
	}
	if animation == "" {
		e["animation"] = "missing value"
	}
	if backgroundColor == "" {
		e["background_color"] = "missing value"
	} else {
		if !isValidHexColor(backgroundColor) {
			e["background_color"] = "wrong formatting, must be HTML hex formatting"
		}
	}
	if len(e) != 0 {
		// If any fields are missing, log an error and return a bad request error.
		a.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Handle `attributes` unmarshalling.
	//

	attrs := make([]*domain.TokenMetadataAttribute, 0)
	if attributes != "" {
		if err := json.Unmarshal([]byte(attributes), &attrs); err != nil {
			// If an error occurs, log an error and return an error.
			a.logger.Error("Failed unmarshal metadata attributes",
				slog.Any("attributes", attributes),
				slog.Any("error", err))
			return nil, fmt.Errorf("failed to deserialize metadata attributete: %v", err)
		}
		a.logger.Debug("attributes",
			slog.Any("attrs", attrs),
			slog.Any("attrs_count", len(attrs)))

		for _, atr := range attrs {
			a.logger.Debug("attribute",
				slog.Any("display_type", atr.DisplayType),
				slog.Any("trait_type", atr.TraitType),
				slog.Any("value", atr.Value))
		}
	}

	//
	// STEP 3: Get related data.
	//

	tokenID, err := a.latestTokenIDRepo.Get()
	if err != nil {
		a.logger.Error("Failed getting latest token ID.",
			slog.Any("error", err))
		return nil, err
	}

	// Please note that in ComicCoin genesis block, we already have a token set
	// at zero. Therefore this increment will work well.
	tokenID++

	// Get our data directory from our app preferences.
	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory

	//
	// STEP 4: Upload `image` file to IPFS.
	//

	imageCID, err := a.nftAssetRepo.PinAddViaFilepath(a.ctx, image)
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed adding to IPFS.",
			slog.Any("filepath", image),
			slog.Any("error", err))
		return nil, err
	}
	if imageCID == "" {
		err := fmt.Errorf("Failed uploading image")
		a.logger.Error("IPFS returned no cid.",
			slog.Any("filepath", image),
			slog.Any("error", err))
		return nil, err
	}
	a.logger.Debug("Image uploaded to ipfs.",
		slog.Any("local", image),
		slog.Any("cid", imageCID))

	//
	// STEP 5: Upload `animation` file to IPFs.
	//

	animationCID, err := a.nftAssetRepo.PinAddViaFilepath(a.ctx, animation)
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed adding animation to IPFs.",
			slog.Any("filepath", animation),
			slog.Any("error", err))
		return nil, err
	}
	a.logger.Debug("Animation uploaded to ipfs.",
		slog.Any("local", animation),
		slog.Any("cid", animationCID))

	//
	// STEP 6:
	// Create token `metadata` file locally.
	//

	metadata := &domain.TokenMetadata{
		Image:           fmt.Sprintf("ipfs://%v", imageCID),
		ExternalURL:     externalURL,
		Description:     description,
		Name:            name,
		Attributes:      attrs,
		BackgroundColor: backgroundColor,
		AnimationURL:    fmt.Sprintf("ipfs://%v", animationCID),
		YoutubeURL:      youtubeURL,
	}

	metadataBytes, err := json.MarshalIndent(metadata, "", "\t")
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed marshal metadata",
			slog.Any("error", err))
		return nil, err
	}

	metadataFilepath := filepath.Join(dataDir, "token_assets", fmt.Sprintf("%v", tokenID), "metadata.json")

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(metadataFilepath), 0755); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed create directories",
			slog.Any("error", err))
		return nil, err
	}

	if err := ioutil.WriteFile(metadataFilepath, metadataBytes, 0644); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed write metadata file",
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 7
	// Copy `image` file so we can consolidate our token assets.
	//

	// Get the base name of the file
	imageFilename := filepath.Base(image)

	// Create new filepath of where to consolidate our file to.
	consolidatedImage := filepath.Join(dataDir, "token_assets", fmt.Sprintf("%v", tokenID), imageFilename)

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(consolidatedImage), 0755); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed create directories",
			slog.Any("error", err))
		return nil, err
	}

	if err := CopyFile(image, consolidatedImage); err != nil {
		a.logger.Error("Failed consolidating image.",
			slog.Any("tokenID", tokenID),
			slog.Any("dataDir", dataDir),
			slog.Any("consolidatedImage", consolidatedImage),
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 8
	// Copy `image` file so we can consolidate our token assets.
	//
	// Get the base name of the file
	animationFilename := filepath.Base(animation)

	// Create new filepath of where to consolidate our file to.
	consolidatedAnimation := filepath.Join(dataDir, "token_assets", fmt.Sprintf("%v", tokenID), animationFilename)

	// Create the directories recursively.
	if err := os.MkdirAll(filepath.Dir(consolidatedAnimation), 0755); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed create directories",
			slog.Any("error", err))
		return nil, err
	}

	if err := CopyFile(animation, consolidatedAnimation); err != nil {
		a.logger.Error("Failed consolidating animation.",
			slog.Any("tokenID", tokenID),
			slog.Any("dataDir", dataDir),
			slog.Any("consolidatedImage", consolidatedImage),
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 9:
	// Upload to IPFs and get the CID.
	//

	metadataCID, err := a.nftAssetRepo.PinAddViaFilepath(a.ctx, metadataFilepath)
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed adding metadata to IPFs.",
			slog.Any("token_id", tokenID),
			slog.Any("filepath", metadataFilepath),
			slog.Any("error", err))
		return nil, err
	}
	a.logger.Debug("Metadata uploaded to ipfs.",
		slog.Any("token_id", tokenID),
		slog.Any("local", metadataFilepath),
		slog.Any("cid", metadataCID))

	//
	// STEP 10:
	// Update our database.
	//

	token := &domain.Token{
		TokenID:     tokenID,
		MetadataURI: fmt.Sprintf("ipfs://%v", metadataCID),
		Metadata:    metadata,
		Timestamp:   uint64(time.Now().UTC().UnixMilli()),
	}

	if err := a.tokenRepo.Upsert(token); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed save to database the new token.",
			slog.Any("error", err))
		return nil, err
	}

	if err := a.latestTokenIDRepo.Set(tokenID); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed save to database the latest token ID.",
			slog.Any("error", err))
		return nil, err
	}

	return token, nil
}
