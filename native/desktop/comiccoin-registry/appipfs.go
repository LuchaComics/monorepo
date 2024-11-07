package main

import (
	"log"
	"log/slog"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

func (a *App) GetIsIPFSRunning() bool {
	// peerID, err := a.remoteIpfsRepo.ID()
	// if err != nil {
	// 	a.logger.Error("failed connecting to IPFS repo to get ID()",
	// 		slog.Any("error", err))
	// 	return false
	// }
	// fmt.Printf("IPFS Node ID: %s\n", peerID)
	log.Fatal("deprecated")
	return true
}

func (a *App) GetFileViaIPFS(ipfsPath string) (*domain.RemoteIPFSGetFileResponse, error) {
	cid := strings.Replace(ipfsPath, "ipfs://", "", -1)
	resp, err := a.remoteIpfsRepo.Get(a.ctx, cid)
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
