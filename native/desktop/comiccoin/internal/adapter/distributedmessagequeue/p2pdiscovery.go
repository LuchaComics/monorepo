package distributedmessagequeue

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func (impl *brokerImpl) discoverPeers(ctx context.Context, h host.Host, topicNameFlags []string) {
	kademliaDHT := impl.initDHT(ctx, h)
	impl.kademliaDHT = kademliaDHT
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)

	// Advertise and look for peers for each topic
	for _, topicNameFlag := range topicNameFlags {
		dutil.Advertise(ctx, routingDiscovery, topicNameFlag)

		// Look for others who have announced and attempt to connect to them
		anyConnected := false
		for !anyConnected {
			peerChan, err := routingDiscovery.FindPeers(ctx, topicNameFlag)
			if err != nil {
				impl.logger.Error("Failed routing discovery finding peers",
					slog.Any("topic", topicNameFlag),
					slog.Any("error", err))
				panic(err)
			}
			for peer := range peerChan {
				if peer.ID == h.ID() {
					continue // No self connection
				}
				err := h.Connect(ctx, peer)
				if err != nil {
					impl.logger.Error("Failed connecting to peer",
						slog.Any("topic", topicNameFlag),
						slog.Any("peer_id", peer.ID),
						slog.Any("error", err),
					)
				} else {
					impl.logger.Debug("Connected successfully to peer",
						slog.Any("topic", topicNameFlag),
						slog.Any("peer_id", peer.ID))
					anyConnected = true

					impl.mu.Lock()
					defer impl.mu.Unlock()

					// Join the topic and save it.
					topic, err := impl.ps.Join(topicNameFlag)
					if err != nil {
						panic(err)
					}
					impl.topics[topicNameFlag] = topic
				}
			}
		}
		impl.logger.Debug("Peer discovery complete for topic",
			slog.Any("topic", topicNameFlag))
	}
}
