package config

// DEVELOPERS NOTE:
// Special thanks to the following link:
// https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/flags.go

import (
	"log"
	"os"
	"strconv"
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
