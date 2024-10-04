package p2p

import (
	"context"
	"log"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

// AdvertiseWithRendezvousString advertises the peer's presence using a rendezvous string.
// This allows other peers to discover and connect to this peer.
func (impl *peerProviderImpl) AdvertiseWithRendezvousString(ctx context.Context, h host.Host, rendezvousString string) {
	// Use the libp2p discovery utility to advertise the peer's presence.
	dutil.Advertise(ctx, impl.routingDiscovery, rendezvousString)
}

// DiscoverPeersAtRendezvousString discovers peers at a rendezvous string and connects to them.
// The peerConnectedFunc callback function is called for each peer that is connected.
func (impl *peerProviderImpl) DiscoverPeersAtRendezvousString(ctx context.Context, h host.Host, rendezvousString string, peerConnectedFunc func(peer.AddrInfo) error) {
	// Initialize the map for the rendezvous string if necessary.
	_, ok := impl.peers[rendezvousString]
	if !ok {
		impl.peers[rendezvousString] = make(map[peer.ID]*peer.AddrInfo, 0)
	}

	// Find peers at the rendezvous string using the routing discovery.
	peerChan, err := impl.routingDiscovery.FindPeers(ctx, rendezvousString)
	if err != nil {
		impl.logger.Error("Failed routing discovery finding peers",
			slog.Any("rendezvous_string", rendezvousString),
			slog.Any("error", err))
		log.Fatalf("Failed routing discovery finding peers: %v", err)
	}

	// Iterate over the peers found at the rendezvous string.
	for peer := range peerChan {
		// Skip self connections.
		if peer.ID == h.ID() {
			continue
		}

		// Check if we have already connected to this peer.
		_, ok := impl.peers[rendezvousString][peer.ID]
		if !ok {
			// Connect to the peer.
			err := h.Connect(ctx, peer)
			if err != nil {
				// Note: This error is "acceptable" because libp2p is designed to
				// "loose" peer information over time. See:
				// https://discuss.libp2p.io/t/disconnecting-removing-peers-form-the-dht-and-peerstore/1932/4
				continue
			} else {
				// Call the peerConnectedFunc callback function.
				if err := peerConnectedFunc(peer); err != nil {
					// impl.logger.Error("failed connecting peer", slog.Any("error", err))
					break
				}

				// Store the peer in the map.
				impl.peers[rendezvousString][peer.ID] = &peer
			}
		}
	}
}
