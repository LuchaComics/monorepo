package config

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	AppServer        serverConf
	DB               dbConfig
	PDFBuilder       pdfBuilderConfig
	Emailer          mailgunConfig
	PaymentProcessor paymentProcessorConfig
	IPFSNode         ipfsConfig
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
	IsDeveloperMode         bool
	Enable2FAOnRegistration bool
}

type dbConfig struct {
	URI  string
	Name string
}

type pdfBuilderConfig struct {
	CBFFTemplatePath  string
	PCTemplatePath    string
	CCIMGTemplatePath string
	CCSCTemplatePath  string
	CCTemplatePath    string
	CCUGTemplatePath  string
	DataDirectoryPath string
}

type mailgunConfig struct {
	APIKey           string
	Domain           string
	APIBase          string
	SenderEmail      string
	MaintenanceEmail string
}

type paymentProcessorConfig struct {
	SecretKey        string
	PublicKey        string
	WebhookSecretKey string
}

type ipfsConfig struct {
	GatewayRPCURL string
	Username      string
	Password      string
}

func New() *Conf {
	var c Conf
	c.AppServer.IsDeveloperMode = getEnvBool("CPS_IPFSSTORE_BACKEND_APP_IS_DEVELOPER_MODE", false, true) // If in doubt assume developer mode!
	c.AppServer.Port = getEnv("CPS_IPFSSTORE_BACKEND_PORT", true)
	c.AppServer.IP = getEnv("CPS_IPFSSTORE_BACKEND_IP", false)
	c.AppServer.HMACSecret = []byte(getEnv("CPS_IPFSSTORE_BACKEND_HMAC_SECRET", true))
	c.AppServer.HasDebugging = getEnvBool("CPS_IPFSSTORE_BACKEND_HAS_DEBUGGING", true, true)
	c.AppServer.InitialAdminEmail = getEnv("CPS_IPFSSTORE_BACKEND_INITIAL_ADMIN_EMAIL", true)
	c.AppServer.InitialAdminPassword = getEnv("CPS_IPFSSTORE_BACKEND_INITIAL_ADMIN_PASSWORD", true)
	c.AppServer.InitialAdminTenantName = getEnv("CPS_IPFSSTORE_BACKEND_INITIAL_ADMIN_ORG_NAME", true)
	c.AppServer.APIDomainName = getEnv("CPS_IPFSSTORE_BACKEND_API_DOMAIN_NAME", true)
	c.AppServer.AppDomainName = getEnv("CPS_IPFSSTORE_BACKEND_APP_DOMAIN_NAME", true)
	c.AppServer.Enable2FAOnRegistration = getEnvBool("CPS_IPFSSTORE_BACKEND_APP_ENABLE_2FA_ON_REGISTRATION", false, false)

	c.DB.URI = getEnv("CPS_IPFSSTORE_BACKEND_DB_URI", true)
	c.DB.Name = getEnv("CPS_IPFSSTORE_BACKEND_DB_NAME", true)

	c.PDFBuilder.CBFFTemplatePath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_CBFF_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.PCTemplatePath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_PC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCIMGTemplatePath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_CCIMG_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCSCTemplatePath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_CCSC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCTemplatePath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_CC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCUGTemplatePath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_CCUG_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.DataDirectoryPath = getEnv("CPS_IPFSSTORE_BACKEND_PDF_BUILDER_DATA_DIRECTORY_PATH", true)

	c.Emailer.APIKey = getEnv("CPS_IPFSSTORE_BACKEND_MAILGUN_API_KEY", true)
	c.Emailer.Domain = getEnv("CPS_IPFSSTORE_BACKEND_MAILGUN_DOMAIN", true)
	c.Emailer.APIBase = getEnv("CPS_IPFSSTORE_BACKEND_MAILGUN_API_BASE", true)
	c.Emailer.SenderEmail = getEnv("CPS_IPFSSTORE_BACKEND_MAILGUN_SENDER_EMAIL", true)
	c.Emailer.MaintenanceEmail = getEnv("CPS_IPFSSTORE_BACKEND_MAILGUN_MAINTENANCE_EMAIL", true)

	c.PaymentProcessor.SecretKey = getEnv("CPS_IPFSSTORE_BACKEND_PAYMENT_PROCESSOR_SECRET_KEY", true)
	c.PaymentProcessor.PublicKey = getEnv("CPS_IPFSSTORE_BACKEND_PAYMENT_PROCESSOR_PUBLIC_KEY", true)
	c.PaymentProcessor.WebhookSecretKey = getEnv("CPS_IPFSSTORE_BACKEND_PAYMENT_PROCESSOR_WEBHOOK_SECRET_KEY", true)

	c.IPFSNode.GatewayRPCURL = getEnv("CPS_IPFSSTORE_BACKEND_IPFS_NODE_RPC_GATEWAY_URL", false)
	c.IPFSNode.Username = getEnv("CPS_IPFSSTORE_BACKEND_IPFS_NODE_RPC_GATEWAY_USERNAME", false)
	c.IPFSNode.Password = getEnv("CPS_IPFSSTORE_BACKEND_IPFS_NODE_RPC_GATEWAY_PASSWORD", false)

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