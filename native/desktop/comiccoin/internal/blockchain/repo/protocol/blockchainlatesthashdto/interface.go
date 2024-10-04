package blockchainlatesthashdto

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

// BlockchainLatestHashDTOProtocol type
type blockchainLatestHashDTOProtocolImpl struct {
	logger *slog.Logger

	host host.Host // local host
	mu   sync.Mutex

	requestChan  chan *BlockchainLatestHashDTORequest
	responseChan chan *BlockchainLatestHashDTOResponse

	protocolIDBlockchainLatestHashDTORequest  protocol.ID
	protocolIDBlockchainLatestHashDTOResponse protocol.ID
}

type BlockchainLatestHashDTOProtocol interface {
	SendRequest(peerID peer.ID, content []byte) error
	SendResponse(peerID peer.ID, content []byte) error

	WaitAndReceiveRequest(ctx context.Context) (*BlockchainLatestHashDTORequest, error)
	WaitAndReceiveResponse(ctx context.Context) (*BlockchainLatestHashDTOResponse, error)
}
