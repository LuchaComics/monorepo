package domain

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

// BlockchainLastestHashDTO represents the data transfer object for the latest hash of the blockchain.
// It is a string that holds the latest hash value.
type BlockchainLastestHashDTO string

// BlockchainLastestHashDTORepository is an interface that defines the methods for interacting with the latest hash of the blockchain.
// It provides methods for sending requests to peers, receiving requests from peers, and sending responses to peers.
type BlockchainLastestHashDTORepository interface {
	// SendRequestToRandomPeer sends a request to a random connected peer in the peer-to-peer network.
	// It takes a context and returns an error if one occurs.
	SendRequestToRandomPeer(ctx context.Context) error

	// ReceiveRequestFromNetwork receives a request from the peer-to-peer network.
	// It takes a context and returns the peer ID of the peer that sent the request and an error if one occurs.
	ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, error)

	// SendResponseToPeer sends a response to a peer that requested block data.
	// It takes a context, the peer ID of the peer to send the response to, and the latest hash data to send.
	// It returns an error if one occurs.
	SendResponseToPeer(ctx context.Context, peerID peer.ID, data BlockchainLastestHashDTO) error

	// ReceiveResponseFromNetwork receives a response from the peer-to-peer network.
	// It takes a context and returns the latest hash data and an error if one occurs.
	ReceiveResponseFromNetwork(ctx context.Context) (BlockchainLastestHashDTO, error)
}
