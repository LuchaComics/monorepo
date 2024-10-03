package domain

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

type BlockchainLastestHashDTO string

type BlockchainLastestHashDTORepository interface {
	// Function will randomly pick a connected peer and send them a request.
	SendRequestToRandomPeer(ctx context.Context) error

	// Function will block your current execution and wait until it receives
	// any request from the peer-to-peer network. Function will return the
	// `peerID` that sent the request and the hash value.
	ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, error)

	// Function will send sync data to the peer whom requested block data.
	SendResponseToPeer(ctx context.Context, peerID peer.ID, data BlockchainLastestHashDTO) error

	ReceiveResponseFromNetwork(ctx context.Context) (BlockchainLastestHashDTO, error)
}
