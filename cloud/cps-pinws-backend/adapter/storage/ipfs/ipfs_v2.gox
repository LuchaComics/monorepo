package ipfs // Special thanks via https://github.com/hsanjuan/ipfs-lite

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"

	ipfslite "github.com/hsanjuan/ipfs-lite"
	"github.com/ipfs/go-cid"
	datastore "github.com/ipfs/go-datastore"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/multiformats/go-multiaddr"
	ma "github.com/multiformats/go-multiaddr"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

type IPFSStorager interface {
	UploadContentFromFilepath(ctx context.Context, filepath string) (string, error)
	GetContentByCID(ctx context.Context, cidString string) ([]byte, error)
	PinContent(ctx context.Context, cidString string) error
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

	// Get the Peer ID
	peerID := h.ID()

	// Get all the multiaddresses for this host
	addrs := h.Addrs()

	// For debugging purposes only.
	logger.Debug("ipfs initialized", slog.Any("peer_id", peerID))

	// Combine the Peer ID with the multiaddresses
	for _, addr := range addrs {
		// Create a multiaddress that includes the Peer ID
		fullAddr := addr.Encapsulate(ma.StringCast(fmt.Sprintf("/p2p/%s", peerID)))

		// For debugging purposes only.
		logger.Debug("ipfs initialized", slog.Any("connection_string", fullAddr.String()))
	}

	// For debugging purpses only, here is an example of opening a local file,
	// add to ipfs and then fetching this file from ipfs again.
	sampleFilepath := "./static/blacklist/ips.json"
	sampleFileCid, err := ipfsStorage.UploadContentFromFilepath(ctx, sampleFilepath)
	if err != nil {
		log.Fatalf("failed to upload sample content to ipfs: %s", err)
	}
	sampleFileContentBytes, err := ipfsStorage.GetContentByCID(ctx, sampleFileCid)
	if err != nil {
		log.Fatalf("failed to fetch sample content from ipfs: %s", err)
	}
	fmt.Println(string(sampleFileContentBytes))
	if err := ipfsStorage.PinContent(ctx, sampleFileCid); err != nil {
		log.Fatalf("failed to pin: %s", err)
	}
	ipfsStorage.peer.Exchange()

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (impl *ipfsStorager) UploadContentFromFilepath(ctx context.Context, filepath string) (string, error) {
	impl.Logger.Debug("opening to be added to ipfs",
		slog.String("filepath", filepath))

	// Open the file at the specified filepath.
	file, err := os.Open(filepath)
	if err != nil {
		impl.Logger.Error("error opening file to be added for ipfs",
			slog.String("filepath", filepath),
			slog.Any("err", err))
		return "", fmt.Errorf("Error opening file: %v", err)
	}
	defer file.Close()

	// Define AddParams
	params := &ipfslite.AddParams{
		Layout:    "balanced",    // Use balanced DAG layout
		Chunker:   "size-262144", // Default chunk size
		RawLeaves: true,          // Store as raw leaves
		HashFun:   "sha2-256",    // Use SHA-256 for hashing
	}

	// Add the file to IPFS
	node, err := impl.peer.AddFile(ctx, file, params)
	if err != nil {
		impl.Logger.Error("error adding file to ipfs",
			slog.String("filepath", filepath),
			slog.Any("params", params),
			slog.Any("err", err))
		return "", fmt.Errorf("Error adding file to IPFS: %v", err)
	}

	impl.Logger.Debug("file added to ipfs",
		slog.String("filepath", filepath),
		slog.String("cid", node.Cid().String()))

	return node.Cid().String(), nil
}

func (impl *ipfsStorager) GetContentByCID(ctx context.Context, cidString string) ([]byte, error) {
	c, _ := cid.Decode(cidString)
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

func (impl *ipfsStorager) PinContent(ctx context.Context, cidString string) error {
	c, err := cid.Decode(cidString)
	if err != nil {
		impl.Logger.Error("error decoding cid", slog.String("cid", cidString), slog.Any("err", err))
		return fmt.Errorf("Error decoding CID: %v", err)
	}

	// Check if the block is already in the blockstore
	has, err := impl.peer.HasBlock(ctx, c)
	if err != nil {
		impl.Logger.Error("error checking if block exists", slog.String("cid", cidString), slog.Any("err", err))
		return fmt.Errorf("Error checking if block exists: %v", err)
	}

	if !has {
		// If not, fetch the block from the network (this pins the content)
		_, err := impl.peer.Get(ctx, c)
		if err != nil {
			impl.Logger.Error("error fetching and pinning content", slog.String("cid", cidString), slog.Any("err", err))
			return fmt.Errorf("Error fetching and pinning content: %v", err)
		}
	}

	impl.Logger.Debug("content pinned successfully", slog.String("cid", cidString))
	return nil
}
