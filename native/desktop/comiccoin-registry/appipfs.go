package main

import (
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

func (a *App) GetIsNFTAssetStoreRunning() bool {
	version, err := a.nftAssetRepo.Version(a.ctx)
	if err != nil {
		a.logger.Error("failed connecting to NFT Asset store repo",
			slog.Any("error", err))
		return false
	}
	return version == "1.0"
}

func (a *App) GetFileViaIPFS(ipfsPath string) (*domain.NFTAsset, error) {
	cid := strings.Replace(ipfsPath, "ipfs://", "", -1)
	resp, err := a.nftAssetRepo.Get(a.ctx, cid)
	if err != nil {
		a.logger.Error("failed getting from cid",
			slog.Any("error", err))
		return nil, err
	}

	a.logger.Debug("GetFileViaIPFS",
		slog.Any("cid", cid),
		slog.Any("content_type", resp.ContentType))

	return resp, nil
}
