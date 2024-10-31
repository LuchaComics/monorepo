package domain

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

type IPFSRepository interface {
	ID() (peer.ID, error)
	Add(fullFilePath string, shouldPin bool) (string, error)
	Pin(cidString string) error
	PinAdd(fullFilePath string) (string, error)
	Get(ctx context.Context, cidString string) ([]byte, string, error)
}
