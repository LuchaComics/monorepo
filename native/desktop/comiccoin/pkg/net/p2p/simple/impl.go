package simple

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/protocol"
)

func NewSimpleMessageProtocol(host host.Host, protocolIDSimpleMessageRequest protocol.ID, protocolIDSimpleMessageResponse protocol.ID) SimpleMessageProtocol {
	impl := &simpleMessageProtocolImpl{
		host:                            host,
		requestChan:                     make(chan *SimpleMessageRequest),
		responseChan:                    make(chan *SimpleMessageResponse),
		protocolIDSimpleMessageRequest:  protocolIDSimpleMessageRequest,
		protocolIDSimpleMessageResponse: protocolIDSimpleMessageResponse,
	}
	host.SetStreamHandler(protocolIDSimpleMessageRequest, impl.onRemoteRequest)
	host.SetStreamHandler(protocolIDSimpleMessageResponse, impl.onRemoteResponse)
	return impl
}

// remote peer requests handler
func (impl *simpleMessageProtocolImpl) onRemoteRequest(s network.Stream) {
	log.Println("onRemoteRequest: received...")

	//
	// STEP 1
	//

	buf := bufio.NewReader(s)
	header, err := buf.ReadByte()
	if err != nil {
		s.Reset() // Important - don't forget!
		log.Printf("onRemoteRequest: failed to read message header: %v\n", err)
		return
	}

	log.Printf("onRemoteRequest: header: %v\n", header)

	//
	// STEP 2
	//

	payloadBytes := make([]byte, header)
	n, err := io.ReadFull(buf, payloadBytes)
	log.Printf("onRemoteRequest: payload has %d bytes\n", n)
	if err != nil {
		s.Reset() // Important - don't forget!
		log.Printf("onRemoteRequest: failed to read message payload: %v\n", err)
		return
	}

	log.Printf("onRemoteRequest: payload: %v\n", payloadBytes)

	//
	// STEP 3
	// Important step, since we finished our stream, then we need to close it.
	//

	s.Close()

	//
	// STEP 4
	//

	// Begin the deserialization
	req, err := NewSimpleMessageRequestFromDeserialize(payloadBytes)
	if err != nil {
		s.Reset()
		log.Printf("failed to deserialize remote request: %v\n", err)
		return
	}

	// Keep track of whom we received this message from.
	req.FromPeerID = s.Conn().RemotePeer()

	log.Printf("onRemoteRequest: payload deserialized: %v\n", req)

	//
	// STEP 5
	//

	impl.requestChan <- req
}

// remote Simple response handler
func (impl *simpleMessageProtocolImpl) onRemoteResponse(s network.Stream) {
	log.Println("onRemoteResponse: received...")

	//
	// STEP 1
	//

	buf := bufio.NewReader(s)
	header, err := buf.ReadByte()
	if err != nil {
		s.Reset() // Important - don't forget!
		log.Printf("onRemoteResponse: failed to read message header: %v\n", err)
		return
	}

	//
	// STEP 2
	//

	payloadBytes := make([]byte, header)
	n, err := io.ReadFull(buf, payloadBytes)
	log.Printf("payload has %d bytes", n)
	if err != nil {
		s.Reset() // Important - don't forget!
		log.Printf("onRemoteResponse: failed to read message payload: %v\n", err)
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
	resp, err := NewSimpleMessageResponseFromDeserialize(payloadBytes)
	if err != nil {
		s.Reset()
		log.Printf("onRemoteResponse: failed to deserialize remote request: %v\n", err)
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
func (impl *simpleMessageProtocolImpl) SendRequest(peerID peer.ID, content []byte) (string, error) {
	log.Printf("%s: Sending Simple to: %s....", impl.host.ID(), peerID)

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
	req := &SimpleMessageRequest{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Content: content,
	}

	//
	// STEP 2
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDSimpleMessageRequest)
	if err != nil {
		log.Printf("SendRequest NewStream error: %v\n", err)
		return "", err
	}
	defer s.Close()

	//
	// STEP 3
	//

	payloadBytes, err := req.Serialize()
	if err != nil {
		log.Printf("SendRequest Serialize error: %v\n", err)
		return "", err
	}

	//
	// STEP 4
	// First stream the length of the message to the peer
	//

	header := []byte{byte(len(payloadBytes))}
	_, err = s.Write(header)
	if err != nil {
		log.Printf("SendRequest: failed to stream message header: %v", err)
		return "", err
	}

	//
	// STEP 5
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(payloadBytes)
	if err != nil {
		log.Printf("SendRequest: failed to stream message payload: %v", err)
		return "", err
	}

	return req.ID, nil
}

// local sends to remote
func (impl *simpleMessageProtocolImpl) SendResponse(peerID peer.ID, content []byte) (string, error) {
	log.Printf("%s: Sending Simple to: %s....", impl.host.ID(), peerID)

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
	resp := &SimpleMessageResponse{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Content: content,
	}

	//
	// STEP 2
	//

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDSimpleMessageResponse)
	if err != nil {
		log.Printf("SendResponse: to open new stream: %v", err)
		return "", err
	}
	defer s.Close()

	//
	// STEP 3
	//

	payloadBytes, err := resp.Serialize()
	if err != nil {
		log.Printf("SendResponse: failed to serialize: %v", err)
		return "", err
	}

	//
	// STEP 4
	// First stream the length of the message to the peer
	//

	header := []byte{byte(len(payloadBytes))}
	_, err = s.Write(header)
	if err != nil {
		log.Printf("SendResponse: failed to stream message header: %v", err)
		return "", err
	}

	//
	// STEP 5
	// Lastely stream the payload of the message to the peer.
	//

	_, err = s.Write(payloadBytes)
	if err != nil {
		log.Printf("SendResponse: failed to stream message payload: %v", err)
		return "", err
	}

	return resp.ID, err
}

func (impl *simpleMessageProtocolImpl) WaitAndReceiveRequest(ctx context.Context) (*SimpleMessageRequest, error) {
	log.Println("WaitForAnyRequests: starting...")
	for {
		select {
		case req := <-impl.requestChan:
			log.Println("WaitForAnyRequests: received request...")
			return req, nil
		case <-ctx.Done():
			log.Println("WaitForAnyRequests: context done")
			return nil, ctx.Err()
		}
	}
}

func (impl *simpleMessageProtocolImpl) WaitAndReceiveResponse(ctx context.Context) (*SimpleMessageResponse, error) {
	log.Println("WaitForAnyResponses: starting...")
	for {
		select {
		case resp := <-impl.responseChan:
			log.Println("WaitForAnyResponses: received response...")
			return resp, nil
		case <-ctx.Done():
			log.Println("WaitForAnyResponses: context done")
			return nil, ctx.Err()
		}
	}
}
