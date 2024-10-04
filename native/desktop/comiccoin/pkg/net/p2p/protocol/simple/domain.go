package simple

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

type SimpleDTORequest struct {
	Content []byte `json:"content"`

	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`
}

type SimpleDTOResponse struct {
	Content []byte `json:"content"`

	// Value set by the receiving node, not the sender in the payload!
	FromPeerID peer.ID `json:"from_peer_id"`
}
