package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
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
	}
	if len(e) != 0 {
		// If any fields are missing, log an error and return a bad request error.
		a.logger.Warn("Failed validating",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3: Lookup and get the latest TokenID and increment by 1.
	//

	tokenID := uint64(1) // TEMPORARY

	//
	// STEP 2: Image upload to IPFS.
	//

	imageUploadResponse, err := a.ipfsRepo.AddAndPinSingleFileFromLocalFileSystem(image)
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed adding to IPFS.",
			slog.Any("filepath", image),
			slog.Any("error", err))
		return nil, err
	}
	a.logger.Debug("Image uploaded to ipfs.",
		slog.Any("local", image),
		slog.Any("cid", imageUploadResponse.Hash))

	//
	// STEP 3: Animation upload to IPFs.
	//

	animationUploadResponse, err := a.ipfsRepo.AddAndPinSingleFileFromLocalFileSystem(animation)
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed adding animation to IPFs.",
			slog.Any("filepath", animation),
			slog.Any("error", err))
		return nil, err
	}
	a.logger.Debug("Animation uploaded to ipfs.",
		slog.Any("local", animation),
		slog.Any("cid", animationUploadResponse.Hash))

	//
	// STEP 4: Attributes.
	//

	attr := make([]*domain.TokenMetadataAttribute, 0)
	if attributes != "" {
		if err := json.Unmarshal([]byte(attributes), &attr); err != nil {
			// If an error occurs, log an error and return an error.
			a.logger.Error("Failed unmarshal metadata attributes",
				slog.Any("attributes", attributes),
				slog.Any("error", err))
			return nil, fmt.Errorf("failed to deserialize metadata attributete: %v", err)
		}
		a.logger.Debug("attributes",
			slog.Any("attr", attr))
	}

	//
	// STEP 5:
	// Create Token metadata file locally.
	//

	metadata := &domain.TokenMetadata{
		Image:           imageUploadResponse.Hash,
		ExternalURL:     externalURL,
		Description:     description,
		Name:            name,
		Attributes:      attr,
		BackgroundColor: backgroundColor,
		AnimationURL:    animationUploadResponse.Hash,
		YoutubeURL:      youtubeURL,
	}

	metadataBytes, err := json.MarshalIndent(metadata, "", "\t")
	if err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed marshal metadata",
			slog.Any("error", err))
		return nil, err
	}

	preferences := PreferencesInstance()
	dataDir := preferences.DataDirectory
	metadataFilepath := filepath.Join(dataDir, "tokens", fmt.Sprintf("%v", tokenID), "metadata.json")

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
	// STEP 6:
	// Upload to IPFs and get the CID.
	//

	metadataUploadResponse, err := a.ipfsRepo.AddAndPinSingleFileFromLocalFileSystem(metadataFilepath)
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
		slog.Any("cid", metadataUploadResponse.Hash))

	//
	// STEP 7:
	// Create Token in our database.
	//

	token := &domain.Token{
		TokenID:     tokenID,
		MetadataURI: fmt.Sprintf("ipfs://%v", metadataUploadResponse.Hash),
		Metadata:    metadata,
	}

	if err := a.tokenRepo.Upsert(token); err != nil {
		// If an error occurs, log an error and return an error.
		a.logger.Error("Failed save to database the token",
			slog.Any("error", err))
		return nil, err
	}
	return token, nil
}
