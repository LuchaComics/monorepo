package simple

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// SimpleMessageProtocol type
type simpleMessageProtocolImpl struct {
	logger *slog.Logger

	host host.Host // local host
	mu   sync.Mutex

	requestChan  chan *SimpleMessageRequest
	responseChan chan *SimpleMessageResponse

	protocolIDSimpleMessageRequest  protocol.ID
	protocolIDSimpleMessageResponse protocol.ID
}

type SimpleMessageProtocol interface {
	SendRequest(peerID peer.ID, content []byte) (string, error)
	SendResponse(peerID peer.ID, content []byte) (string, error)

	WaitAndReceiveRequest(ctx context.Context) (*SimpleMessageRequest, error)
	WaitAndReceiveResponse(ctx context.Context) (*SimpleMessageResponse, error)
}
