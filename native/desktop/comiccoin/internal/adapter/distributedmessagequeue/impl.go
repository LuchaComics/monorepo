package distributedmessagequeue

import (
	"bufio"
	"context"
	"fmt"
	"log/slog"
	"os"
	"sync"

	dht "github.com/libp2p/go-libp2p-kad-dht"
	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/discovery/routing"

	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
)

type brokerImpl struct {
	cfg           *config.Config
	logger        *slog.Logger
	keypairStorer keypair_ds.KeypairStorer

	// mu field is used to protect access to the subs and closed fields using a mutex.
	mu sync.Mutex

	topics map[string]*pubsub.Topic

	// The quit field is a channel that is closed when the `message queue broker` is closed, allowing goroutines that are blocked on the channel to unblock and exit.
	quit chan struct{}

	// The closed field is a flag that indicates whether the `message queue broker` has been closed.
	closed bool

	host             host.Host
	kademliaDHT      *dht.IpfsDHT
	routingDiscovery *routing.RoutingDiscovery

	ps *pubsub.PubSub
}

func NewDistributedMessageQueueAdapter(cfg *config.Config, logger *slog.Logger, kp keypair_ds.KeypairStorer) DistributedMessageQueueBroker {
	ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
	node := &brokerImpl{
		cfg:           cfg,
		logger:        logger,
		keypairStorer: kp,
		topics:        make(map[string]*pubsub.Topic, 0),
	}
	h, err := node.newHostWithPredictableIdentifier()
	node.host = h

	// topic1 := flag.String("topicName", "applesauce", "name of topic to join")
	topics := []string{"mempool", "sync"}
	go node.discoverPeers(ctx, h, topics)

	ps, err := pubsub.NewGossipSub(ctx, h)
	if err != nil {
		panic(err)
	}
	node.ps = ps

	return node
}

func (impl *brokerImpl) Subscribe(ctx context.Context, topicName string) []byte {
	topic, ok := impl.topics[topicName]
	if !ok {
		panic("failed Subscribe")
	}
	sub, err := topic.Subscribe()
	if err != nil {
		panic(err)
	}

	m, err := sub.Next(ctx)
	if err != nil {
		panic(err)
	}
	return m.Message.Data
}

func (impl *brokerImpl) Publish(ctx context.Context, topicName string, msg []byte) {
	topic, ok := impl.topics[topicName]
	if !ok {
		panic("failed Publish")
	}

	if err := topic.Publish(ctx, msg); err != nil {
		fmt.Println("### Publish error:", err)
	}
}

func (impl *brokerImpl) Close() {
	// TODO: IMPL.
}

func writeMessagesToQueue(ctx context.Context, queue chan<- string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')
		if err != nil {
			panic(err)
		}
		queue <- s
	}
}

func readMessagesFromQueue(ctx context.Context, queue <-chan string) {
	for {
		select {
		case s := <-queue:
			fmt.Println("Received message:", s)
		case <-ctx.Done():
			return
		}
	}
}
