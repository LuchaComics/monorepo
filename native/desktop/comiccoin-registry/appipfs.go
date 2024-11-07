package main

import (
	"log"
	"log/slog"
	"strings"
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

type IPFSFileResponse struct {
	Data          []byte `json:"data"`
	ContentType   string `json:"content_type"`
	ContentLength uint64 `json:"content_length"`
}

func (a *App) GetFileViaIPFS(ipfsPath string) (*IPFSFileResponse, error) {
	cid := strings.Replace(ipfsPath, "ipfs://", "", -1)
	content, contentType, err := a.remoteIpfsRepo.Get(a.ctx, cid)
	if err != nil {
		a.logger.Error("failed getting from cid",
			slog.Any("error", err))
		return nil, err
	}

	a.logger.Debug("GetFileViaIPFS",
		slog.Any("cid", cid),
		slog.Any("content_type", contentType))

	return &IPFSFileResponse{
		Data:        content,
		ContentType: contentType,
	}, nil
}
