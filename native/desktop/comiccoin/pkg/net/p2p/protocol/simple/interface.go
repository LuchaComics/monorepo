package simple

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// SimpleDTOProtocol type
type simpleProtocolImpl struct {
	logger *slog.Logger

	host host.Host // local host
	mu   sync.Mutex

	requestChan  chan *SimpleDTORequest
	responseChan chan *SimpleDTOResponse

	protocolIDSimpleDTORequest  protocol.ID
	protocolIDSimpleDTOResponse protocol.ID
}

type SimpleDTOProtocol interface {
	SendRequest(peerID peer.ID, content []byte) error
	SendResponse(peerID peer.ID, content []byte) error

	WaitAndReceiveRequest(ctx context.Context) (*SimpleDTORequest, error)
	WaitAndReceiveResponse(ctx context.Context) (*SimpleDTOResponse, error)
}
