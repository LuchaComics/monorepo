package domain

import (
	"context"

	"github.com/libp2p/go-libp2p/core/peer"
)

// ConsensusRepository is an interface that defines the methods for:
// Individual nodes to query the distributed / peer-to-peer network to get
// the latest hash of the ComicCoin blockchain
//
// What is consensus?
// Consensus is the process by which nodes in a distributed system (such as a
// peer-to-peer network) agree on a single source of truth or state, even when
// some nodes may fail, behave maliciously, or hold incorrect data.
type ConsensusRepository interface {
	// BroadcastRequestToNetwork function used to broadcast to the peer-to-peer
	// network that this node requests a consensus on the blockchain.
	BroadcastRequestToNetwork(ctx context.Context) error

	// ReceiveRequestFromNetwork function receives any consensus request from
	// the peer-to-peer network and returns the `peerID` that called upon the
	// network for the consensus.
	ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, error)

	SendResponseToPeer(ctx context.Context, peerID peer.ID, blockchainHash string) error

	ReceiveMajorityVoteConsensusResponseFromNetwork(ctx context.Context) (string, error)
}
