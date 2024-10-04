package p2p

import (
	"context"
	"log"
	"log/slog"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
)

// newKademliaDHT creates a new instance of the Kademlia DHT.
// The DHT is used for peer discovery and data storage.
func (node *peerProviderImpl) newKademliaDHT(ctx context.Context) *dht.IpfsDHT {
	// Create a list of bootstrap peers from the configuration.
	var options []dht.Option
	bootstrapPeers := make([]peer.AddrInfo, len(node.cfg.Peer.BootstrapPeers))
	for i, addr := range node.cfg.Peer.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		bootstrapPeers[i] = *peerinfo
	}

	// If no bootstrap peers are specified, run the node in host mode.
	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
		node.logger.Info("Running p2p node in host mode")
		node.isHostMode = true
	} else {
		node.logger.Info("Running p2p node in dialer mode")
	}

	// Add the bootstrap peers to the DHT options.
	options = append(options, dht.BootstrapPeers(bootstrapPeers...))

	// Create a new instance of the Kademlia DHT.
	kademliaDHT, err := dht.New(ctx, node.host, options...)
	if err != nil {
		node.logger.Debug("Failed to create new dht",
			slog.Any("error", err))
		log.Fatal(err)
	}

	// Bootstrap the DHT to populate the peer table.
	node.logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		node.logger.Debug("Failed bootstrapping the dht",
			slog.Any("error", err))
		log.Fatal(err)
	}

	// Wait for the bootstrapping process to complete.
	time.Sleep(1 * time.Second)

	return kademliaDHT
}

// PutDataToKademliaDHT puts data into the Kademlia DHT.
// The data is stored under the given key.
func (impl *peerProviderImpl) PutDataToKademliaDHT(key string, bytes []byte) error {
	// Put the data into the DHT.
	if err := impl.kademliaDHT.PutValue(context.Background(), key, bytes); err != nil {
		impl.logger.Error("failed putting to kademlia dht",
			slog.Any("error", err),
		)
		return err
	}
	return nil
}

// GetDataFromKademliaDHT gets data from the Kademlia DHT.
// The data is retrieved under the given key.
func (impl *peerProviderImpl) GetDataFromKademliaDHT(key string) ([]byte, error) {
	// Get the data from the DHT.
	bytes, err := impl.kademliaDHT.GetValue(context.Background(), key)
	if err != nil {
		impl.logger.Error("failed getting from kademlia dht",
			slog.Any("error", err),
		)
		return nil, err
	}
	return bytes, nil
}

// RemoveDataFromKademliaDHT removes data from the Kademlia DHT.
// The data is removed under the given key.
func (impl *peerProviderImpl) RemoveDataFromKademliaDHT(key string) error {
	// Remove the data from the DHT by putting an empty value under the key.
	if err := impl.kademliaDHT.PutValue(context.Background(), key, []byte{}); err != nil {
		impl.logger.Error("failed removing from kademlia dht",
			slog.Any("error", err),
		)
		return err
	}
	return nil
}
