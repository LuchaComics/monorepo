package simple

import (
	"github.com/libp2p/go-libp2p/core/peer"
)

// SimpleDTORequest represents a request message sent between nodes in the network.
// It contains the payload content and the ID of the sending peer.
type SimpleDTORequest struct {
	// The content of the request message.
	Content []byte `json:"content"`

	// The ID of the peer that sent this request.
	// Note: This field is set by the receiving node, not the sender, when the request is processed.
	FromPeerID peer.ID `json:"from_peer_id"`
}

// SimpleDTOResponse represents a response message sent between nodes in the network.
// It contains the payload content and the ID of the sending peer.
type SimpleDTOResponse struct {
	// The content of the response message.
	Content []byte `json:"content"`

	// The ID of the peer that sent this response.
	// Note: This field is set by the receiving node, not the sender, when the response is processed.
	FromPeerID peer.ID `json:"from_peer_id"`
}
