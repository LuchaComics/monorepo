package peer

import (
	"bufio"
	"context"
	"crypto/rand"
	"fmt"
	"io"
	"log"
	mrand "math/rand"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/multiformats/go-multiaddr"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/adapter/keyvaluestore"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type peerInputPort struct {
	cfg        *config.Config
	kvs        keyvaluestore.KeyValueStorer
	blockchain *blockchain.Blockchain
	host       host.Host
}

func NewInputPort(
	cfg *config.Config,
	kvs keyvaluestore.KeyValueStorer,
	bc *blockchain.Blockchain,
) inputport.InputPortServer {
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
}

func (peer *peerInputPort) Run() {
	startPeer(context.Background(), peer.host, peer.handleStream)
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

func startPeer(ctx context.Context, h host.Host, streamHandler network.StreamHandler) {
	// Set a function as stream handler.
	// This function is called when a peer connects, and starts a stream with this protocol.
	// Only applies on the receiving side.
	h.SetStreamHandler("/comic-coin/1.0.0", streamHandler)

	// Let's get the actual TCP port from our listen multiaddr, in case we're using 0 (default; random available port).
	var port string
	for _, la := range h.Network().ListenAddresses() {
		if p, err := la.ValueForProtocol(multiaddr.P_TCP); err == nil {
			port = p
			break
		}
	}

	if port == "" {
		log.Println("was not able to find actual local port")
		return
	}

	log.Printf("peer address: /ip4/127.0.0.1/tcp/%v/p2p/%s\n", port, h.ID())
	log.Println("peer note: you can replace 127.0.0.1 with public IP as well.")
	log.Println("peer: waiting for incoming connection")
	log.Println()
}

func (peer *peerInputPort) handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go peer.synchLocalBlockchainWithNetwork(rw) // a.k.a. read from network
	go peer.synchNetworkWithLocalBlockchain(rw) // a.k.a. writer to network

	// stream 's' will stay open until you close it (or the other side closes it).
}

func (peer *peerInputPort) synchLocalBlockchainWithNetwork(rw *bufio.ReadWriter) {
	for {
		str, _ := rw.ReadString('\n')

		if str == "" {
			return
		}
		if str != "\n" {
			// Green console colour: 	\x1b[32m
			// Reset console colour: 	\x1b[0m
			fmt.Printf("\x1b[32m%s\x1b[0m> ", str)
		}

	}
}

// synchNetworkWithLocalBlockchain listens for new blocks from our local blockchain and if new blocks come in then we broadcast them to the P2P network
func (peer *peerInputPort) synchNetworkWithLocalBlockchain(rw *bufio.ReadWriter) {
	log.Println("Waiting to receive new blocks from the local blockchain so we can publish to p2p network...")
	for newBlock := range peer.blockchain.Subscribe() {
		fmt.Printf("New local block received: %v\n", newBlock)

		sendData := newBlock.Serialize()

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()

		fmt.Println("New local block sent to network")
	}
}
