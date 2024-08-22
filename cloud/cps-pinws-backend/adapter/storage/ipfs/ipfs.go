package ipfs // Special thanks via https://github.com/ipfs/kubo/blob/master/docs/examples/kubo-as-a-library/README.md

import (
	"context"
	"log"
	"log/slog"
	"mime/multipart"

	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core"
	icore "github.com/ipfs/kubo/core/coreiface"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

type IPFSStorager interface {
	UploadContent(ctx context.Context, objectKey string, content []byte) error
	UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error
	DeleteByKeys(ctx context.Context, key []string) error
}

type ipfsStorager struct {
	ipfs       icore.CoreAPI
	node       *core.IpfsNode
	Logger     *slog.Logger
	BucketName string
}

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	// DEVELOPERS NOTE:
	// How can I use the AWS SDK v2 for Go with DigitalOcean Spaces? via https://stackoverflow.com/a/74284205
	logger.Debug("ipfs initializing...")

	// https://github.com/hsanjuan/ipfs-lite

	api, err := rpc.NewPathApi("ipfs:4001")
	logger.Debug("-> %v", api)
	if err != nil {
		log.Fatalf("failed to connect to node: %s", err)
	}

	/// --- Part I: Getting a IPFS node running

	logger.Debug("getting an ipfs node running")

	logger.Debug("IPFS node is running")

	// Create our storage handler.
	ipfsStorage := &ipfsStorager{}

	// For debugging purposes only.
	logger.Debug("ipfs initialized")

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (s *ipfsStorager) UploadContent(ctx context.Context, objectKey string, content []byte) error {
	return nil
}

func (s *ipfsStorager) UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error {
	return nil
}

func (s *ipfsStorager) DeleteByKeys(ctx context.Context, objectKeys []string) error {
	return nil
}
