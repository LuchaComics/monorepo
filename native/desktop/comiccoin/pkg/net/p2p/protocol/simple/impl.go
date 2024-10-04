package simple

import (
	"bufio"
	"context"
	"encoding/binary"
	"io"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func NewSimpleDTOProtocol(logger *slog.Logger, host host.Host, reqProtocolID protocol.ID, respProtocolID protocol.ID) SimpleDTOProtocol {
	req := reqProtocolID
	resp := respProtocolID
	impl := &simpleProtocolImpl{
		logger:                      logger,
		host:                        host,
		requestChan:                 make(chan *SimpleDTORequest),
		responseChan:                make(chan *SimpleDTOResponse),
		protocolIDSimpleDTORequest:  req,
		protocolIDSimpleDTOResponse: resp,
	}
	host.SetStreamHandler(req, impl.onRemoteRequest)
	host.SetStreamHandler(resp, impl.onRemoteResponse)
	return impl
}

// remote peer requests handler
func (impl *simpleProtocolImpl) onRemoteRequest(s network.Stream) {
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
	req := &SimpleDTORequest{
		Content:    payloadBytes,
		FromPeerID: s.Conn().RemotePeer(),
	}

	//
	// STEP 5
	//

	impl.requestChan <- req
}

// remote Simple response handler
func (impl *simpleProtocolImpl) onRemoteResponse(s network.Stream) {
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

	resp := &SimpleDTOResponse{
		Content:    payloadBytes,
		FromPeerID: s.Conn().RemotePeer(), // Keep track of whom we received this message from.
	}

	//
	// STEP 5
	//

	impl.responseChan <- resp
}

// local sends to remote
func (impl *simpleProtocolImpl) SendRequest(peerID peer.ID, content []byte) error {
	//
	// STEP 1
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDSimpleDTORequest)
	if err != nil {
		impl.logger.Error("SendRequest: newstream error",
			slog.Any("error", err))
		return err
	}
	defer s.Close()

	// DEVELOPERS NOTE:
	// It's OK if `content` is empty. Do not handle any defensive coding as
	// there might be cases in which you want to send a request without any
	// payload.
	if content == nil {
		content = []byte(string(""))
	}

	//
	// STEP 4
	// First stream the length of the message to the peer
	//

	var lengthBuffer [4]byte
	binary.LittleEndian.PutUint32(lengthBuffer[:], uint32(len(content)))
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

	_, err = s.Write(content)
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message payload",
			slog.Any("error", err))
		return err
	}

	return nil
}

// local sends to remote
func (impl *simpleProtocolImpl) SendResponse(peerID peer.ID, content []byte) error {
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

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDSimpleDTOResponse)
	if err != nil {
		impl.logger.Error("SendResponse: failed to open stream",
			slog.Any("error", err))
		return err
	}
	defer s.Close()

	//
	// STEP 2
	// First stream the length of the message to the peer
	//

	var lengthBuffer [4]byte
	binary.LittleEndian.PutUint32(lengthBuffer[:], uint32(len(content)))
	_, err = s.Write(lengthBuffer[:])
	if err != nil {
		impl.logger.Error("SendRequest: failed to stream message header",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(content)
	if err != nil {
		impl.logger.Error("SendResponse: failed to stream message payload",
			slog.Any("error", err))
		return err
	}

	return err
}

func (impl *simpleProtocolImpl) WaitAndReceiveRequest(ctx context.Context) (*SimpleDTORequest, error) {
	for {
		select {
		case req := <-impl.requestChan:
			return req, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

func (impl *simpleProtocolImpl) WaitAndReceiveResponse(ctx context.Context) (*SimpleDTOResponse, error) {
	for {
		select {
		case resp := <-impl.responseChan:
			return resp, nil
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}
