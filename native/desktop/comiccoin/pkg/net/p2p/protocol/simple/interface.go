package simple

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// SimpleDTOProtocol is an implementation of the SimpleDTOProtocol interface.
// It provides methods for sending and receiving SimpleDTO messages between nodes in a libp2p network.
type simpleProtocolImpl struct {
	// The logger used for logging events and errors.
	logger *slog.Logger

	// The local host node.
	host host.Host

	// A mutex used to protect access to the protocol's internal state.
	mu sync.Mutex

	// Channels used to receive incoming requests and responses.
	requestChan  chan *SimpleDTORequest
	responseChan chan *SimpleDTOResponse

	// The protocol IDs used for sending and receiving SimpleDTO requests and responses.
	protocolIDSimpleDTORequest  protocol.ID
	protocolIDSimpleDTOResponse protocol.ID
}

// SimpleDTOProtocol defines the interface for sending and receiving SimpleDTO messages.
// It provides methods for sending requests and responses, as well as waiting for incoming requests and responses.
type SimpleDTOProtocol interface {
	// Send a SimpleDTO request to a peer.
	SendRequest(peerID peer.ID, content []byte) error

	// Send a SimpleDTO response to a peer.
	SendResponse(peerID peer.ID, content []byte) error

	// Wait for an incoming SimpleDTO request and return it.
	WaitAndReceiveRequest(ctx context.Context) (*SimpleDTORequest, error)

	// Wait for an incoming SimpleDTO response and return it.
	WaitAndReceiveResponse(ctx context.Context) (*SimpleDTOResponse, error)
}
