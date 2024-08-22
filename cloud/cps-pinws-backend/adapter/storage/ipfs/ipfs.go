package ipfs // Special thanks via https://github.com/hsanjuan/ipfs-lite

import (
	"context"
	"io"
	"log"
	"log/slog"
	"mime/multipart"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

type IPFSStorager interface {
	UploadContent(ctx context.Context, objectKey string, content []byte) error
	UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error
	DeleteByKeys(ctx context.Context, key []string) error
	GetContentByCID(ctx context.Context, cidString string) ([]byte, error)
}

type ipfsStorager struct {
	datastore datastore.Batching
	peer      *ipfslite.Peer
	Logger    *slog.Logger
}

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	logger.Debug("ipfs initializing...")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load up the in-memory storage system so all fetched and saved files will
	// exists in memory and will be lost when the application terminates.
	ds := ipfslite.NewInMemoryDatastore()
	priv, _, err := crypto.GenerateKeyPair(crypto.RSA, 2048)
	if err != nil {
		panic(err)
	}
	logger.Debug("ipfs in-memory datastore ready")

	// Launches an IPFS-Lite peer.
	listen, _ := multiaddr.NewMultiaddr("/ip4/0.0.0.0/tcp/4005")
	h, dht, err := ipfslite.SetupLibp2p(
		ctx,
		priv,
		nil,
		[]multiaddr.Multiaddr{listen},
		ds,
		ipfslite.Libp2pOptionsExtra...,
	)
	if err != nil {
		log.Fatalf("failed to setup lib p2p: %s", err)
	}
	logger.Debug("ipfs connected to p2p network")
	lite, err := ipfslite.New(ctx, ds, nil, h, dht, nil)
	if err != nil {
		log.Fatalf("failed to setup ipfs peer: %s", err)
	}
	lite.Bootstrap(ipfslite.DefaultBootstrapPeers())
	logger.Debug("ipfs is peer in p2p network")

	// Create our storage handler for IPFS.
	ipfsStorage := &ipfsStorager{
		datastore: ds,
		peer:      lite,
		Logger:    logger,
	}

	// For debugging purposes only.
	logger.Debug("ipfs initialized")

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (impl *ipfsStorager) UploadContent(ctx context.Context, objectKey string, content []byte) error {
	return nil
}

func (impl *ipfsStorager) UploadContentFromMulipart(ctx context.Context, objectKey string, file multipart.File) error {
	return nil
}

func (impl *ipfsStorager) GetContentByCID(ctx context.Context, cidString string) ([]byte, error) {
	c, _ := cid.Decode("QmWATWQ7fVPP2EFGu71UkfnqhYXDYH566qy47CnJDgvs8u")
	rsc, err := impl.peer.GetFile(ctx, c)
	if err != nil {
		return nil, err
	}
	defer rsc.Close()
	content, err := io.ReadAll(rsc)
	if err != nil {
		return nil, err
	}

	// fmt.Println(string(content))
	return content, nil
}

func (impl *ipfsStorager) DeleteByKeys(ctx context.Context, objectKeys []string) error {
	return nil
}
