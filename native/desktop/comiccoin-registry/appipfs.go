package main

import (
	"fmt"
	"log/slog"
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
