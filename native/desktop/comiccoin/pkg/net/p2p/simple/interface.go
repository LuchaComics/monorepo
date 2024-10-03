package simple

import (
	"context"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// SimpleMessageProtocol type
type simpleMessageProtocolImpl struct {
	host host.Host // local host
	mu   sync.Mutex

	requests map[peer.ID][]*SimpleMessageRequest

	responses map[peer.ID][]*SimpleMessageResponse

	protocolIDSimpleMessageRequest  protocol.ID
	protocolIDSimpleMessageResponse protocol.ID
}

type SimpleMessageProtocol interface {
	SendRequest(peerID peer.ID, content []byte) (string, error)
	SendResponse(peerID peer.ID, content []byte) (string, error)

	WaitForAnyResponses(ctx context.Context) (map[peer.ID][]*SimpleMessageResponse, error)
	WaitForAnyRequests(ctx context.Context) (map[peer.ID][]*SimpleMessageRequest, error)

	// ReceiveResponses() []*SimpleMessageResponse
	// ReceiveRequests() []*SimpleMessageRequest
	//
	// WaitForRequest(requestID string) (*SimpleMessageRequest, error)
	// WaitForResponse(responseID string) (*SimpleMessageResponse, error)
}
