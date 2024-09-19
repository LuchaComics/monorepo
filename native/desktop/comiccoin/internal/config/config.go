package config

type Config struct {
	BlockchainDifficulty int
	AppPort              int
	DB                   DBConfig
}

type DBConfig struct {
	DataDir string
}
