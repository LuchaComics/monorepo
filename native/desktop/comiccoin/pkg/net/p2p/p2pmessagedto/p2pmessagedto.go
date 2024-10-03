package p2pmessagedto

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"math/rand"
	"time"

	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

const (
	P2PMessageDTOTypeRequest  = 1
	P2PMessageDTOTypeResponse = 2
)

// P2PMessageDTO respents data-transfer object that is used for bi-directional
// communication between peers in our network. This is used primarily for
// direct peer to peer communication.
type P2PMessageDTO struct {
	// The user whom sent this message.
	PeerID peer.ID

	// You set whatever function id you want.
	FunctionID string

	Type    int
	Content []byte
}

func (b *P2PMessageDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return result.Bytes(), nil
}

func NewP2PMessageDTOFromDeserialize(data []byte) (*P2PMessageDTO, error) {
	// Variable we will use to return.
	dto := &P2PMessageDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}

type P2PMessageDTORepository interface {
	// Function will randomly pick a connected peer and send them a request.
	SendToRandomPeerInNetwork(ctx context.Context, dto *P2PMessageDTO) error

	// Function will connect to specific peer and send a request.
	SendToSpecificPeerInNetwork(ctx context.Context, peerID peer.ID, dto *P2PMessageDTO) error

	// Function will block the current thread it's in and unblock when a single message has been received from the network.
	WaitAndReceiveFromNetwork(ctx context.Context) (*P2PMessageDTO, error)
}

type P2PMessageDTORepo struct {
	logger        *slog.Logger
	libP2PNetwork p2p.LibP2PNetwork

	rendezvousString string
	protocolID       protocol.ID

	// The list of connected peer addresses
	peers map[peer.ID]*peer.AddrInfo

	// The list of connected peers with a direct stream with tem.
	streams map[peer.ID]network.Stream
}

