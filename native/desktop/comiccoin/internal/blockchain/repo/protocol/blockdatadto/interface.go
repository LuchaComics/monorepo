package blockdatadto

import (
	"context"
	"log/slog"
	"sync"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

// BlockDataDTOProtocol type
type blockDataDTOProtocolImpl struct {
	logger *slog.Logger

	host host.Host // local host
	mu   sync.Mutex

	requestChan  chan *BlockDataDTORequest
	responseChan chan *BlockDataDTOResponse

	protocolIDBlockDataDTORequest  protocol.ID
	protocolIDBlockDataDTOResponse protocol.ID
}

type BlockDataDTOProtocol interface {
	SendRequest(peerID peer.ID, blockDataHash string) error
	SendResponse(peerID peer.ID, blockData *domain.BlockDataDTO) error

	WaitAndReceiveRequest(ctx context.Context) (*BlockDataDTORequest, error)
	WaitAndReceiveResponse(ctx context.Context) (*BlockDataDTOResponse, error)
}
