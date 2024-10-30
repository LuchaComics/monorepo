package main

import (
	"fmt"
	"log/slog"
	"strings"
)

func (a *App) GetIsIPFSRunning() bool {
	identity, err := a.ipfsRepo.ID()
	if err != nil {
		a.logger.Error("failed connecting to IPFS repo to get ID()",
			slog.Any("error", err))
		return false
	}
	fmt.Printf("IPFS Node ID: %s\n", identity.ID)

	return true
}

type IPFSFileResponse struct {
	Data          []byte `json:"data"`
	ContentType   string `json:"content_type"`
	ContentLength uint64 `json:"content_length"`
}

func (a *App) GetFileViaIPFS(ipfsPath string) (*IPFSFileResponse, error) {
	cid := strings.Replace(ipfsPath, "ipfs://", "", -1)
	bytes, contentType, contentLength, err := a.ipfsRepo.Cat(cid)
	if err != nil {
		a.logger.Error("failed getting from cid",
			slog.Any("error", err))
		return nil, err
	}

	a.logger.Debug("",
		slog.Any("cid", cid),
		slog.Any("contentType", contentType),
		slog.Any("contentLength", contentLength))

	return &IPFSFileResponse{
		Data:          bytes,
		ContentType:   contentType,
		ContentLength: contentLength,
	}, nil
}
