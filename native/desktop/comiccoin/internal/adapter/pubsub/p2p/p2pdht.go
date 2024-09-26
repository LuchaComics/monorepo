package p2p

import (
	"context"
	"log"
	"time"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
)

func (node *pubSubBrokerImpl) initDHT(ctx context.Context, h host.Host) *dht.IpfsDHT {
	var options []dht.Option
	bootstrapPeers := make([]peer.AddrInfo, len(node.cfg.Peer.BootstrapPeers))
	for i, addr := range node.cfg.Peer.BootstrapPeers {
		peerinfo, _ := peer.AddrInfoFromP2pAddr(addr)
		bootstrapPeers[i] = *peerinfo
	}
	if len(bootstrapPeers) == 0 {
		options = append(options, dht.Mode(dht.ModeServer))
		node.logger.Info("Running p2p node as dht server")
	}
	options = append(options, dht.BootstrapPeers(bootstrapPeers...))

	kademliaDHT, err := dht.New(ctx, node.host, options...)
	if err != nil {
		log.Fatalf("failed createing new dht: %v", err)
	}

	// Bootstrap the DHT. In the default configuration, this spawns a Background
	// thread that will refresh the peer table every five minutes.
	node.logger.Debug("Bootstrapping the DHT")
	if err = kademliaDHT.Bootstrap(ctx); err != nil {
		log.Fatalf("failed boostraping new dht: %v", err)
	}

	// Wait a bit to let bootstrapping finish (really bootstrap should block until it's ready, but that isn't the case yet.)
	time.Sleep(1 * time.Second)

	return kademliaDHT

	//
	// Source: https://github.com/libp2p/go-libp2p/blob/master/examples/chat-with-rendezvous/chat.go#L112
	//
}