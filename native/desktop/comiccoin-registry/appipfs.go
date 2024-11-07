package main

import (
	"log"
	"log/slog"
	"strings"

	pkgdomain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

func (a *App) GetIsIPFSRunning() bool {
	// peerID, err := a.nftAssetRepo.ID()
	// if err != nil {
	// 	a.logger.Error("failed connecting to IPFS repo to get ID()",
	// 		slog.Any("error", err))
	// 	return false
	// }
	// fmt.Printf("IPFS Node ID: %s\n", peerID)
	log.Fatal("deprecated")
	return true
}

func (a *App) GetFileViaIPFS(ipfsPath string) (*pkgdomain.NFTAsset, error) {
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
