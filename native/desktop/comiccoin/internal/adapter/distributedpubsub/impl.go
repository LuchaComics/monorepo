package distributedpubsub

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"

	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type pubSubBrokerImpl struct {
	cfg           *config.Config
	logger        *slog.Logger
	keypairStorer keypair_ds.KeypairStorer

	// mu field is used to protect access to the subs and closed fields using a mutex.
	mu sync.Mutex

	// subs field is a map that holds a list of channels for each topic, allowing subscribers to receive messages published to that topic.
	subs map[string][]chan []byte

	// The quit field is a channel that is closed when the `message queue broker` is closed, allowing goroutines that are blocked on the channel to unblock and exit.
	quit chan struct{}

	// The closed field is a flag that indicates whether the `message queue broker` has been closed.
	closed bool

	topics           map[string]*pubsub.Topic
	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *routing.RoutingDiscovery

	gossipPubSub *pubsub.PubSub
}

func NewAdapter(cfg *config.Config, logger *slog.Logger, kp keypair_ds.KeypairStorer) PublishSubscribeBroker {
	ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
	node := &pubSubBrokerImpl{
		cfg:           cfg,
		logger:        logger,
		keypairStorer: kp,
		topics:        make(map[string]*pubsub.Topic, 0),
	}

	// Run the code which will setup our peer-to-peer node in either `host mode`
	// or `dial mode`.
	h, err := node.newHostWithPredictableIdentifier()
	node.host = h

	topics := []string{"mempool"}
	go node.discoverPeers(ctx, h, topics)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	node.gossipPubSub = ps

	return node
}

func (impl *pubSubBrokerImpl) Subscribe(ctx context.Context, topicName string) <-chan []byte {
	impl.mu.Lock()
	defer impl.mu.Unlock()

	if impl.closed {
		return nil
	}

	ch := make(chan []byte)
	impl.subs[topicName] = append(impl.subs[topicName], ch)
	return ch
}

func (impl *pubSubBrokerImpl) Publish(ctx context.Context, topicName string, msg []byte) error {
	topic, ok := impl.topics[topicName]
	if !ok {
		impl.logger.Error("Failed to get topic because d.n.e.",
			slog.String("topic_name", topicName))
		return fmt.Errorf("Failed to get topic because d.n.e.: %s", topicName)
	}
	if err := topic.Publish(ctx, msg); err != nil {
		impl.logger.Error("Failed to publish",
			slog.String("topic_name", topicName),
			slog.Any("error", err))
		return fmt.Errorf("failed to publish: %s", topicName)
	}
	impl.logger.Debug("Published",
		slog.Any("topic", topicName))

	return nil
}

func (impl *pubSubBrokerImpl) Close() {
	impl.mu.Lock()
	defer impl.mu.Unlock()

	if impl.closed {
		return
	}

	impl.closed = true
	close(impl.quit)

	// Close our active channels.
	for _, ch := range impl.subs {
		for _, sub := range ch {
			close(sub)
		}
	}

	// Close our network channels.
	for _, topic := range impl.topics {
		topic.Close()

	}
}
