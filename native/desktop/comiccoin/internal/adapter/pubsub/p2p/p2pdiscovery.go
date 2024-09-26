package p2p

import (
	"context"
	"log/slog"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	drouting "github.com/libp2p/go-libp2p/p2p/discovery/routing"
	dutil "github.com/libp2p/go-libp2p/p2p/discovery/util"
)

func (impl *pubSubBrokerImpl) discoverPeers(ctx context.Context, h host.Host, topicNames []string) {
	kademliaDHT := impl.initDHT(ctx, h)
	impl.kademliaDHT = kademliaDHT
	routingDiscovery := drouting.NewRoutingDiscovery(kademliaDHT)

	// Advertise and look for peers for each topic
	for _, topicName := range topicNames {
		dutil.Advertise(ctx, routingDiscovery, topicName)

		// Look for others who have announced and attempt to connect to them
		anyConnected := false
		for !anyConnected {
			peerChan, err := routingDiscovery.FindPeers(ctx, topicName)
			if err != nil {
				impl.logger.Error("Failed routing discovery finding peers",
					slog.Any("topic", topicName),
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
						slog.Any("topic", topicName),
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

					//
					// STEP 2
					// Join the topic and save it.
					//

					topic, err := impl.gossipPubSub.Join(topicName)
					if err != nil {
						panic(err)
					}
					impl.topics[topicName] = topic

					impl.logger.Debug("Joined",
						slog.Any("topic", topicName),
						slog.Any("peer_id", peer.ID))

					//
					// STEP 3
					//

					sub, err := topic.Subscribe()
					if err != nil {
						panic(err)
					}

					impl.logger.Debug("Subscribed",
						slog.Any("topic", topicName),
						slog.Any("peer_id", peer.ID))

					go impl.streamSubscribeResponsesFromNetwork(ctx, topicName, sub)
				}
			}
		}
		impl.logger.Debug("Peer discovery complete for topic",
			slog.Any("topic", topicName))
	}
}

func (impl *pubSubBrokerImpl) streamSubscribeResponsesFromNetwork(ctx context.Context, topicName string, sub *pubsub.Subscription) {
	for {
		//
		// STEP 1
		// Block the flow of execution in this function until we receive a
		// response from the networker publisher and then release flow
		// execution in this function.

		m, err := sub.Next(ctx)
		if err != nil {
			impl.logger.Error("Failed to get next response",
				slog.Any("error", err))

			// Restart the loop by stopping execution in this current loop and
			// then reset to the start of this loop.
			continue
		}

		//
		// STEP 2
		// Perform goroutine coordination by waiting flow of execution.
		//

		impl.mu.Lock()
		defer impl.mu.Unlock()

		if impl.closed {
			return
		}

		//
		// STEP 3
		// Take the publisher's subscribe contents and send them to all the
		// subscribers we are managing.
		//

		for _, ch := range impl.subs[topicName] {
			ch <- m.Message.Data
		}
	}
}
