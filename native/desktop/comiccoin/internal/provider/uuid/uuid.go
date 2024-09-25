package uuid

import (
	"fmt"

	uuid "github.com/segmentio/ksuid"
)

// Provider provides interface for abstracting UUID generation.
type Provider interface {
	NewUUID(namespace string) string
}

type uuidProvider struct {
}

// NewProvider constructor that returns the default UUID generator.
func NewProvider() Provider {
	return uuidProvider{}
}

// NewUUID generates a new UUID.
func (u uuidProvider) NewUUID(namespace string) string {
	return fmt.Sprintf("%s_%s", namespace, uuid.New().String())
}
