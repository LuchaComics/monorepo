package peer

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	mrand "math/rand"
	"strings"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type peerInputPort struct {
	cfg           *config.Config
	kvs           keyvaluestore.KeyValueStorer
	blockchain    *blockchain.Blockchain
	host          host.Host
	peerAddresses []*peer.AddrInfo
}

func NewInputPort(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	bc *blockchain.Blockchain,
) inputport.InputPortServer {

	if cfg.Peer.BootstrapPeers == "" {
		//
		// CASE 1 OF 2:
		// HOST MODE
		//

		// Create a new P2P host on the listening port, else fatally error app.
		h, err := makeBasicHost(cfg.Peer.ListenPort, false, cfg.Peer.RandomSeed)
		if err != nil {
			log.Fatalf("failed make p2p host: %v", err)
		}

		return &peerInputPort{
			cfg:        cfg,
			kvs:        kvs,
			blockchain: bc,
			host:       h,
		}
	} else {
		//
		// CASE 2 OF 2:
		// DIAL MODE
		//

		h, err := libp2p.New(
			libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", cfg.Peer.ListenPort)),
			// libp2p.Identity(priv),
		)
		if err != nil {
			log.Fatalf("failed make p2p peer: %v", err)
		}

		// Get our peer addresses.
		peerAddresses := make([]*peer.AddrInfo, 0)
		for _, bootstrapPeer := range strings.Split(cfg.Peer.BootstrapPeers, ",") {
			// Note:
			// https://github.com/libp2p/go-libp2p/blob/master/examples/chat/chat.go#L195C1-L234C2

			// Turn the `bootstrapPeer` into a multiaddr.
			maddr, err := multiaddr.NewMultiaddr(bootstrapPeer)
			if err != nil {
				log.Fatalf("failed to create new multi-addr: %v", err)
			}

			// Extract the peer ID from the multiaddr.
			info, err := peer.AddrInfoFromP2pAddr(maddr)
			if err != nil {
				log.Fatalf("failed to get info from p2p addr: %v", err)
			}

			// Keep a record of our peer addresses.
			peerAddresses = append(peerAddresses, info)
		}

		return &peerInputPort{
			cfg:           cfg,
			kvs:           kvs,
			blockchain:    bc,
			host:          h,
			peerAddresses: peerAddresses,
		}
	}
}

func (peer *peerInputPort) Run() {
	if len(peer.peerAddresses) == 0 {
		//
		// CASE 1 OF 2:
		// HOST MODE
		//

		// Set a function as stream handler.
		// This function is called when a peer connects, and starts a stream with this protocol.
		// Only applies on the receiving side.
		peer.host.SetStreamHandler("/p2p/1.0.0", peer.handleStream)

		log.Printf("peer address: /ip4/127.0.0.1/tcp/<PORT>/p2p/%s\n", peer.host.ID())
		log.Println("peer: waiting for incoming connection")
		log.Println()
	} else {
		//
		// CASE 2 OF 2:
		// PEER MODE
		//

		for _, bootstrapPeer := range peer.peerAddresses {
			// Note:
			// https://github.com/libp2p/go-libp2p/blob/master/examples/chat/chat.go#L195C1-L234C2

			s, err := startPeerAndConnect(context.Background(), peer.host, bootstrapPeer)
			if err != nil {
				log.Fatalf("failed to start and connect to peer: %v", err)
			}
			peer.handleStream(s)
		}
	}
}

func (s *peerInputPort) Shutdown() {
	log.Println("peer: shutting down")
	s.host.Close()
}

// makeBasicHost creates a LibP2P host with a random peer ID listening on the
// given multiaddress. It will use secio if secio is true.
func makeBasicHost(listenPort int, secio bool, randseed int64) (host.Host, error) {

	// If the seed is zero, use real cryptographic randomness. Otherwise, use a
	// deterministic randomness source to make generated keys stay the same
	// across multiple runs
	var r io.Reader
	if randseed == 0 {
		r = rand.Reader
	} else {
		// Do not use this in production
		r = mrand.New(mrand.NewSource(randseed))
	}

	// Generate a key pair for this host. We will use it
	// to obtain a valid host ID.
	priv, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return nil, err
	}

	basicHost, err := libp2p.New(
		libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", listenPort)),
		libp2p.Identity(priv),
	)
	if err != nil {
		return nil, err
	}
	//
	// // Build host multiaddress
	// hostAddr, _ := ma.NewMultiaddr(fmt.Sprintf("/ipfs/%s", basicHost.ID()))
	//
	// // Now we can build a full multiaddress to reach this host
	// // by encapsulating both addresses:
	// addr := basicHost.Addrs()[0]
	// fullAddr := addr.Encapsulate(hostAddr)
	// log.Printf("I am %s\n", fullAddr)
	// if secio {
	// 	log.Printf("Now run \"go run main.go -l %d -d %s -secio\" on a different terminal\n", listenPort+1, fullAddr)
	// } else {
	// 	log.Printf("Now run \"go run main.go -l %d -d %s\" on a different terminal\n", listenPort+1, fullAddr)
	// }

	return basicHost, nil
}

func startPeerAndConnect(ctx context.Context, h host.Host, info *peer.AddrInfo) (network.Stream, error) {
	// Note:
	// https://github.com/libp2p/go-libp2p/blob/master/examples/chat/chat.go#L195C1-L234C2

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := h.NewStream(ctx, info.ID, "/p2p/1.0.0")
	if err != nil {
		return nil, fmt.Errorf("failed to create new stream: %v", err)
	}
	log.Println("Established connection to destination")

	return s, nil
}

func (peer *peerInputPort) handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go peer.responseHandler(rw) // a.k.a. read from network
	go peer.requestHandler(rw)  // a.k.a. writer to network

	// stream 's' will stay open until you close it (or the other side closes it).
}

const (
	StreamMessageTypeRequestLatestBlock = 0
	StreamMessageTypeRespondLatestBlock = 1

	StreamMessageTypeRequestBlock = 2
	StreamMessageTypeRespondBlock = 3
)

type StreamMessageIDO struct {
	Type int `json:"type"`

	Hash string `json:"hash"`

	Block *blockchain.Block `json:"block"`
}
