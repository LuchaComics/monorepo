package peer

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func (impl *peerProviderImpl) DiscoverPeersChannel(ctx context.Context, h host.Host, rendezvousString string) <-chan peer.AddrInfo {
	ch := make(chan peer.AddrInfo)
	func() {
		dutil.Advertise(ctx, impl.routingDiscovery, rendezvousString)

		// Look for others who have announced and attempt to connect to them
		anyConnected := false
		for !anyConnected {
			peerChan, err := impl.routingDiscovery.FindPeers(ctx, rendezvousString)
			if err != nil {
				impl.logger.Error("Failed routing discovery finding peers",
					slog.Any("rendezvous_string", rendezvousString),
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
						slog.Any("rendezvous_string", rendezvousString),
						slog.Any("peer_id", peer.ID),
						slog.Any("error", err),
					)
				} else {
					//
					// STEP 1
					//

					impl.logger.Debug("Connected successfully to peer",
						slog.Any("peer_id", peer.ID))

					anyConnected = true

					impl.mu.Lock()
					defer impl.mu.Unlock()

					ch <- peer
				}
			}
		}
		impl.logger.Debug("Peer discovery complete for topic")
	}()

	return ch
}
