package p2p

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/p2p/net/connmgr"
	"github.com/libp2p/go-libp2p/p2p/security/noise"
	libp2ptls "github.com/libp2p/go-libp2p/p2p/security/tls"
	ma "github.com/multiformats/go-multiaddr"
)

// newHostWithPredictableIdentifier function will create a host with an
// identifier that never changes - it remains the same due to the fact that
// we are using custom private-public key pairs saved locally in the system.
func (node *peerProviderImpl) newHostWithPredictableIdentifier() (host.Host, error) {
	// ctx := context.Background()

	if node.cfg.Peer.ListenPort == 0 {
		return nil, fmt.Errorf("missing to p2p listen port: %v", node.cfg.Peer.ListenPort)
	}

	//
	// STEP 1:
	// We want to keep the same identifier every time the server restarts
	// or restarts so we will reuse the key-pair we have saved locally.
	//

	priv := node.identityPrivateKey

	//
	// STEP 2:
	// Loadup our p2p host server.
	//

	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		panic(err)
	}

	basicHost, err := libp2p.New(
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", node.cfg.Peer.ListenPort), // regular tcp connections
			fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic-v1", node.cfg.Peer.ListenPort),
		),
		// Use the keypair we generated
		libp2p.Identity(priv),
		// support TLS connections
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// support noise connections
		libp2p.Security(noise.ID, noise.New),
		// support any other default transports (TCP)
		libp2p.DefaultTransports,
		// Let's prevent our peer from having too many
		// connections by attaching a connection manager.
		libp2p.ConnectionManager(connmgr),

		// Attempt to open ports using uPNP for NATed hosts.
		// libp2p.NATPortMap(),// TODO: UNCOMMENT WHEN WE ARE READY TO LAUNCH THIS APP INTO PRODUCTION

		// If you want to help other peers to figure out if they are behind
		// NATs, you can launch the server-side of AutoNAT too (AutoRelay
		// already runs the client)
		//
		// This service is highly rate-limited and should not cause any
		// performance issues.
		// libp2p.EnableNATService(), // TODO: UNCOMMENT WHEN WE ARE READY TO LAUNCH THIS APP INTO PRODUCTION

		// DEVELOPERS NOTE:
		// See more options via this link:
		// https://github.com/libp2p/go-libp2p/blob/master/examples/libp2p-host/host.go#L56
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
		slog.Any("full_address", fullAddr),
		slog.Any("host_id", basicHost.ID()),
		slog.String("note", "you can replace `127.0.0.1` with your public ip when connecting to this peer"),
	)

	return basicHost, nil
}
