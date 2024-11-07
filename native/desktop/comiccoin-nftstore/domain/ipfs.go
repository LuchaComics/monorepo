// Copy and pasted from `comiccoin`
package domain

import (
	"context"
	"os"

	"github.com/ipfs/boxo/files"
	"github.com/libp2p/go-libp2p/core/peer"
)

type IPFSRepository interface {
	ID() (peer.ID, error)
	AddViaFilePath(fullFilePath string, shouldPin bool) (string, error)
	AddViaFileContent(fileContent []byte, shouldPin bool) (string, error)
	AddViaFile(file *os.File, shouldPin bool) (string, error)
	AddViaReaderFile(node files.File, shouldPin bool) (string, error)
	Pin(cidString string) error
	PinAddViaFilePath(fullFilePath string) (string, error)
	Get(ctx context.Context, cidString string) ([]byte, string, error)
}
