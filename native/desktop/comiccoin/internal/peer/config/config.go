package config

import (
	maddr "github.com/multiformats/go-multiaddr"
)

type Config struct {
	App  AppConfig
	Peer PeerConfig
	DB   DBConfig
}

type AppConfig struct {
	// DirPath variable is the path to the directory where all the files for
	// this appliction to
	// save to.
	DirPath string

	// HttpPort variable is the port that the HTTP JSON API server will listen
	// on for this application. Do not expose to public!
	HTTPPort int

	// HttpIP variable is the address to bind the HTTP JSON API server onto.
	HTTPIP string
}

type PeerConfig struct {
	ListenPort       int
	KeyName          string
	RendezvousString string
	BootstrapPeers   []maddr.Multiaddr
	ListenAddresses  []maddr.Multiaddr
}

type DBConfig struct {
	// Location of were to save the database files.
	DataDir string
}
