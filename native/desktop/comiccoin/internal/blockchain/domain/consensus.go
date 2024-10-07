package domain

import (
	"context"
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
	// Method broadcasts a request for latest hash to every connected peer in
	// the network at the given time and waits for everyone to submit their
	// response. Upon receiving the responses, this function will return the
	// hash that was agreed upon by `consenus` of the network.
	QueryLatestHashByConsensus(ctx context.Context) (string, error)

	// // Method handles receiving network requests of nodes whom want to figure
	// // out `consensus` for their machine. Method is execution blocking and will
	// // unblock and return a result only when it received a request from the
	// // network.
	// ReceiveRequestFromNetwork(ctx context.Context) (peer.ID, error)
	//
	// // Method used by nodes to submit their latest hash for `consensus`.
	// SubmitQueryResponseToPeer(ctx context.Context, peerID peer.ID, latestHash string) error
}
