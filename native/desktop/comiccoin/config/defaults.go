package config

import (
	"log"
	"os"
	"path/filepath"
)

// GetDefaultDataDirectory function returns the *recommended* location of were
// to save all our blockchain data for this application.
func GetDefaultDataDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed get home dir: %v\n", err)
	}
	return filepath.Join(homeDir, "ComicCoin")
}

func GetDefaultPeerToPeerListenPort() int {
	return 26642
}
