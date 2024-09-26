package config

import (
	maddr "github.com/multiformats/go-multiaddr"
)

type Config struct {
	//
	Blockchain BlockchainConfig

	//
	App AppConfig

	//
	Peer PeerConfig

	//
	DB DBConfig
}

type BlockchainConfig struct {
	ChainID       uint16 `json:"chain_id"`        // The chain id represents an unique id for this running instance.
	TransPerBlock uint16 `json:"trans_per_block"` // The maximum number of transactions that can be in a block.
	Difficulty    uint16 `json:"difficulty"`      // How difficult it needs to be to solve the work problem.
	MiningReward  uint64 `json:"mining_reward"`   // Reward for mining a block.
	GasPrice      uint64 `json:"gas_price"`       // Fee paid for each transaction mined into a block.
	UnitsOfGas    uint64 `json:"units_of_gas"`
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
	DataDir string
}
