package p2p

import (
	"log/slog"

	"github.com/libp2p/go-libp2p/core/network"
)

const fetchProtocolVersion = "/fetch/1.0.0"

type FetchProtocol struct {
	node *nodeInputPort
}

func NewFetchProtocol(node *nodeInputPort, stream network.Stream) *FetchProtocol {
	f := &FetchProtocol{
		node: node,
	}

	node.logger.Debug("launched fetch protocol multiplex handler")

	go f.onFetchRequest(stream)
	go f.onFetchResponse(stream)

	return f
}

func (f *FetchProtocol) onFetchRequest(s network.Stream) {
	f.node.logger.Debug("fetch request",
		slog.String("protocol", "fetch"))
	s.Close()
}

func (f *FetchProtocol) onFetchResponse(s network.Stream) {
	f.node.logger.Debug("fetch response",
		slog.String("protocol", "fetch"))
	s.Close()
}
