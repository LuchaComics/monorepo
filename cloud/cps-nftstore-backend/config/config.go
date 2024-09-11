package config

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	AppServer          serverConf
	DB                 dbConfig
	Emailer            mailgunConfig
	EthereumBlockchain ethConfig
}

type serverConf struct {
	Port                    string
	IP                      string
	HMACSecret              []byte
	HasDebugging            bool
	InitialAdminEmail       string
	InitialAdminPassword    string
	InitialAdminTenantName  string
	APIDomainName           string
	AppDomainName           string
	Enable2FAOnRegistration bool
}

type dbConfig struct {
	URI  string
	Name string
}

type mailgunConfig struct {
	APIKey           string
	Domain           string
	APIBase          string
	SenderEmail      string
	MaintenanceEmail string
}

type ethConfig struct {
	NodeURL              string
	SmartContractAddress string
	OwnerAddress         string
	OwnerPrivateKey      string
}

func New() *Conf {
	var c Conf
	c.AppServer.Port = getEnv("CPS_NFTSTORE_BACKEND_PORT", true)
	c.AppServer.IP = getEnv("CPS_NFTSTORE_BACKEND_IP", false)
	c.AppServer.HMACSecret = []byte(getEnv("CPS_NFTSTORE_BACKEND_HMAC_SECRET", true))
	c.AppServer.HasDebugging = getEnvBool("CPS_NFTSTORE_BACKEND_HAS_DEBUGGING", true, true)
	c.AppServer.InitialAdminEmail = getEnv("CPS_NFTSTORE_BACKEND_INITIAL_ADMIN_EMAIL", true)
	c.AppServer.InitialAdminPassword = getEnv("CPS_NFTSTORE_BACKEND_INITIAL_ADMIN_PASSWORD", true)
	c.AppServer.InitialAdminTenantName = getEnv("CPS_NFTSTORE_BACKEND_INITIAL_ADMIN_ORG_NAME", true)
	c.AppServer.APIDomainName = getEnv("CPS_NFTSTORE_BACKEND_API_DOMAIN_NAME", true)
	c.AppServer.AppDomainName = getEnv("CPS_NFTSTORE_BACKEND_APP_DOMAIN_NAME", true)
	c.AppServer.Enable2FAOnRegistration = getEnvBool("CPS_NFTSTORE_BACKEND_APP_ENABLE_2FA_ON_REGISTRATION", false, false)

	c.DB.URI = getEnv("CPS_NFTSTORE_BACKEND_DB_URI", true)
	c.DB.Name = getEnv("CPS_NFTSTORE_BACKEND_DB_NAME", true)

	c.Emailer.APIKey = getEnv("CPS_NFTSTORE_BACKEND_MAILGUN_API_KEY", true)
	c.Emailer.Domain = getEnv("CPS_NFTSTORE_BACKEND_MAILGUN_DOMAIN", true)
	c.Emailer.APIBase = getEnv("CPS_NFTSTORE_BACKEND_MAILGUN_API_BASE", true)
	c.Emailer.SenderEmail = getEnv("CPS_NFTSTORE_BACKEND_MAILGUN_SENDER_EMAIL", true)
	c.Emailer.MaintenanceEmail = getEnv("CPS_NFTSTORE_BACKEND_MAILGUN_MAINTENANCE_EMAIL", true)

	c.EthereumBlockchain.NodeURL = getEnv("CPS_NFTSTORE_BACKEND_ETH_NODE_URL", false)
	c.EthereumBlockchain.SmartContractAddress = getEnv("CPS_NFTSTORE_BACKEND_ETH_SMART_CONTRACT_ADDRESS", false)
	c.EthereumBlockchain.OwnerAddress = getEnv("CPS_NFTSTORE_BACKEND_ETH_OWNER_ADDRESS", false)
	c.EthereumBlockchain.OwnerPrivateKey = getEnv("CPS_NFTSTORE_BACKEND_ETH_PRIVATE_KEY_ADDRESS", false)

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
