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
)

// newHostWithPredictableIdentifier creates a host with a predictable identifier.
// This is achieved by reusing a custom private-public key pair saved locally in the system.
func (node *peerProviderImpl) newHostWithPredictableIdentifier() (host.Host, error) {
	// Check if the listen port is configured.
	if node.cfg.Peer.ListenPort == 0 {
		return nil, fmt.Errorf("missing p2p listen port: %v", node.cfg.Peer.ListenPort)
	}

	// Load the private key from the node's identity.
	priv := node.identityPrivateKey

	// Create a connection manager to limit the number of connections.
	connmgr, err := connmgr.NewConnManager(
		100, // Lowwater
		400, // HighWater,
		connmgr.WithGracePeriod(time.Minute),
	)
	if err != nil {
		// Panic if the connection manager cannot be created.
		panic(err)
	}

	// Create a new host with the specified configuration.
	basicHost, err := libp2p.New(
		// Listen on the configured port for TCP and UDP connections.
		libp2p.ListenAddrStrings(
			fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", node.cfg.Peer.ListenPort), // regular tcp connections
			fmt.Sprintf("/ip4/127.0.0.1/udp/%d/quic-v1", node.cfg.Peer.ListenPort),
		),
		// Use the loaded private key for the host's identity.
		libp2p.Identity(priv),
		// Support TLS connections.
		libp2p.Security(libp2ptls.ID, libp2ptls.New),
		// Support noise connections.
		libp2p.Security(noise.ID, noise.New),
		// Support any other default transports (TCP).
		libp2p.DefaultTransports,
		// Attach the connection manager to the host.
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
		// Log an error if the host cannot be created.
		node.logger.Error("failed starting basic host",
			slog.Any("error", err))
		return nil, err
	}

	return basicHost, nil
}
