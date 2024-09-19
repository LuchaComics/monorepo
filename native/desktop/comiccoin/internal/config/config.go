package config

type Config struct {
	BlockchainDifficulty int
	DB                   DBConfig
}

type DBConfig struct {
	DataDir string
}
