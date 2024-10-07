package repo

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"sync"
	"time"

	pubsub "github.com/libp2p/go-libp2p-pubsub"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p/protocol/simple"
)

// MajorityVoteConsensusRepoImpl represents the implementation of the
// `ConsensusRepository` interface tailer for a `majority voting` consensum
// algorithm.
type MajorityVoteConsensusRepoImpl struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork

	mu sync.Mutex

	// The list of connected peer addresses. This is important for consenus
	// because we will need to ensure that at least more then 50% of the
	// connected peers at any given time report the same hash for the
	// network to estable a majority voting consensus.
	peers map[peer.ID]*peer.AddrInfo

	// Establish a simple DTO protocol to listen for requests from netowrk
	// and make direct response to the peers whom made consensus requsts.
	dtoProtocol simple.SimpleDTOProtocol

	// Variable that will be shared by this node to the network so whenever
	// another peer node asks for the current blockchain hash, this node will
	// automatically respond with this value.
	currentBlockchainHash string

	// Pub-Sub topic used to send/receive broadcast announcment that a node
	// requests a consensus on the network.
	topic *pubsub.Topic

	// The subscription to the network broadcasts for consensus.
	sub *pubsub.Subscription
}

func NewMajorityVoteConsensusRepoImpl(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.ConsensusRepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/consensus"

	//
	// STEP 1
	// Initialize our instance
	//

	impl := &MajorityVoteConsensusRepoImpl{
		config:        cfg,
		logger:        logger,
		libP2PNetwork: libP2PNetwork,
		peers:         make(map[peer.ID]*peer.AddrInfo, 0),
		topic:         nil,
		sub:           nil,
	}

	//
	// STEP 2:
	// Create and advertise our `rendezvousString` which is essentially telling
	// our P2P network that clients can meet and communicate in our app at this
	// specific location.
	//

	// This is like your friend telling you the location to meet you.
	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), rendezvousString)

	//
	// STEP 3:
	// Load up all the stream handlers by this peer.
	//

	h := libP2PNetwork.GetHost()

	//
	// STEP 4:
	// In a peer-to-peer network we need to be away of when peers disconnect
	// our network, the following code will callback when a peer disconnects so
	// our repository can remove the peer from our records.
	//

	//Remove disconnected peer
	h.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected",
				slog.Any("peer_id", peerID),
				slog.String("dto", "consensus"),
			)
			delete(impl.peers, peerID)

		},
	})

	//
	// STEP 5:
	// Settup our simple DTO protocol.
	//

	dtoProtocol := simple.NewSimpleDTOProtocol(logger, h, "/consensus/req/1.0.0", "/consensus/resp/1.0.0")
	impl.dtoProtocol = dtoProtocol

	//
	// STEP 6:
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
	// STEP 7:
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
	// STEP 8:
	// Subscribe our peer to this topic.
	//

	sub, err := topic.Subscribe()
	if err != nil {
		impl.logger.Error("failed subscribing to our topic",
			slog.Any("error", err))
		log.Fatalf("failed subscribing to our topic: %v", err)
	}
	if sub == nil {
		err := fmt.Errorf("failed subscribing to our topic: %v", "topic does not exist")
		impl.logger.Error("failed subscribing to our topic",
			slog.Any("error", err))
		log.Fatalf("failed subscribing to our topic: %v", err)
	}
	impl.sub = sub

	//
	// STEP 9:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect...")

		for {

			//
			// STEP 10:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), h, rendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 11
				// Connect our peer with this peer and keep a record of it.
				//

				// Keep a record of our peers b/c we will need it later.
				impl.peers[p.ID] = &p

				impl.logger.Debug("peer connected",
					slog.String("dto", "consensus"),
					slog.Any("peer_id", p.ID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	//
	// Developers Note:
	// When we load up the consensus repository, we want to have a background
	// goroutine running which will check for any broadcast messages and
	// automatically send the current blockchain hash to the specific peer
	// right away.
	//
	go func(ctx context.Context, impl *MajorityVoteConsensusRepoImpl) {
		for {
			// Developers Note:
			// https://github.com/libp2p/go-libp2p/blob/master/examples/pubsub/basic-chat-with-rendezvous/main.go#L121

			msg, err := impl.sub.Next(ctx)
			if err != nil {
				impl.logger.Error("Failed to receive message",
					slog.Any("error", err))
				continue
			}
			if msg != nil {
				peerID := msg.ReceivedFrom

				impl.mu.Lock()
				defer impl.mu.Unlock()

				dataBytes := []byte(impl.currentBlockchainHash)
				if err := impl.dtoProtocol.SendResponse(peerID, dataBytes); err != nil {
					impl.logger.Error("Failed to send response",
						slog.Any("error", err))
					continue
				}
			}
		}
	}(context.Background(), impl)

	return impl
}
func (impl *MajorityVoteConsensusRepoImpl) QueryLatestHashByConsensus(ctx context.Context) (string, error) {
	// Attach a `1 minute` timeout so if we don't acheive consensus within that
	// time limit then we will need to abandon this request.
	ctx, cancel := context.WithTimeout(ctx, 1*time.Minute)
	defer cancel()

	// Variable used to keep track of the networks response.
	voteResults := make(map[peer.ID]string, len(impl.peers))

	// Variable used to synchronize all the go routines running in
	// background outside of this function.
	var wg sync.WaitGroup

	// Variable used to lock / unlock access when the goroutines want to
	// perform writes to our output response.
	var reqmu sync.Mutex

	// Load up the number of workers our waitgroup will need to handle.
	numOfPeers := len(impl.peers)
	wg.Add(numOfPeers)

	// Create a channel to collect errors from goroutines.
	errCh := make(chan error, numOfPeers)

	for peerID := range impl.peers {
		if err := impl.dtoProtocol.SendRequest(peerID, []byte("")); err != nil {
			return "", err
		}

		go func(mu *sync.Mutex) {
			defer wg.Done() // We are done this background task.

			resp, err := impl.dtoProtocol.WaitAndReceiveResponse(ctx)
			if err != nil {
				if ctx.Err() == context.DeadlineExceeded {
					impl.logger.Warn("timeout occurred")
				} else {
					errCh <- err
				}
				return
			}
			// Deserialize the result.
			currentBlockchainHash := string(resp.Content)

			// Lock our votes and add our new vote from a peer.
			mu.Lock()
			voteResults[resp.FromPeerID] = currentBlockchainHash
			mu.Unlock()
		}(&reqmu)
	}

	// Block the current execution until all our goroutines finish.
	go func() {
		wg.Wait()
		close(errCh)
	}()

	// Check if any errors occurred.
	select {
	case err := <-errCh:
		return "", err
	case <-errCh:
		// No errors occurred.
	}

	votes := make(map[string]int, 0)

	for _, hash := range voteResults {
		votes[hash]++
	}

	// Find the most common hash
	var maxCount int
	var mostCommonHash string
	for hash, count := range votes {
		if count > maxCount {
			maxCount = count
			mostCommonHash = hash
		}
	}

	return mostCommonHash, nil
}

// SetCurrentBlockchainHash method used to set the blockchain hash that you
// want to auto-submit to the peer-to-peer network from your peer when the
// network requests a consensus on the state of the blockchain pertaining to
// what is the hash of the latest block.
func (impl *MajorityVoteConsensusRepoImpl) CastVoteForLatestHashConsensus(newHash string) error {
	impl.mu.Lock()
	defer impl.mu.Unlock()
	impl.currentBlockchainHash = newHash
	return nil
}
