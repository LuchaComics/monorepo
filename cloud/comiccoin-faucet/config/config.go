package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/ethereum/go-ethereum/common"

	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
)

type Configuration struct {
	App        serverConf
	Blockchain BlockchainConfig
	DB         dbConfig
}

type serverConf struct {
	DataDirectory string
	Port          string
	IP            string
	HTTPAddress   string
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

	// (Only set by PoA node)
	ProofOfAuthorityAccountAddress *common.Address

	// (Only set by PoA node)
	ProofOfAuthorityWalletPassword *sstring.SecureString
}

type dbConfig struct {
	URI  string
	Name string
}

func NewProvider() *Configuration {
	var c Configuration

	// Application section.
	c.App.DataDirectory = getEnv("COMICCOIN_FAUCET_APP_DATA_DIRECTORY", true)
	c.App.Port = getEnv("COMICCOIN_FAUCET_PORT", true)
	c.App.IP = getEnv("COMICCOIN_FAUCET_IP", false)
	c.App.HTTPAddress = fmt.Sprintf("%v:%v", c.App.IP, c.App.Port)

	// Blockchain section.
	chainID, _ := strconv.ParseUint(getEnv("COMICCOIN_FAUCET_BLOCKCHAIN_CHAIN_ID", true), 10, 16)
	c.Blockchain.ChainID = uint16(chainID)
	transPerBlock, _ := strconv.ParseUint(getEnv("COMICCOIN_FAUCET_BLOCKCHAIN_TRANS_PER_BLOCK", true), 10, 16)
	c.Blockchain.TransPerBlock = uint16(transPerBlock)
	difficulty, _ := strconv.ParseUint(getEnv("COMICCOIN_FAUCET_BLOCKCHAIN_DIFFICULTY", true), 10, 16)
	c.Blockchain.Difficulty = uint16(difficulty)
	c.Blockchain.MiningReward, _ = strconv.ParseUint(getEnv("COMICCOIN_FAUCET_BLOCKCHAIN_MINING_REWARD", false), 10, 64)
	c.Blockchain.GasPrice, _ = strconv.ParseUint(getEnv("COMICCOIN_FAUCET_BLOCKCHAIN_GAS_PRICE", false), 10, 64)
	c.Blockchain.UnitsOfGas, _ = strconv.ParseUint(getEnv("COMICCOIN_FAUCET_BLOCKCHAIN_UNITS_OF_GAS", false), 10, 64)
	proofOfAuthorityAccountAddress := getEnv("COMICCOIN_FAUCET_ACCOUNT_ADDRESS", false)
	if proofOfAuthorityAccountAddress != "" {
		address := common.HexToAddress(proofOfAuthorityAccountAddress)
		c.Blockchain.ProofOfAuthorityAccountAddress = &address
	}
	c.Blockchain.ProofOfAuthorityWalletPassword = getSecureStringEnv("COMICCOIN_FAUCET_WALLET_PASSWORD", false)
	// Database section.
	c.DB.URI = getEnv("COMICCOIN_FAUCET_DB_URI", true)
	c.DB.Name = getEnv("COMICCOIN_FAUCET_DB_NAME", true)

	return &c
}

func getEnv(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return value
}

func getSecureStringEnv(key string, required bool) *sstring.SecureString {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	ss, err := sstring.NewSecureString(value)
	if err != nil {
		log.Fatalf("Environment variable `%v` failed to secure: %v", key, err)
	}
	return ss
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
