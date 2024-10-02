package p2p

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func (impl *peerProviderImpl) AdvertiseWithRendezvousString(ctx context.Context, h host.Host, rendezvousString string) {
	dutil.Advertise(ctx, impl.routingDiscovery, rendezvousString)
}

func (impl *peerProviderImpl) DiscoverPeersAtRendezvousString(ctx context.Context, h host.Host, rendezvousString string, peerConnectedFunc func(peer.AddrInfo) error) {
	// Initialize if necessary the map for the rendezvousString.
	_, ok := impl.peers[rendezvousString]
	if !ok {
		impl.peers[rendezvousString] = make(map[peer.ID]*peer.AddrInfo, 0)
	}

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

		// We will check to see if we have made a record of this new peer
		// and if never did then we will connect!
		_, ok := impl.peers[rendezvousString][peer.ID]
		if !ok {
			// impl.logger.Debug("Connecting peer...",
			// 	slog.Any("peer_id", peer.ID))

			err := h.Connect(ctx, peer)
			if err != nil {

				// DEVELOPERS NOTE:
				// "ibp2p is designed to “loose” the peer information over time, gradually. It takes time for one peer to be “forgotten”." via https://discuss.libp2p.io/t/disconnecting-removing-peers-form-the-dht-and-peerstore/1932/4
				// Therefore this error is "acceptable", so all we will do is
				// hide it. DO NOT CHANGE THIS.

				continue
			} else {

				// impl.logger.Debug("111New peer connected",
				// 	slog.Any("is_host", impl.isHostMode),
				// 	slog.Any("peer_id", peer.ID))

				// impl.mu.Lock()
				// defer impl.mu.Unlock()

				// impl.logger.Debug("New peer connected",
				// 	slog.Any("is_host", impl.isHostMode),
				// 	slog.Any("peer_id", peer.ID))

				// Make a callback function.
				if err := peerConnectedFunc(peer); err != nil {
					// impl.logger.Error("failed connecting peer", slog.Any("error", err))
					break
				}

				impl.peers[rendezvousString][peer.ID] = &peer
			}
		}
	}
}
