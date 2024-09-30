package repo

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
)

const (
	topicName        = "signedtxdto"
	rendezvousString = "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/signedtxdto"
)

type signedTransactionDTORepoImpl struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	topic         *pubsub.Topic
	subs          map[peer.ID]*pubsub.Subscription
}

func NewSignedTransactionDTORepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.SignedTransactionDTORepository {
	//
	// STEP 1
	// Initialize our instance
	//

	impl := &signedTransactionDTORepoImpl{
		config:        cfg,
		logger:        logger,
		libP2PNetwork: libP2PNetwork,
		topic:         nil,
		subs:          make(map[peer.ID]*pubsub.Subscription, 0),
	}

	//
	// STEP 2:
	// Create and advertise our `rendezvousString` which is essentially telling
	// our P2P network that clients can meet and communicate in our app at this
	// specific location.
	//

	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), rendezvousString)

	//
	// STEP 3:
	// We want to implement broadcast sort of system for this signed
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

	topic, err := psObj.Join(rendezvousString)
	if err != nil {
		log.Fatalf("failed joining pub-sub for topic: %v", err)
	}
	if topic == nil {
		log.Fatal("joined pub-sub topic does not exist")
	}
	impl.topic = topic

	//
	// STEP 5:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect to topic...",
			slog.String("topic_name", topicName))

		for {

			//
			// STEP 1:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), rendezvousString, func(p peer.AddrInfo) error {
				impl.logger.Debug("connecting...",
					slog.Any("peer_id", p.ID),
					slog.String("topic_name", topicName))

				//
				// STEP 2:
				// Subscribe our peer to this topic.
				//

				sub, err := impl.topic.Subscribe()
				if err != nil {
					impl.logger.Error("failed subscribing to our topic",
						slog.Any("error", err),
						slog.Any("peer_id", p.ID),
						slog.String("topic_name", topicName))
					return err
				}
				if sub == nil {
					err := fmt.Errorf("failed subscribing to our topic: %v", "topic does not exist")
					impl.logger.Error("failed subscribing to our topic",
						slog.Any("error", err),
						slog.Any("peer_id", p.ID),
						slog.String("topic_name", topicName))
					return err
				}

				//
				// STEP 3:
				// Save our new peer subscription so we can use it later.
				//

				impl.subs[p.ID] = sub

				impl.logger.Debug("subscribed",
					slog.Any("peer_id", p.ID),
					slog.String("topic", topicName))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (impl *signedTransactionDTORepoImpl) BroadcastToP2PNetwork(ctx context.Context, bd *domain.SignedTransactionDTO) error {
	//
	// STEP 1
	// Marshal the DTO into a binary payload which we can send over the network.
	//

	bdBytes, err := bd.Serialize()
	if err != nil {
		impl.logger.Error("Failed to publish",
			slog.String("topic_name", topicName),
			slog.Any("error", err))
		return err
	}

	if err := impl.topic.Publish(ctx, bdBytes); err != nil {
		impl.logger.Error("Failed to publish",
			slog.String("topic_name", topicName),
			slog.Any("error", err))
		return fmt.Errorf("failed to publish: %s", topicName)
	}
	impl.logger.Debug("Published",
		slog.Any("topic", topicName))

	return nil
}

func (impl *signedTransactionDTORepoImpl) ReceiveFromP2PNetwork(ctx context.Context) (*domain.SignedTransactionDTO, error) {
	//
	// STEP 1:
	// We will iterate over each subscription we have and listen for incoming
	// messages.
	//

	for _, sub := range impl.subs {
		//
		// STEP 2:
		// We will receive the incoming message from our P2P network.
		//

		msg, err := sub.Next(ctx)
		if err != nil {
			impl.logger.Error("Failed to receive message",
				slog.Any("error", err),
				slog.String("topic_name", topicName))
			continue
		}

		//
		// STEP 3:
		// We need to deserialize the incoming message into our DTO.
		//

		stxDTO, err := domain.NewSignedTransactionDTOFromDeserialize(msg.Data)
		if err != nil {
			impl.logger.Error("Failed to deserialize message",
				slog.Any("error", err),
				slog.String("topic_name", topicName))
			continue
		}

		//
		// STEP 4:
		// Return the DTO.
		//

		return stxDTO, nil
	}

	//
	// STEP 5:
	// If we didn't receive any messages then return without any problems.
	//

	return nil, nil
}
