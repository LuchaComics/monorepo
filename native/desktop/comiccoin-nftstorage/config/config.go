package config

import (
	maddr "github.com/multiformats/go-multiaddr"

	sbytes "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/securebytes"
	sstring "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/securestring"
)

// Config represents the configuration for the application.
// It contains settings for the blockchain, application, database, and peer connections.
type Config struct {
	// Blockchain configuration.
	Blockchain BlockchainConfig

	// Application configuration.
	App AppConfig

	// Database configuration.
	DB DBConfig

	// Peer configuration.
	Peer PeerConfig

	// IPFS configuration.
	IPFS IPFSConfig
}

// BlockchainConfig represents the configuration for the blockchain.
// It contains settings for the chain ID, transactions per block, difficulty, mining reward, gas price, and units of gas.
type BlockchainConfig struct {
	// ChainID is the unique ID for this blockchain instance.
	ChainID uint16 `json:"chain_id"`
}

// AppConfig represents the configuration for the application.
// It contains settings for the directory path, HTTP address, and RPC address.
type AppConfig struct {
	// DirPath is the path to the directory where all files for this application are saved.
	DirPath string

	HMACSecret *sbytes.SecureBytes
	AppSecret  *sstring.SecureString

	// HTTPAddress is the address and port that the HTTP JSON API server will listen on.
	// Do not expose to the public!
	HTTPAddress string
}

// DBConfig represents the configuration for the database.
// It contains the location of the database files.
type DBConfig struct {
	// DataDir is the location of the database files.
	DataDir string
}

// PeerConfig represents the configuration for peer connections.
// It contains settings for the listen port, key name, rendezvous string, bootstrap peers, and listen addresses.
type PeerConfig struct {
	// ListenPort is the port that the peer will listen on.
	ListenPort int

	// KeyName is the name of the key used for encryption.
	KeyName string

	// RendezvousString is the string used for rendezvous connections.
	RendezvousString string

	// BootstrapPeers is a list of multiaddresses for bootstrap peers.
	BootstrapPeers []maddr.Multiaddr

	// ListenAddresses is a list of multiaddresses that the peer will listen on.
	ListenAddresses []maddr.Multiaddr
}

type IPFSConfig struct {
	// RemoteIP is the IP address of the peer node running IPFS.
	RemoteIP string

	// RemotePort is the port to the peer node running IPFS.
	RemotePort string

	// PublicGatewayDomain is the HTTP domain to use as a fall-back if peer node is not running.
	PublicGatewayDomain string
}
