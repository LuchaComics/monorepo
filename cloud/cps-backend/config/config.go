package config

import (
	"log"
	"os"
	"strconv"
)

type Conf struct {
	AppServer        serverConf
	DB               dbConfig
	AWS              awsConfig
	PDFBuilder       pdfBuilderConfig
	Emailer          mailgunConfig
	PaymentProcessor paymentProcessorConfig
}

type serverConf struct {
	Port                    string
	IP                      string
	HMACSecret              []byte
	HasDebugging            bool
	InitialAdminEmail       string
	InitialAdminPassword    string
	InitialAdminStoreName   string
	APIDomainName           string
	AppDomainName           string
	IsDeveloperMode         bool
	Enable2FAOnRegistration bool
}

type dbConfig struct {
	URI  string
	Name string
}

type awsConfig struct {
	AccessKey  string
	SecretKey  string
	Endpoint   string
	Region     string
	BucketName string
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

func New() *Conf {
	var c Conf
	c.AppServer.IsDeveloperMode = getEnvBool("CPS_BACKEND_APP_IS_DEVELOPER_MODE", false, true) // If in doubt assume developer mode!
	c.AppServer.Port = getEnv("CPS_BACKEND_PORT", true)
	c.AppServer.IP = getEnv("CPS_BACKEND_IP", false)
	c.AppServer.HMACSecret = []byte(getEnv("CPS_BACKEND_HMAC_SECRET", true))
	c.AppServer.HasDebugging = getEnvBool("CPS_BACKEND_HAS_DEBUGGING", true, true)
	c.AppServer.InitialAdminEmail = getEnv("CPS_BACKEND_INITIAL_ADMIN_EMAIL", true)
	c.AppServer.InitialAdminPassword = getEnv("CPS_BACKEND_INITIAL_ADMIN_PASSWORD", true)
	c.AppServer.InitialAdminStoreName = getEnv("CPS_BACKEND_INITIAL_ADMIN_ORG_NAME", true)
	c.AppServer.APIDomainName = getEnv("CPS_BACKEND_API_DOMAIN_NAME", true)
	c.AppServer.AppDomainName = getEnv("CPS_BACKEND_APP_DOMAIN_NAME", true)
	c.AppServer.Enable2FAOnRegistration = getEnvBool("CPS_BACKEND_APP_ENABLE_2FA_ON_REGISTRATION", false, false)

	c.DB.URI = getEnv("CPS_BACKEND_DB_URI", true)
	c.DB.Name = getEnv("CPS_BACKEND_DB_NAME", true)

	c.AWS.AccessKey = getEnv("CPS_BACKEND_AWS_ACCESS_KEY", true)
	c.AWS.SecretKey = getEnv("CPS_BACKEND_AWS_SECRET_KEY", true)
	c.AWS.Endpoint = getEnv("CPS_BACKEND_AWS_ENDPOINT", true)
	c.AWS.Region = getEnv("CPS_BACKEND_AWS_REGION", true)
	c.AWS.BucketName = getEnv("CPS_BACKEND_AWS_BUCKET_NAME", true)

	c.PDFBuilder.CBFFTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CBFF_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.PCTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_PC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCIMGTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CCIMG_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCSCTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CCSC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CC_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.CCUGTemplatePath = getEnv("CPS_BACKEND_PDF_BUILDER_CCUG_TEMPLATE_FILE_PATH", true)
	c.PDFBuilder.DataDirectoryPath = getEnv("CPS_BACKEND_PDF_BUILDER_DATA_DIRECTORY_PATH", true)

	c.Emailer.APIKey = getEnv("CPS_BACKEND_MAILGUN_API_KEY", true)
	c.Emailer.Domain = getEnv("CPS_BACKEND_MAILGUN_DOMAIN", true)
	c.Emailer.APIBase = getEnv("CPS_BACKEND_MAILGUN_API_BASE", true)
	c.Emailer.SenderEmail = getEnv("CPS_BACKEND_MAILGUN_SENDER_EMAIL", true)
	c.Emailer.MaintenanceEmail = getEnv("CPS_BACKEND_MAILGUN_MAINTENANCE_EMAIL", true)

	c.PaymentProcessor.SecretKey = getEnv("CPS_BACKEND_PAYMENT_PROCESSOR_SECRET_KEY", true)
	c.PaymentProcessor.PublicKey = getEnv("CPS_BACKEND_PAYMENT_PROCESSOR_PUBLIC_KEY", true)
	c.PaymentProcessor.WebhookSecretKey = getEnv("CPS_BACKEND_PAYMENT_PROCESSOR_WEBHOOK_SECRET_KEY", true)

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
