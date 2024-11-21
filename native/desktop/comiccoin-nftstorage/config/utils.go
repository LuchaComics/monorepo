package config

// DEVELOPERS NOTE:
// Special thanks to the following link:
// https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/flags.go

import (
	"log"
	"os"
	"strconv"

	sbytes "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/securebytes"
	sstring "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/security/securestring"
)

func GetEnvString(key string, required bool) string {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return value
}

func GetEnvBytes(key string, required bool) []byte {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	return []byte(value)
}

func GetEnvBool(key string, required bool, defaultValue bool) bool {
	valueStr := GetEnvString(key, required)
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		log.Fatalf("Invalid boolean value for environment variable %s", key)
	}
	return value
}

func GetSecureStringEnv(key string, required bool) *sstring.SecureString {
	value := os.Getenv(key)
	if required && value == "" {
		log.Fatalf("Environment variable not found: %s", key)
	}
	ss, err := sstring.NewSecureString(value)
	if err != nil {
		log.Fatalf("Environment variable failed to secure: %v", err)
	}
	return ss
}

func GetSecureBytesEnv(key string, required bool) *sbytes.SecureBytes {
	value := GetEnvBytes(key, required)
	ss, err := sbytes.NewSecureBytes(value)
	if err != nil {
		log.Fatalf("Environment variable failed to secure: %v", err)
	}
	return ss
}
