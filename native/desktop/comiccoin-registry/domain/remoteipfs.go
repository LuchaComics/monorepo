package domain

import "context"

// RemoteIPFSGetFileResponse represents the structure of a response when retrieving a file from IPFS.
// This struct includes metadata and the actual file content returned by the IPFS server.
type RemoteIPFSGetFileResponse struct {
	Filename      string `json:"filename"`       // The name of the file being retrieved, derived from Content-Disposition headers.
	Content       []byte `json:"content"`        // The raw file content in bytes.
	ContentType   string `json:"content_type"`   // The MIME type of the file, determined by server response headers.
	ContentLength uint64 `json:"content_length"` // The length of the content in bytes.
}

// RemoteIPFSRepository defines an interface for interacting with a remote IPFS server.
// This interface abstracts common IPFS operations, including retrieving the IPFS server version,
// pinning a file by its filepath, and fetching a file by CID.
type RemoteIPFSRepository interface {
	// Version fetches the current version of the remote IPFS server.
	// It returns the version as a string and any error encountered.
	Version(ctx context.Context) (string, error)

	// PinAddViaFilepath pins a file located at the specified filepath to the IPFS server.
	// The filepath should be a full local path, and it returns the CID of the pinned file or an error.
	PinAddViaFilepath(ctx context.Context, filepath string) (string, error)

	// Get retrieves a file from the IPFS server using its CID (Content Identifier).
	// It returns a RemoteIPFSGetFileResponse containing file metadata, content, and any error encountered.
	Get(ctx context.Context, cidString string) (*RemoteIPFSGetFileResponse, error)
}
