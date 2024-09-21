package p2p

import (
	"context"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p/core/network"
)

const fetchProtocolVersion = "/fetch/1.0.0"

type FetchProtocol struct {
	node   *nodeInputPort
	stream network.Stream
}

func NewFetchProtocol(node *nodeInputPort, stream network.Stream) *FetchProtocol {
	f := &FetchProtocol{
		node:   node,
		stream: stream,
	}
	return f
}

func (f *FetchProtocol) onFetchRequest(ctx context.Context, s network.Stream) {
	f.node.logger.Debug("fetch request",
		slog.String("protocol", "fetch"))

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			f.node.logger.Debug("finished fetching latest block(s) from p2p network")
			s.Close()
			return
		case <-ticker.C:
			// This is where you send the request to the network
			f.node.logger.Debug("fetching latest block(s) from p2p network...")

			// ALGORITHM.
			// 1. Lookup in my local blockchain for the latest hash
			// 2. If no latest hash found, send
			// 2. Submit request
			// 3. Repeat above every minute
		}
	}

}

func (f *FetchProtocol) onFetchResponse(ctx context.Context, s network.Stream) {
	f.node.logger.Debug("fetch response",
		slog.String("protocol", "fetch"))
	s.Close()
}

func (f *FetchProtocol) Run(ctx context.Context) {
	f.node.logger.Debug("launched fetch protocol multiplex handler")
	go f.onFetchRequest(ctx, f.stream)
	go f.onFetchResponse(ctx, f.stream)
}
