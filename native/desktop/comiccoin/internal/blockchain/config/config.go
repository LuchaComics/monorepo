package config

import (
	"github.com/ethereum/go-ethereum/common"
	maddr "github.com/multiformats/go-multiaddr"
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
}

// BlockchainConfig represents the configuration for the blockchain.
// It contains settings for the chain ID, transactions per block, difficulty, mining reward, gas price, and units of gas.
type BlockchainConfig struct {
	// ChainID is the unique ID for this blockchain instance.
	ChainID uint16 `json:"chain_id"`

	// TransPerBlock is the maximum number of transactions that can be included in a block.
	TransPerBlock uint16 `json:"trans_per_block"`

	// Difficulty represents how difficult it should be to solve the work problem.
	Difficulty uint16 `json:"difficulty"`

	// MiningReward is the reward for mining a block.
	MiningReward uint64 `json:"mining_reward"`

	// GasPrice is the fee paid for each transaction included in a block.
	GasPrice uint64 `json:"gas_price"`

	// UnitsOfGas represents the units of gas for each transaction.
	UnitsOfGas uint64 `json:"units_of_gas"`

	// The delay time (in minutes) before this mode will poll the blockchain
	// network to request a consensus as to what the latest block is.
	ConsensusPollingDelayInMinutes int64 `json:"consensus_polling_delay_in_minutes"`

	// Control whether to have the miner running in the background for this node.
	EnableMiner bool `json:"enable_miner"`

	// Used to set what protocol to for mining and coordinating latest blockchain.
	ConsensusProtocol string `json:"consensus_protocol"`

	// (Only set by PoA node)
	ProofOfAuthorityAccountAddress *common.Address

	// (Only set by PoA node)
	ProofOfAuthorityWalletPassword string
}

// AppConfig represents the configuration for the application.
// It contains settings for the directory path, HTTP address, and RPC address.
type AppConfig struct {
	// DirPath is the path to the directory where all files for this application are saved.
	DirPath string

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
