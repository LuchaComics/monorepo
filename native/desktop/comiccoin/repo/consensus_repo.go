package repo

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/net/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/net/p2p/protocol/simple"
)

type ConsensusRepoImpl struct {
	config        *config.Config
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork
	dtoProtocol   simple.SimpleDTOProtocol

	rendezvousString string

	mu sync.Mutex

	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo
}

func NewConsensusRepoImpl(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) domain.ConsensusRepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain/consensus"

	//
	// STEP 1
	// Initialize our instance
	//

	impl := &ConsensusRepoImpl{
		config:           cfg,
		logger:           logger,
		libP2PNetwork:    libP2PNetwork,
		rendezvousString: rendezvousString,
		peers:            make(map[peer.ID]*peer.AddrInfo, 0),
	}

	//
	// STEP 2:
	// Create and advertise our `rendezvousString` which is essentially telling
	// our P2P network that clients can meet and communicate in our app at this
	// specific location.
	//

	// This is like your friend telling you the location to meet you.
	impl.libP2PNetwork.AdvertiseWithRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), impl.rendezvousString)

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
				slog.Any("local_peer_id", h.ID()),
				slog.Any("remote_peer_id", peerID),
				slog.String("dto", "consensus"),
			)
			delete(impl.peers, peerID)

		},
	})

	//
	dtoProtocol := simple.NewSimpleDTOProtocol(logger, h, "/consensus/req/1.0.0", "/consensus/resp/1.0.0")
	impl.dtoProtocol = dtoProtocol

	//
	// STEP 5:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect...")

		for {

			//
			// STEP 6:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), h, impl.rendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 7
				// Connect our peer with this peer and keep a record of it.
				//

				// Keep a record of our peers b/c we will need it later.
				impl.peers[p.ID] = &p

				impl.logger.Debug("peer connected",
					slog.String("dto", "consensus"),
					slog.Any("local_peer_id", h.ID()),
					slog.Any("remote_peer_id", p.ID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (impl *ConsensusRepoImpl) BroadcastRequestToNetwork(ctx context.Context) error {
	// Defensive code: Do not continue if we have no connections.
	if len(impl.peers) == 0 {
		return fmt.Errorf("error: %v", "no peers connected")
	}

	// impl.logger.Debug("consensus mechanism sending request to all peers...")
	for peerID := range impl.peers {
		// Note: Send empty request because we don't want anything.
		if err := impl.dtoProtocol.SendRequest(peerID, []byte("")); err != nil {
			impl.logger.Error("failed sent consensus request to peer",
				slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
				slog.Any("remote_peer_id", peerID),
				slog.Any("error", err))
			return err
		}
		// impl.logger.Debug("consensus mechanism sent request to peer",
		// 	slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
		// 	slog.Any("remote_peer_id", peerID))
	}
	// impl.logger.Debug("consensus mechanism finished sending request to all peers")
	return nil
}

func (impl *ConsensusRepoImpl) ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, error) {
	req, err := impl.dtoProtocol.WaitAndReceiveRequest(ctx)
	if err != nil {
		impl.logger.Error("failed receiving request from network",
			slog.Any("error", err))
		return "", err
	}
	// impl.logger.Debug("consensus mechanism received request from network",
	// 	slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
	// 	slog.Any("remote_peer_id", req.FromPeerID))
	return req.FromPeerID, nil
}

func (impl *ConsensusRepoImpl) SendResponseToPeer(ctx context.Context, peerID peer.ID, blockchainHash string) error {
	dataBytes := []byte(blockchainHash)
	if err := impl.dtoProtocol.SendResponse(peerID, dataBytes); err != nil {
		impl.logger.Error("failed sent consensus vote response to peer",
			slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
			slog.Any("remote_peer_id", peerID),
			slog.Any("error", err))
		return err
	}
	// impl.logger.Debug("consensus mechanism sent response to peer",
	// 	slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
	// 	slog.Any("remote_peer_id", peerID))
	return nil

}

func (impl *ConsensusRepoImpl) ReceiveIndividualResponseFromNetwork(ctx context.Context) (string, error) {
	resp, err := impl.dtoProtocol.WaitAndReceiveResponse(ctx)
	if err != nil {
		impl.logger.Error("failed receiving individiual consensus response from network",
			slog.Any("error", err))
		return "", err
	}

	hash := string(resp.Content)
	// impl.logger.Debug("consensus mechanism received response from network",
	// 	slog.String("hash", hash),
	// 	slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
	// 	slog.Any("remote_peer_id", resp.FromPeerID))
	return hash, nil
}

func (impl *ConsensusRepoImpl) ReceiveMajorityVoteConsensusResponseFromNetwork(ctx context.Context) (string, error) {
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

	for range impl.peers {
		go func(mu *sync.Mutex) {
			defer wg.Done() // We are done this background task.

			resp, err := impl.dtoProtocol.WaitAndReceiveResponse(ctx)
			if err != nil {
				impl.logger.Error("failed receiving consensus response from network",
					slog.Any("error", err))
				return
			}

			currentBlockchainHash := string(resp.Content)

			// If the response is not empty then save it.
			if currentBlockchainHash != "" {
				impl.logger.Debug("consensus mechanism received response from peer",
					slog.String("hash", currentBlockchainHash),
					slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()),
					slog.Any("remote_peer_id", resp.FromPeerID))

				// Lock our votes and add our new vote from a peer.
				mu.Lock()
				voteResults[resp.FromPeerID] = currentBlockchainHash
				mu.Unlock()
			}

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
		if err != nil {
			impl.logger.Error("network connectivity issue",
				slog.Any("error", err),
				slog.Any("local_peer_id", impl.libP2PNetwork.GetHost().ID()))
			return "", err
		}
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

	// impl.logger.Debug("consensus returned",
	// 	slog.Any("votes", votes),
	// 	slog.String("most_common_hash", mostCommonHash))

	return mostCommonHash, nil
}
