package main

import (
	"encoding/json"
	"fmt"
	"log"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
	"github.com/wailsapp/wails/v2/pkg/runtime"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

func (a *App) GetNFTs() ([]*domain.NFT, error) {
	res, err := a.nftRepo.ListAll()
	if err != nil {
		return make([]*domain.NFT, 0), err
	}
	if res == nil {
		res = make([]*domain.NFT, 0)
	}
	return res, nil
}

func (a *App) GetImageFilePathFromDialog() string {
	// Initialize Wails runtime
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Please select the image for this NFT",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Images (*.png;*.jpg)",
				Pattern:     "*.png;*.jpg",
			},
		},
	})
	if err != nil {
		a.logger.Error("Failed opening file dialog",
			slog.Any("error", err))
		log.Fatalf("%v", err)
	}
	return result
}

func (a *App) GetVideoFilePathFromDialog() string {
	// Initialize Wails runtime
	result, err := runtime.OpenFileDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Please select the image for this NFT",
		Filters: []runtime.FileFilter{
			{
				DisplayName: "Videos (*.mov;*.mp4;*.webm)",
				Pattern:     "*.mov;*.mp4;*.webm",
			},
		},
	})
	if err != nil {
		a.logger.Error("Failed opening file dialog",
			slog.Any("error", err))
		log.Fatalf("%v", err)
	}
	return result
}

func (a *App) CreateNFT(
	name string,
	description string,
	image string,
	animation string,
	youtubeURL string,
	externalURL string,
	attributes string,
	backgroundColor string,
) (*domain.NFT, error) {
	// For debugging purposes only.
	a.logger.Debug("received",
		slog.String("name", name),
		slog.String("image", image),
		slog.String("animation", animation),
		slog.String("youtubeURL", youtubeURL),
		slog.String("externalURL", externalURL),
		slog.Any("attributes", attributes),
		slog.String("backgroundColor", backgroundColor),
	)

	// Defensive code purposes.
	e := make(map[string]string)
	if name == "" {
		e["name"] = "missing value"
	}
	if len(e) != 0 {
		a.logger.Warn("Failed creating new account",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	if attributes == nil {
		// attr
		if err := json.Unmarshal(attributes, &wallet); err != nil {
			// Return an error if the unmarshaling fails.
			return nil, fmt.Errorf("failed to deserialize wallet: %v", err)
		}
	}

	nft := &domain.NFT{}
	return nft, nil
}
