package config

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	EthServer ethConf
}

type ethConf struct {
	NodeURL         string
	OwnerAddress    string
	OwnerPrivateKey string
}

func New() *Conf {
	var c Conf
	c.EthServer.NodeURL = getEnv("CPS_NFTSTORE_CLI_ETH_NODE_URL", true)
	c.EthServer.OwnerAddress = getEnv("CPS_NFTSTORE_CLI_ETH_OWNER_ADDRESS", true)
	c.EthServer.OwnerPrivateKey = getEnv("CPS_NFTSTORE_CLI_ETH_OWNER_PRIVATE_KEY", true)

	return &c
}

func getEnv(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return value
}

func getEnvBool(key string, required bool, defaultValue bool) bool {
	valueStr := getEnv(key, required)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatalf("Invalid boolean value for environment variable %s", key)
	}
	return value
}
