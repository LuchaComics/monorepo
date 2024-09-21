package p2p

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	ma "github.com/multiformats/go-multiaddr"
)

// newHostWithPredictableIdentifier function will create a host with an
// identifier that never changes - it remains the same due to the fact that
// we are using custom private-public key pairs saved locally in the system.
func (node *nodeInputPort) newHostWithPredictableIdentifier() (host.Host, error) {
	ctx := context.Background()

	//
	// STEP 1:
	// We want to keep the same identifier every time the server restarts
	// or restarts so we will reuse the key-pair we have saved locally.
	//

	priv, _, err := node.keypairStorer.GetByName(ctx, node.cfg.Peer.KeyName)
	if err != nil {
		node.logger.Error("failed getting keypair by name",
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 2:
	// Loadup our p2p host server.
	//

	basicHost, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", node.cfg.Peer.ListenPort)),
		libp2p.Identity(priv),
	)
	if err != nil {
		node.logger.Error("failed starting basic host",
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 3:
	// Output to user's console some useful information.
	//

	// Build host multiaddress
	hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID()))

	// Now we can build a full multiaddress to reach this host
	// by encapsulating both addresses:
	addr := basicHost.Addrs()[0]
	fullAddr := addr.Encapsulate(hostAddr)

	node.logger.Info("host ready to accept peers",
		slog.Any("host_address", hostAddr),
		slog.Any("full_address", fullAddr),
		slog.String("note", "you can replace `127.0.0.1` with your public ip when connecting to this peer"),
	)

	return basicHost, nil
}
