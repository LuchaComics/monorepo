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
		requests:                        make(map[peer.ID][]*SimpleMessageRequest),
		responses:                       make(map[peer.ID][]*SimpleMessageResponse),
		protocolIDSimpleMessageRequest:  protocolIDSimpleMessageRequest,
		protocolIDSimpleMessageResponse: protocolIDSimpleMessageResponse,
	}
	host.SetStreamHandler(protocolIDSimpleMessageRequest, impl.onRemoteRequest)
	host.SetStreamHandler(protocolIDSimpleMessageResponse, impl.onRemoteResponse)
	return impl
}

// remote peer requests handler
func (impl *simpleMessageProtocolImpl) onRemoteRequest(s network.Stream) {
	// get request data
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	req, err := NewSimpleMessageRequestFromDeserialize(buf)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	impl.mu.Lock()
	arr := impl.requests[s.Conn().RemotePeer()]
	arr = append(arr, req)
	impl.requests[s.Conn().RemotePeer()] = arr
	impl.mu.Unlock()
}

// remote Simple response handler
func (impl *simpleMessageProtocolImpl) onRemoteResponse(s network.Stream) {
	// get request data
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	resp, err := NewSimpleMessageResponseFromDeserialize(buf)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	impl.mu.Lock()
	arr := impl.responses[s.Conn().RemotePeer()]
	arr = append(arr, resp)
	impl.responses[s.Conn().RemotePeer()] = arr
	impl.mu.Unlock()
}

// local sends to remote
func (impl *simpleMessageProtocolImpl) SendRequest(peerID peer.ID, content []byte) (string, error) {
	log.Printf("%s: Sending Simple to: %s....", impl.host.ID(), peerID)

	// create message data
	req := &SimpleMessageRequest{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Content: content,
		From:    impl.host.ID(),
		To:      peerID,
	}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if string(content) == "" {
		return "", nil
	}

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDSimpleMessageRequest)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer s.Close()

	buf := bufio.NewWriter(s)

	reqBytes, err := req.Serialize()
	if err != nil {
		return "", err
	}

	bytesLen, err := buf.WriteString(fmt.Sprintf("%s\n", reqBytes))
	if err != nil {
		return "", err
	}

	fmt.Println("sent:", bytesLen)

	return req.ID, nil
}

// local sends to remote
func (impl *simpleMessageProtocolImpl) SendResponse(peerID peer.ID, content []byte) (string, error) {
	log.Printf("%s: Sending Simple to: %s....", impl.host.ID(), peerID)

	// create message data
	resp := &SimpleMessageResponse{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Content: content,
		From:    impl.host.ID(),
		To:      peerID,
	}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if string(content) == "" {
		return "", nil
	}

	s, err := impl.host.NewStream(context.Background(), peerID, impl.protocolIDSimpleMessageResponse)
	if err != nil {
		log.Println(err)
		return "", err
	}
	defer s.Close()

	buf := bufio.NewWriter(s)

	respBytes, err := resp.Serialize()
	if err != nil {
		return "", err
	}

	bytesLen, err := buf.WriteString(fmt.Sprintf("%s\n", respBytes))
	if err != nil {
		return "", err
	}

	fmt.Println("sent:", bytesLen)

	return resp.ID, err
}

func (impl *simpleMessageProtocolImpl) WaitForAnyRequests(ctx context.Context) (map[peer.ID][]*SimpleMessageRequest, error) {
	for {
		impl.mu.Lock()
		defer impl.mu.Unlock()
		if len(impl.requests) > 0 {
			reqCopy := copyMapDeepRequests(impl.requests)
			impl.requests = make(map[peer.ID][]*SimpleMessageRequest)
			return reqCopy, nil
		}
		time.Sleep(1 * time.Second)
	}
}

func (impl *simpleMessageProtocolImpl) WaitForAnyResponses(ctx context.Context) (map[peer.ID][]*SimpleMessageResponse, error) {
	for {
		impl.mu.Lock()
		defer impl.mu.Unlock()
		if len(impl.responses) > 0 {
			reqCopy := copyMapDeepResponses(impl.responses)
			impl.responses = make(map[peer.ID][]*SimpleMessageResponse)
			return reqCopy, nil
		}
		time.Sleep(1 * time.Second)
	}
}

//	func (impl *simpleMessageProtocolImpl) WaitForResponse() map[peer.ID][]*SimpleMessageResponse {
//		for {
//			impl.mu.Lock()
//			defer impl.mu.Unlock()
//		}
//
//		return impl.responses
//	}
//
//	func (impl *simpleMessageProtocolImpl) ReceiveRequests() map[peer.ID][]*SimpleMessageRequest {
//		impl.mu.Lock()
//		defer impl.mu.Unlock()
//
//		return impl.requests
//	}
//
//	func (impl *simpleMessageProtocolImpl) WaitForResponse(responseID string) (*SimpleMessageResponse, error) {
//		for {
//			reponses := impl.GetResponses()
//			resp, ok := reponses[responseID]
//			if !ok {
//				time.Sleep(5 * time.Second)
//				continue
//			}
//			if resp != nil {
//				return resp, nil
//			}
//		}
//	}
//
//	func (impl *simpleMessageProtocolImpl) WaitForRequest(requestID string) (*SimpleMessageRequest, error) {
//		for {
//			requests := impl.GetRequests()
//			req, ok := requests[requestID]
//			if !ok {
//				time.Sleep(5 * time.Second)
//				continue
//			}
//			if req != nil {
//				return req, nil
//			}
//		}
//	}
