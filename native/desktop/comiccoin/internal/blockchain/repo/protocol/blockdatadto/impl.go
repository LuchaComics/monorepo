package blockdatadto

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"

	"github.com/fxamacker/cbor/v2"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
)

func NewBlockDataDTOProtocol(logger *slog.Logger, host host.Host) BlockDataDTOProtocol {
	req := protocol.ID("/blockdatadto/req/1.0.0")
	resp := protocol.ID("/blockdatadto/resp/1.0.0")
	impl := &blockDataDTOProtocolImpl{
		logger:                         logger,
		host:                           host,
		requestChan:                    make(chan *BlockDataDTORequest),
		responseChan:                   make(chan *BlockDataDTOResponse),
		protocolIDBlockDataDTORequest:  req,
		protocolIDBlockDataDTOResponse: resp,
	}
	host.SetStreamHandler(req, impl.onRemoteRequest)
	host.SetStreamHandler(resp, impl.onRemoteResponse)
	return impl
}

// remote peer requests handler
func (impl *blockDataDTOProtocolImpl) onRemoteRequest(s network.Stream) {
	//
	// STEP 1
	//

	buf := bufio.NewReader(s)
	var lengthBuffer [4]byte
	_, err := io.ReadFull(buf, lengthBuffer[:])
	if err != nil {
		s.Reset() // Important - don't forget!
		impl.logger.Error("onRemoteResponse: failed to read message header",
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	payloadLength := int(binary.LittleEndian.Uint32(lengthBuffer[:]))

	//
	// STEP 2
	//

	payloadBytes := make([]byte, payloadLength)

	n, err := io.ReadFull(buf, payloadBytes)
	if err != nil {
		s.Reset() // Important - don't forget!
		impl.logger.Error("onRemoteRequest: failed to read message payload",
			slog.Any("payload_bytes_length", n),
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	//
	// STEP 3
	// Important step, since we finished our stream, then we need to close it.
	//

	s.Close()

	//
	// STEP 4
	//

	// Begin the deserialization
	req, err := NewBlockDataDTORequestFromDeserialize(payloadBytes)
	if err != nil {
		s.Reset()
		impl.logger.Error("onRemoteRequest: failed to deserialize remote request",
			slog.Any("payload", string(payloadBytes)),
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	// Keep track of whom we received this message from.
	req.FromPeerID = s.Conn().RemotePeer()

	//
	// STEP 5
	//

	impl.requestChan <- req
}

// remote Simple response handler
func (impl *blockDataDTOProtocolImpl) onRemoteResponse(s network.Stream) {
	//
	// STEP 1
	//

	buf := bufio.NewReader(s)

	var lengthBuffer [4]byte
	_, err := io.ReadFull(buf, lengthBuffer[:])
	if err != nil {
		s.Reset() // Important - don't forget!
		impl.logger.Error("onRemoteResponse: failed to read message header",
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	payloadLength := int(binary.LittleEndian.Uint32(lengthBuffer[:]))

	//
	// STEP 2
	//

	payloadBytes := make([]byte, payloadLength)

	n, err := io.ReadFull(buf, payloadBytes)
	if err != nil {
		s.Reset() // Important - don't forget!
		impl.logger.Error("onRemoteResponse: failed to read message payload",
			slog.Any("payload_bytes_length", n),
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	//
	// STEP 3
	// Important step, since we finished our stream, then we need to close it.
	//

	s.Close()

	//
	// STEP 4
	//

	// Variable we will use to return.
	dto := &domain.BlockDataDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if payloadBytes == nil {
		return
	}

	if err := cbor.Unmarshal(payloadBytes, &dto); err != nil {
		fmt.Printf("failed to deserialize stream message dto: %v\n", err)
		return
	}

	resp := &BlockDataDTOResponse{
		Payload: dto,
	}

	//
	// STEP 5
	//

	impl.responseChan <- resp
}

// local sends to remote
func (impl *blockDataDTOProtocolImpl) SendRequest(peerID peer.ID, blockDataHash string) error {
	//
	// STEP 1
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDBlockDataDTORequest)
	if err != nil {
		impl.logger.Error("SendRequest: newstream error",
			slog.Any("error", err))
		return err
	}
	defer s.Close()

	//
	// STEP 2
	//

	// create message data
	req := &BlockDataDTORequest{
		FromPeerID: s.Conn().LocalPeer(),
		ParamHash:  blockDataHash,
	}

	//
	// STEP 3
	//

	payloadBytes, err := req.Serialize()
	if err != nil {
		impl.logger.Error("SendRequest: failed to Serialize",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 4
	// First stream the length of the message to the peer
	//

	var lengthBuffer [4]byte
	binary.LittleEndian.PutUint32(lengthBuffer[:], uint32(len(payloadBytes)))
	_, err = s.Write(lengthBuffer[:])
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message header",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 5
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(payloadBytes)
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message payload",
			slog.Any("error", err))
		return err
	}

	return nil
}

// local sends to remote
func (impl *blockDataDTOProtocolImpl) SendResponse(peerID peer.ID, payload *domain.BlockDataDTO) error {
	//
	// STEP 1
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDBlockDataDTOResponse)
	if err != nil {
		impl.logger.Error("SendResponse: failed to open stream",
			slog.Any("error", err))
		return err
	}
	defer s.Close()

	//
	// STEP 2
	//

	payloadBytes, err := cbor.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to serialize stream message dto: %v", err)
	}

	//
	// STEP 4
	// First stream the length of the message to the peer
	//

	var lengthBuffer [4]byte
	binary.LittleEndian.PutUint32(lengthBuffer[:], uint32(len(payloadBytes)))
	_, err = s.Write(lengthBuffer[:])
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message header",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 5
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(payloadBytes)
	if err != nil {
		impl.logger.Error("SendResponse: failed to stream message payload",
			slog.Any("error", err))
		return err
	}

	return err
}

func (impl *blockDataDTOProtocolImpl) WaitAndReceiveRequest(ctx context.Context) (*BlockDataDTORequest, error) {
	for {
		select {
		case req := <-impl.requestChan:
			return req, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (impl *blockDataDTOProtocolImpl) WaitAndReceiveResponse(ctx context.Context) (*BlockDataDTOResponse, error) {
	for {
		select {
		case resp := <-impl.responseChan:
			return resp, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
