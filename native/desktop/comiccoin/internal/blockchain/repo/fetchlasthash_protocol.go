package repo

import (
	"bufio"
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"io"
	"log"
	"sync"
	"time"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
)

// Node type - a p2p host implementing one or more p2p protocols
// https://github.com/libp2p/go-libp2p/blob/master/examples/multipro/node.go#L22
type Node struct {
	host.Host              // lib-p2p host
	*FetchLastHashProtocol // fetchLastHash protocol impl
	// add other protocols here...
}

// Create a new node with its implemented protocols
func NewNode(host host.Host, done chan bool) *Node {
	node := &Node{Host: host}
	node.FetchLastHashProtocol = NewFetchLastHashProtocol(node, done)
	return node
}

// pattern: /protocol-name/request-or-response-message/version
const ProtocolIDFetchLastHashRequest = "/fetchlasthash/fetchlasthashreq/0.0.1"
const ProtocolIDFetchLastHashResponse = "/fetchlasthash/fetchlasthashres/0.0.1"

// FetchLastHashProtocol type
type FetchLastHashProtocol struct {
	node      *Node // local host
	mu        sync.Mutex
	requests  map[string]*FetchLastHashRequest
	responses map[string]*FetchLastHashResponse
	done      chan bool
}

type FetchLastHashRequest struct {
	ID      string
	Hash    string
	Message string
}

type FetchLastHashResponse struct {
	ID      string
	Hash    string
	Message string
}

func (b *FetchLastHashRequest) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return result.Bytes(), nil
}

func NewFetchLastHashRequestFromDeserialize(data []byte) (*FetchLastHashRequest, error) {
	// Variable we will use to return.
	dto := &FetchLastHashRequest{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}

func (b *FetchLastHashResponse) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return result.Bytes(), nil
}

func NewFetchLastHashResponseFromDeserialize(data []byte) (*FetchLastHashResponse, error) {
	// Variable we will use to return.
	dto := &FetchLastHashResponse{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}

func NewFetchLastHashProtocol(node *Node, done chan bool) *FetchLastHashProtocol {
	p := &FetchLastHashProtocol{
		node:     node,
		requests: make(map[string]*FetchLastHashRequest),
		done:     done,
	}
	node.SetStreamHandler(ProtocolIDFetchLastHashRequest, p.onRemoteRequest)
	node.SetStreamHandler(ProtocolIDFetchLastHashResponse, p.onRemoteResponse)
	return p
}

// remote peer requests handler
func (p *FetchLastHashProtocol) onRemoteRequest(s network.Stream) {
	// get request data
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	req, err := NewFetchLastHashRequestFromDeserialize(buf)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	log.Printf("%s: Received fetch last hash request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), req.Message)

	p.mu.Lock()
	p.requests[req.ID] = req
	p.mu.Unlock()

	p.done <- true
}

// remote fetchLastHash response handler
func (p *FetchLastHashProtocol) onRemoteResponse(s network.Stream) {
	// get request data
	buf, err := io.ReadAll(s)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	s.Close()

	resp, err := NewFetchLastHashResponseFromDeserialize(buf)
	if err != nil {
		s.Reset()
		log.Println(err)
		return
	}
	log.Printf("%s: Received fetch last hash request from %s. Message: %s", s.Conn().LocalPeer(), s.Conn().RemotePeer(), resp.Message)

	p.mu.Lock()
	p.responses[resp.ID] = resp
	p.mu.Unlock()

	p.done <- true
}

func (p *FetchLastHashProtocol) GetResponse() map[string]*FetchLastHashResponse {
	p.mu.Lock()
	defer p.mu.Unlock()
	return p.responses
}

// local sends to remote
func (p *FetchLastHashProtocol) SendRequest(peerID peer.ID, hash string) (string, error) {
	log.Printf("%s: Sending fetchLastHash to: %s....", p.node.ID(), peerID)

	// create message data
	req := &FetchLastHashRequest{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Hash:    hash,
		Message: fmt.Sprintf("Fetch last hash from %s", p.node.ID()),
	}

	// store ref request so response handler has access to it
	p.mu.Lock()
	p.requests[req.ID] = req
	p.mu.Unlock()

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if hash == "" {
		return "", nil
	}

	s, err := p.node.Host.NewStream(context.Background(), peerID, ProtocolIDFetchLastHashRequest)
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
func (p *FetchLastHashProtocol) SendResponse(peerID peer.ID, hash string) bool {
	log.Printf("%s: Sending fetchLastHash to: %s....", p.node.ID(), peerID)

	// create message data
	resp := &FetchLastHashResponse{
		ID:      fmt.Sprintf("%v", time.Now().Unix()),
		Hash:    hash,
		Message: fmt.Sprintf("Fetch last hash from %s", p.node.ID()),
	}

	// store ref request so response handler has access to it
	p.mu.Lock()
	p.responses[resp.ID] = resp
	p.mu.Unlock()

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if hash == "" {
		return false
	}

	s, err := p.node.Host.NewStream(context.Background(), peerID, ProtocolIDFetchLastHashResponse)
	if err != nil {
		log.Println(err)
		return false
	}
	defer s.Close()

	buf := bufio.NewWriter(s)

	respBytes, err := resp.Serialize()
	if err != nil {
		return false
	}

	bytesLen, err := buf.WriteString(fmt.Sprintf("%s\n", respBytes))
	if err != nil {
		return false
	}

	fmt.Println("sent:", bytesLen)

	return true
}
