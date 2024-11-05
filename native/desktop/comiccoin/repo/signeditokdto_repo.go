package repo

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/net/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

const (
	signedIssuedTokenDTOTopicName        = "signedissuedtokendto"
	signedIssuedTokenDTORendezvousString = "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain/signedissuedtokendto"
)

type signedIssuedTokenDTORepoImpl struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	topic         *pubsub.Topic
	sub           *pubsub.Subscription
}

func NewSignedIssuedTokenDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.SignedIssuedTokenDTORepository {
	//
	// STEP 1
	// Initialize our instance
	//

	impl := &signedIssuedTokenDTORepoImpl{
		config:        cfg,
		logger:        logger,
		libP2PNetwork: libP2PNetwork,
		topic:         nil,
		sub:           nil,
	}

	//
	// STEP 2:
	// Create and advertise our `signedIssuedTokenDTORendezvousString` which is essentially telling
	// our P2P network that clients can meet and communicate in our app at this
	// specific location.
	//

	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), signedIssuedTokenDTORendezvousString)

	//
	// STEP 3:
	// We want to implement broadcast sort of system for this mempool
	// transaction data-transfer object; meaning, any one node in the P2P
	// network can transmit to all the nodes on the P2P network this data.
	//
	// To accomplish this functionality we will use the `pub-sub` pattern.
	// Begin by getting out pub-sub instance.
	//

	psObj := impl.libP2PNetwork.GetPubSubSingletonInstance()
	if psObj == nil {
		log.Fatal("failed to get pub-sub object")
	}

	//
	// STEP 4:
	// Join the `topic` in the pub-sub.
	//

	topic, err := psObj.Join(signedIssuedTokenDTORendezvousString)
	if err != nil {
		log.Fatalf("failed joining pub-sub for topic: %v", err)
	}
	if topic == nil {
		log.Fatal("joined pub-sub topic does not exist")
	}
	impl.topic = topic

	//
	// STEP 5:
	// Subscribe our peer to this topic.
	//

	sub, err := topic.Subscribe()
	if err != nil {
		impl.logger.Error("failed subscribing to our topic",
			slog.Any("error", err),
			slog.String("topic_name", signedIssuedTokenDTOTopicName))
		log.Fatalf("failed subscribing to our topic: %v", err)
	}
	if sub == nil {
		err := fmt.Errorf("failed subscribing to our topic: %v", "topic does not exist")
		impl.logger.Error("failed subscribing to our topic",
			slog.Any("error", err),
			slog.String("topic_name", signedIssuedTokenDTOTopicName))
		log.Fatalf("failed subscribing to our topic: %v", err)
	}
	impl.sub = sub

	//
	// STEP 5:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect to topic...",
			slog.String("topic_name", signedIssuedTokenDTOTopicName))

		for {

			//
			// STEP 1:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), signedIssuedTokenDTORendezvousString, func(p peer.AddrInfo) error {

				impl.logger.Debug("subscribed",
					slog.Any("peer_id", p.ID),
					slog.String("dto", "signedissuedtokendto"),
					slog.String("topic", signedIssuedTokenDTOTopicName))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (impl *signedIssuedTokenDTORepoImpl) BroadcastToP2PNetwork(ctx context.Context, bd *domain.SignedIssuedTokenDTO) error {
	//
	// STEP 1
	// Marshal the DTO into a binary payload which we can send over the network.
	//

	bdBytes, err := bd.Serialize()
	if err != nil {
		impl.logger.Error("Failed to publish",
			slog.String("topic_name", signedIssuedTokenDTOTopicName),
			slog.Any("error", err))
		return err
	}

	// Developers Note:
	// https://github.com/libp2p/go-libp2p/blob/master/examples/pubsub/basic-chat-with-rendezvous/main.go#L115

	if err := impl.topic.Publish(ctx, bdBytes); err != nil {
		impl.logger.Error("Failed to publish",
			slog.String("topic_name", signedIssuedTokenDTOTopicName),
			slog.Any("error", err))
		return fmt.Errorf("failed to publish: %s", signedIssuedTokenDTOTopicName)
	}
	impl.logger.Debug("Published",
		slog.Any("topic", signedIssuedTokenDTOTopicName))

	return nil
}

func (impl *signedIssuedTokenDTORepoImpl) ReceiveFromP2PNetwork(ctx context.Context) (*domain.SignedIssuedTokenDTO, error) {
	//
	// STEP 2:
	// We will receive the incoming message from our P2P network.
	//

	// Developers Note:
	// https://github.com/libp2p/go-libp2p/blob/master/examples/pubsub/basic-chat-with-rendezvous/main.go#L121

	msg, err := impl.sub.Next(ctx)
	if err != nil {
		impl.logger.Error("Failed to receive message",
			slog.Any("error", err),
			slog.String("topic_name", signedIssuedTokenDTOTopicName))
		return nil, err
	}

	//
	// STEP 3:
	// We need to deserialize the incoming message into our DTO.
	//

	stxDTO, err := domain.NewSignedIssuedTokenDTOFromDeserialize(msg.Data)
	if err != nil {
		impl.logger.Error("Failed to deserialize message",
			slog.Any("error", err),
			slog.String("topic_name", signedIssuedTokenDTOTopicName))
		return nil, err
	}

	//
	// STEP 4:
	// Return the DTO.
	//

	return stxDTO, nil
}
