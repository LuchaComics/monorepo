package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
)

type Preferences struct {
	// DataDirectory variable holds the location of were the entire application
	// will be saved on the user's computer.
	DataDirectory string `json:"data_directory"`

	// DefaultWalletAddress holds the address of the wallet that will be
	// automatically opend every time the application loads up. This is selected
	// by the user.
	DefaultWalletAddress string `json:"default_wallet_address"`

	// NFTStorageAddress variable holds the full address to the location
	// of the NFTStore on the network. Example: https://example.com or
	// http://127.0.0.1:8080.
	NFTStorageAddress string `json:"nft_storage_address"`

	// ChainID variable keeps track which blockchain version we are using.
	ChainID uint16 `json:"chain_id"`

	// AuthorityAddress holds the address of the ComicCoin blockchain authority
	// address that our client will communicate with.
	AuthorityAddress string `json:"authority_address"`
}

var (
	instance            *Preferences
	once                sync.Once
	FilePathPreferences string
)

func PreferencesInstance() *Preferences {
	once.Do(func() {
		// Either reads the file if the file exists or creates an empty.
		file, err := os.OpenFile(FilePathPreferences, os.O_RDONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Fatalf("failed open file %v because of error: %v\n", FilePathPreferences, err)
		}

		var preferences Preferences
		err = json.NewDecoder(file).Decode(&preferences)
		file.Close() // Close the file after you're done with it
		if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
			log.Fatalf("failed decode file: %v\n", err)
		}
		instance = &preferences
	})
	return instance
}

func GetDefaultDataDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed get home dir: %v\n", err)
	}
	return filepath.Join(homeDir, "ComicCoin")
}

// AbortOnValidationFailure method will check the preferences and if any field
// was not set then trigger a failure.
func (pref *Preferences) RunFatalIfHasAnyMissingFields() {
	if pref.DataDirectory == "" {
		log.Fatal("Missing configuration for ComicCoin: `DataDirectory` was not set. Please run in your console: `./comiccoin init`\n")
	}
	if pref.ChainID == 0 {
		log.Fatal("You have already configured ComicCoin: `ChainID` was set. Please run in your console: `./comiccoin init`\n")
	}
	if pref.AuthorityAddress == "" {
		log.Fatal("You have already configured ComicCoin: `AuthorityAddress` was set. Please run in your console: `./comiccoin init`\n")
	}
	if pref.NFTStorageAddress == "" {
		log.Fatal("You have already configured ComicCoin: `NFTStorageAddress` was set. Please run in your console: `./comiccoin init`\n")
	}
}

func (pref *Preferences) SetDataDirectory(dataDir string) error {
	pref.DataDirectory = dataDir
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetDefaultWalletAddress(newAdrs string) error {
	pref.DefaultWalletAddress = newAdrs
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetChainID(chainID uint16) error {
	pref.ChainID = chainID
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetAuthorityAddress(authorityAddress string) error {
	pref.AuthorityAddress = authorityAddress
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) SetNFTStorageAddress(remoteAddress string) error {
	pref.NFTStorageAddress = remoteAddress
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(FilePathPreferences, data, 0666)
}

func (pref *Preferences) GetFilePathOfPreferencesFile() string {
	return FilePathPreferences
}
