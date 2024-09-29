package p2p

import (
	"context"
	"log"
	"log/slog"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (node *peerProviderImpl) newKademliaDHT(ctx context.Context) *dht.IpfsDHT {
	var options []dht.Option
	bootstrapPeers := make([]peer.AddrInfo, len(node.cfg.Peer.BootstrapPeers))
	for i, addr := range node.cfg.Peer.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		bootstrapPeers[i] = *peerinfo
	}
	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
		node.logger.Info("Running p2p node in host mode")
		node.isHostMode = true
	} else {
		node.logger.Info("Running p2p node in dialer mode")
	}
	options = append(options, dht.BootstrapPeers(bootstrapPeers...))

	kademliaDHT, err := dht.New(ctx, node.host, options...)
	if err != nil {
		node.logger.Debug("Failed to create new dht",
			slog.Any("error", err))
		log.Fatal(err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	node.logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		node.logger.Debug("Failed bootstraping the dht",
			slog.Any("error", err))
		log.Fatal(err)
	}

	// Wait a bit to let bootstrapping finish (really bootstrap should block until it's ready, but that isn't the case yet.)
	time.Sleep(1 * time.Second)

	return kademliaDHT

	//
	// Source: https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/chat.go#L112
	//
}
