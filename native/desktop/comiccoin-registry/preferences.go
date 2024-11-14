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

func init() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed get home dir: %v\n", err)
	}
	FilePathPreferences = filepath.Join(homeDir, ".comiccoin.registry")
}

type Preferences struct {
	// DataDirectory variable holds the location of were the entire application
	// will be saved on the user's computer.
	DataDirectory string `json:"data_directory"`

	// DefaultWalletAddress holds the address of the wallet that will be
	// automatically opend every time the application loads up. This is selected
	// by the user.
	DefaultWalletAddress string `json:"default_wallet_address"`

	// NFTStoreSettings holds the configuration for connecting to
	// our remote IPFS node which handles IPFS submission.
	NFTStoreSettings map[string]string `json:"nft_store_settings"`
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
			log.Fatalf("failed open file: %v\n", err)
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

func (pref *Preferences) GetDefaultDataDirectory() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("failed get home dir: %v\n", err)
	}
	return filepath.Join(homeDir, "ComicCoinRegistry")
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

func (pref *Preferences) GetEmptyNFTStoreSettings() map[string]string {
	remoteIPFSNodeConfiguration := map[string]string{}
	remoteIPFSNodeConfiguration["apiVersion"] = ""
	remoteIPFSNodeConfiguration["accessKeyId"] = ""
	remoteIPFSNodeConfiguration["secretAccessKey"] = ""
	remoteIPFSNodeConfiguration["endpoint"] = ""
	remoteIPFSNodeConfiguration["region"] = ""
	remoteIPFSNodeConfiguration["s3ForcePathStyle"] = "true"
	return remoteIPFSNodeConfiguration
}

func (pref *Preferences) SetNFTStoreSettings(config map[string]string) error {
	pref.NFTStoreSettings = config
	data, err := json.MarshalIndent(pref, "", "\t")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(FilePathPreferences, data, 0666)
}

//
// func (pref *Preferences) SetNFTStoreAPIKey(apiKey string) error {
// 	pref.NFTStoreAPIKey = apiKey
// 	data, err := json.MarshalIndent(pref, "", "\t")
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(FilePathPreferences, data, 0666)
// }
//
// func (pref *Preferences) SetNFTStoreRemoteAddress(remoteAddress string) error {
// 	pref.NFTStoreRemoteAddress = remoteAddress
// 	data, err := json.MarshalIndent(pref, "", "\t")
// 	if err != nil {
// 		return err
// 	}
// 	return ioutil.WriteFile(FilePathPreferences, data, 0666)
// }
