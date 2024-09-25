package p2ptransport

import (
	"context"
	"log/slog"

	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"

	mempool_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/mempool/controller"
)

// MempoolNode represents our `mempool` node in our distributed network.
type MempoolNode struct {
	logger     *slog.Logger
	controller mempool_c.MempoolController

	host   host.Host
	stream network.Stream
}

// NewHandler Constructor
func NewNode(loggerp *slog.Logger, c mempool_c.MempoolController) *MempoolNode {
	return &MempoolNode{
		logger:     loggerp,
		controller: c,
	}
}

func (node *MempoolNode) Handle(ctx context.Context, host host.Host, peerID peer.ID) {
	node.host = host

	if peerID == "" {
		node.host.SetStreamHandler("/mempool-purpose/1.0.0", func(stream network.Stream) {
			node.logger.Info("Got a new stream via host")
			// fp := NewFetchProtocol(node, stream)
			// go fp.Run(ctx)
			// // 'stream' will stay open until you close it (or the other side closes it).
		})

	} else {
		stream, err := node.host.NewStream(ctx, peerID, "/mempool-purpose/1.0.0")
		if err != nil {
			node.logger.Warn("Connection failed:",
				slog.Any("error", err))
			return
		} else {
			node.logger.Info("Got a new stream via peer!")

			// 'stream' will stay open until you close it (or the other side closes it).
		}
		_ = stream
	}

	// node.host.SetStreamHandler("/mempool/purpose/1.0.0", func(stream network.Stream) {
	// 	node.logger.Info("Got a new stream!")
	// })
}
