package blockchainlatesthashdto

import (
	"bufio"
	"context"
	"encoding/binary"
	"fmt"
	"io"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func NewBlockchainLatestHashDTOProtocol(logger *slog.Logger, host host.Host) BlockchainLatestHashDTOProtocol {
	req := protocol.ID("/lastblockdatahash/req/1.0.0")
	resp := protocol.ID("/lastblockdatahash/resp/1.0.0")
	impl := &blockchainLatestHashDTOProtocolImpl{
		logger:                                   logger,
		host:                                     host,
		requestChan:                              make(chan *BlockchainLatestHashDTORequest),
		responseChan:                             make(chan *BlockchainLatestHashDTOResponse),
		protocolIDBlockchainLatestHashDTORequest: req,
		protocolIDBlockchainLatestHashDTOResponse: resp,
	}
	host.SetStreamHandler(req, impl.onRemoteRequest)
	host.SetStreamHandler(resp, impl.onRemoteResponse)
	return impl
}

// remote peer requests handler
func (impl *blockchainLatestHashDTOProtocolImpl) onRemoteRequest(s network.Stream) {
	//
	// STEP 1
	//

	buf := bufio.NewReader(s)
	header, err := buf.ReadByte()
	if err != nil {
		s.Reset() // Important - don't forget!
		impl.logger.Error("onRemoteRequest: failed to read message header",
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	//
	// STEP 2
	//

	payloadBytes := make([]byte, header)

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
	req, err := NewBlockchainLatestHashDTORequestFromDeserialize(payloadBytes)
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
func (impl *blockchainLatestHashDTOProtocolImpl) onRemoteResponse(s network.Stream) {
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

	// Begin the deserialization
	resp, err := NewBlockchainLatestHashDTOResponseFromDeserialize(payloadBytes)
	if err != nil {
		s.Reset()
		impl.logger.Error("onRemoteResponse: failed to deserialize remote request",
			slog.Any("payload", string(payloadBytes)),
			slog.Any("peer_id", s.Conn().RemotePeer()),
			slog.Any("error", err))
		return
	}

	// Keep track of whom we received this message from.
	resp.FromPeerID = s.Conn().RemotePeer()

	//
	// STEP 5
	//

	impl.responseChan <- resp
}

// local sends to remote
func (impl *blockchainLatestHashDTOProtocolImpl) SendRequest(peerID peer.ID, content []byte) (string, error) {
	// DEVELOPERS NOTE:
	// It's OK if `content` is empty. Do not handle any defensive coding as
	// there might be cases in which you want to send a request without any
	// payload.
	if content == nil {
		content = []byte(string(""))
	}

	//
	// STEP 1
	//

	// create message data
	req := &BlockchainLatestHashDTORequest{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Content: content,
	}

	//
	// STEP 2
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDBlockchainLatestHashDTORequest)
	if err != nil {
		impl.logger.Error("SendRequest: newstream error",
			slog.Any("error", err))
		return "", err
	}
	defer s.Close()

	//
	// STEP 3
	//

	payloadBytes, err := req.Serialize()
	if err != nil {
		impl.logger.Error("SendRequest: failed to Serialize",
			slog.Any("error", err))
		return "", err
	}

	//
	// STEP 4
	// First stream the length of the message to the peer
	//

	header := []byte{byte(len(payloadBytes))}
	_, err = s.Write(header)
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message header",
			slog.Any("error", err))
		return "", err
	}

	//
	// STEP 5
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(payloadBytes)
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message payload",
			slog.Any("error", err))
		return "", err
	}

	return req.ID, nil
}

// local sends to remote
func (impl *blockchainLatestHashDTOProtocolImpl) SendResponse(peerID peer.ID, content []byte) (string, error) {
	// DEVELOPERS NOTE:
	// It's OK if `content` is empty. Do not handle any defensive coding as
	// there might be cases in which you want to send a request without any
	// payload.
	if content == nil {
		content = []byte(string(""))
	}

	//
	// STEP 1
	//

	// create message data
	resp := &BlockchainLatestHashDTOResponse{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Content: content,
	}

	//
	// STEP 2
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDBlockchainLatestHashDTOResponse)
	if err != nil {
		impl.logger.Error("SendResponse: failed to open stream",
			slog.Any("error", err))
		return "", err
	}
	defer s.Close()

	//
	// STEP 3
	//

	payloadBytes, err := resp.Serialize()
	if err != nil {
		impl.logger.Error("SendResponse: failed to serialize",
			slog.Any("error", err))
		return "", err
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
		return "", err
	}

	//
	// STEP 5
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(payloadBytes)
	if err != nil {
		impl.logger.Error("SendResponse: failed to stream message payload",
			slog.Any("error", err))
		return "", err
	}

	return resp.ID, err
}

func (impl *blockchainLatestHashDTOProtocolImpl) WaitAndReceiveRequest(ctx context.Context) (*BlockchainLatestHashDTORequest, error) {
	for {
		select {
		case req := <-impl.requestChan:
			return req, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (impl *blockchainLatestHashDTOProtocolImpl) WaitAndReceiveResponse(ctx context.Context) (*BlockchainLatestHashDTOResponse, error) {
	for {
		select {
		case resp := <-impl.responseChan:
			return resp, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
