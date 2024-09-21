package config

import (
	maddr "github.com/multiformats/go-multiaddr"
)

type Config struct {
	BlockchainDifficulty int
	Peer                 PeerConfig
	DB                   DBConfig
}

type PeerConfig struct {
	ListenPort       int
	KeyName          string
	RendezvousString string
	BootstrapPeers   []maddr.Multiaddr
	ListenAddresses  []maddr.Multiaddr
}

type DBConfig struct {
	DataDir string
}
