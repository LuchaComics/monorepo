package repo

import (
	"context"
	"log"
	"log/slog"
	"sync"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	p2p "github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/net/p2p/p2pmessagedto"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

type BlockchainSyncRepo struct {
	config             *config.Config
	logger             *slog.Logger
	libP2PNetwork      p2p.LibP2PNetwork
	p2pMessengeDTORepo p2pmessagedto.P2PMessageDTORepository

	mu              sync.Mutex
	requestsBuffer  []*p2pmessagedto.P2PMessageDTO
	responsesBuffer []*p2pmessagedto.P2PMessageDTO
}

type BlockchainSyncRepository interface {
	// Function will randomly pick a connected peer and send them a request.
	SendRequestToRandomPeer(ctx context.Context, hash string) error

	// Function will block your current execution and wait until it receives
	// any request from the peer-to-peer network. Function will return the
	// `peerID` that sent the request and the hash value.
	ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, string, error)

	// Function will send sync data to the peer whom requested block data.
	SendResponseToPeer(ctx context.Context, peerID peer.ID, data *domain.BlockDataDTO) error

	ReceiveResponseFromNetwork(ctx context.Context) (*domain.BlockDataDTO, error)
}

func (impl *BlockchainSyncRepo) handleResponse() {
	ctx := context.Background()
	for {
		log.Println("wait and receive")
		res, err := impl.p2pMessengeDTORepo.WaitAndReceiveFromNetwork(ctx)
		if err != nil {
			impl.logger.Error("failed getting from p2p messages", slog.Any("error", err))
			continue
		}
		if res != nil {
			log.Println("ok")
			impl.mu.Lock()
			defer impl.mu.Unlock()

			if res.Type == p2pmessagedto.P2PMessageDTOTypeRequest {
				impl.requestsBuffer = append(impl.requestsBuffer, res)
			}
			if res.Type == p2pmessagedto.P2PMessageDTOTypeResponse {
				impl.responsesBuffer = append(impl.responsesBuffer, res)
			}
		}
	}
}

func NewBlockchainSyncRepo(cfg *config.Config, logger *slog.Logger, libP2PNetwork p2p.LibP2PNetwork) BlockchainSyncRepository {
	rendezvousString := "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain/blockdatadto"
	protocolID := protocol.ID("/sync/1.0.0")

	responsesBuffer := make([]*p2pmessagedto.P2PMessageDTO, 0)
	requestsBuffer := make([]*p2pmessagedto.P2PMessageDTO, 0)
	p2pMessengeDTORepo := p2pmessagedto.NewP2PMessageDTORepo(logger, libP2PNetwork, rendezvousString, protocolID)

	impl := &BlockchainSyncRepo{
		config:             cfg,
		logger:             logger,
		libP2PNetwork:      libP2PNetwork,
		p2pMessengeDTORepo: p2pMessengeDTORepo,
		responsesBuffer:    responsesBuffer,
		requestsBuffer:     requestsBuffer}

	go func() {
		impl.handleResponse()
	}()

	return impl
}

func (impl *BlockchainSyncRepo) SendRequestToRandomPeer(ctx context.Context, hash string) error {
	dto := &p2pmessagedto.P2PMessageDTO{
		FunctionID: "SendRequestToRandomPeer",
		Type:       p2pmessagedto.P2PMessageDTOTypeRequest,
		Content:    []byte(hash),
	}
	return impl.p2pMessengeDTORepo.SendToRandomPeerInNetwork(ctx, dto)
}

func (impl *BlockchainSyncRepo) ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, string, error) {
	if len(impl.requestsBuffer) == 0 {
		return "", "", nil
	}
	// Get the first request from the buffer
	firstRequest := impl.requestsBuffer[0]
	// Delete the first request from the buffer
	impl.requestsBuffer = impl.requestsBuffer[1:]
	return firstRequest.PeerID, string(firstRequest.Content), nil
}

func (impl *BlockchainSyncRepo) SendResponseToPeer(ctx context.Context, peerID peer.ID, data *domain.BlockDataDTO) error {
	content, err := data.Serialize()
	if err != nil {
		return err
	}
	dto := &p2pmessagedto.P2PMessageDTO{
		FunctionID: "SendResponseToPeer",
		Type:       p2pmessagedto.P2PMessageDTOTypeResponse,
		Content:    content,
	}
	return impl.p2pMessengeDTORepo.SendToSpecificPeerInNetwork(ctx, peerID, dto)

}

func (impl *BlockchainSyncRepo) ReceiveResponseFromNetwork(ctx context.Context) (*domain.BlockDataDTO, error) {
	if len(impl.responsesBuffer) == 0 {
		return nil, nil
	}
	// Get the first request from the buffer
	firstResponse := impl.responsesBuffer[0]
	// Delete the first request from the buffer
	impl.responsesBuffer = impl.responsesBuffer[1:]

	data, err := domain.NewBlockDataDTOFromDeserialize(firstResponse.Content)
	if err != nil {
		return nil, err
	}
	return data, nil
}
