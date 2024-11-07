package domain

import "context"

type RemoteIPFSRepository interface {
	Version(ctx context.Context) (string, error)
	PinAddViaFilepath(ctx context.Context, filepath string) (string, error)
	Get(ctx context.Context, cidString string) ([]byte, string, error)
}
