package config

type Config struct {
	DB DBConfig
}

type DBConfig struct {
	DataDir string
}
