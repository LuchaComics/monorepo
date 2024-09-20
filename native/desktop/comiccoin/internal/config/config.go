package config

type Config struct {
	BlockchainDifficulty int
	Peer                 PeerConfig
	DB                   DBConfig
}

type PeerConfig struct {
	ListenPort     int
	RandomSeed     int64
	BootstrapPeers string
}

type DBConfig struct {
	DataDir string
}
