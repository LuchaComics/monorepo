package config

type Config struct {
	BlockchainDifficulty int
	Peer                 PeerConfig
	DB                   DBConfig
}

type PeerConfig struct {
	ListenPort     int
	KeyName        string
	BootstrapPeers string
}

type DBConfig struct {
	DataDir string
}
