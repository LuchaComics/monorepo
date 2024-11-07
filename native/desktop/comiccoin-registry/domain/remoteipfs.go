package domain

import "context"

type RemoteIPFSGetFileResponse struct {
	Filename      string `json:"filename"`
	Content       []byte `json:"content"`
	ContentType   string `json:"content_type"`
	ContentLength uint64 `json:"content_length"`
}

type RemoteIPFSRepository interface {
	Version(ctx context.Context) (string, error)
	PinAddViaFilepath(ctx context.Context, filepath string) (string, error)
	Get(ctx context.Context, cidString string) (*RemoteIPFSGetFileResponse, error)
}