func NewP2PMessageDTORepo(logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork, rendezvousString string, protocolID protocol.ID) *P2PMessageDTORepo {
	//
	// STEP 1
	// Initialize our instance
	//

	impl := &P2PMessageDTORepo{
		logger:           logger,
		libP2PNetwork:    libP2PNetwork,
		rendezvousString: rendezvousString,
		protocolID:       protocolID,
		peers:            make(map[peer.ID]*peer.AddrInfo, 0),
		streams:          make(map[peer.ID]network.Stream, 0),
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

	host := libP2PNetwork.GetHost()

	//
	// STEP 4:
	// In a peer-to-peer network we need to be away of when peers disconnect
	// our network, the following code will callback when a peer disconnects so
	// our repository can remove the peer from our records.
	//

	//Remove disconnected peer
	host.Network().Notify(&network.NotifyBundle{
		DisconnectedF: func(_ network.Network, c network.Conn) {
			peerID := c.RemotePeer()
			impl.logger.Warn("peer disconnected", slog.Any("peer_id", peerID))
			delete(impl.peers, peerID)

			impl.logger.Warn("stream closed", slog.Any("peer_id", peerID))
			stream, ok := impl.streams[peerID]
			if ok {
				stream.Close()
				delete(impl.streams, peerID)

			}
		},
	})

	// Set a function as stream handler. This function is called when a peer
	// initiates a connection and starts a stream with this peer.
	host.SetStreamHandler(impl.protocolID, func(stream network.Stream) {
		// Handle incoming streams
		switch stream.Protocol() {
		case impl.protocolID:
			impl.streams[host.ID()] = stream
		default:
			// Handle unknown protocol
			log.Fatalf("Unknown protocol: %v", stream.Protocol())
		}
	})

	//
	// STEP 5:
	// When our repository loads up, we need to create a background goroutine
	// which will wait for new connections and get our list of peers that
	// connect in real-time to our application for this particular repository.
	//

	go func() {

		impl.logger.Debug("waiting for peers to connect...",
			slog.Any("protocol_id", impl.protocolID))

		for {

			//
			// STEP 6:
			// Wait to connect with new peers.
			//

			impl.libP2PNetwork.DiscoverPeersAtRendezvousString(context.Background(), impl.libP2PNetwork.GetHost(), impl.rendezvousString, func(p peer.AddrInfo) error {

				//
				// STEP 7
				// Connect our peer with this peer and keep a record of it.
				//

				// Keep a record of our peers b/c we will need it later.
				impl.peers[p.ID] = &p

				ctx := context.Background()
				stream, err := host.NewStream(ctx, p.ID, protocol.ID(impl.protocolID))
				if err != nil {
					// logger.Warn("Connection failed", slog.Any("error", err))
					return err
				} else {
					impl.streams[p.ID] = stream
				}

				impl.logger.Debug("peer connected",
					slog.Any("peer_id", p.ID),
					slog.Any("protocol_id", impl.protocolID))

				// Return nil to indicate success (no errors occured).
				return nil
			})
		}
	}()

	return impl
}

func (r *P2PMessageDTORepo) SendToRandomPeerInNetwork(ctx context.Context, dto *P2PMessageDTO) error {
	randomPeerID := r.randomPeerID()
	if randomPeerID == "" {
		return nil
	}

	return r.SendToSpecificPeerInNetwork(ctx, randomPeerID, dto)
}

func (r *P2PMessageDTORepo) SendToSpecificPeerInNetwork(ctx context.Context, peerID peer.ID, dto *P2PMessageDTO) error {
	dto.PeerID = peerID
	stream, ok := r.streams[peerID]
	if !ok {
		r.logger.Debug("stream does not exist",
			slog.Any("peer_id", peerID))
		return fmt.Errorf("stream does not exist for peer_id: %v", peerID)
	}

	r.logger.Debug("random peer selected, making request now...",
		slog.Any("peer_id", peerID))

	dtoBytes, err := dto.Serialize()
	if err != nil {
		return err
	}

	// Append a newline character to the serialized DTO bytes
	dtoBytes = append(dtoBytes, '\n')

	buf := bufio.NewWriter(stream)

	bytesLen, err := buf.WriteString(fmt.Sprintf("%s\n", dtoBytes))
	if err != nil {
		return err
	}

	r.logger.Debug("sent",
		slog.Any("bytes", bytesLen),
		slog.Any("peer_id", peerID))

	return nil
}

func (r *P2PMessageDTORepo) WaitAndReceiveFromNetwork(ctx context.Context) (*P2PMessageDTO, error) {
	dtos := make(chan *P2PMessageDTO)
	errs := make(chan error)

	if len(r.streams) == 0 {
		log.Println("WaitAndReceiveFromNetwork: empty")
		return nil, nil
	}

	for peerID, stream := range r.streams {

		log.Println("WaitAndReceiveFromNetwork: start")
		go func(peerID peer.ID, stream network.Stream) {
			log.Println("WaitAndReceiveFromNetwork: go")
			for {
				dto, err := r.receiveDTOFromStream(stream)
				log.Println("dto:", dto, err)
				if err != nil {
					errs <- err
					return
				}
				log.Println("RCV:", dto)
				dtos <- dto
			}
		}(peerID, stream)
	}

	select {
	case dto := <-dtos:
		return dto, nil
	case err := <-errs:
		return nil, err
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

var ErrStreamClosed = errors.New("stream is closed")

func (r *P2PMessageDTORepo) receiveDTOFromStream(stream network.Stream) (*P2PMessageDTO, error) {
	buf := bufio.NewReader(stream)

	for {
		log.Println("xxxx")
		dtoStr, err := buf.ReadString('\n')
		if err == io.EOF {
			log.Println("--->", ErrStreamClosed)
			return nil, ErrStreamClosed
		}
		if err != nil {
			log.Println("err---> closed")
			return nil, err
		}

		log.Println("--->", string(dtoStr))

		if len(dtoStr) > 0 {
			dto, err := NewP2PMessageDTOFromDeserialize([]byte(dtoStr))
			if err != nil {
				return nil, err
			}

			return dto, nil
		}
	}
}

func (r *P2PMessageDTORepo) randomPeerID() peer.ID {
	// Seed the random number generator
	rand.Seed(time.Now().UnixNano())

	// Get a list of peer IDs
	peerIDs := make([]peer.ID, 0, len(r.peers))
	for id := range r.peers {
		peerIDs = append(peerIDs, id)
	}

	// Select a random peer ID
	if len(peerIDs) == 0 {
		// Handle the case where there are no peers
		return ""
	}
	return peerIDs[rand.Intn(len(peerIDs))]
}

func (r *P2PMessageDTORepo) getRandomStream() (network.Stream, error) {
	peerID := r.randomPeerID()
	if peerID == "" {
		return nil, errors.New("no valid peers available")
	}

	s, _ := r.streams[peerID]
	return s, nil
}

func (r *P2PMessageDTORepo) getRandomPeer() (*peer.AddrInfo, error) {
	peerID := r.randomPeerID()
	if peerID == "" {
		return nil, errors.New("no valid peers available")
	}

	// Connect to a random peer
	peer, _ := r.peers[peerID]
	if peer == nil {
		return nil, errors.New("no peers connected")
	}
	return peer, nil
}
