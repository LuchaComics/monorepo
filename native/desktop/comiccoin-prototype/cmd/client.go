package cmd

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/host"
	"github.com/libp2p/go-libp2p/core/network"
	"github.com/libp2p/go-libp2p/core/peer"
	"github.com/libp2p/go-libp2p/core/peerstore"
	"github.com/multiformats/go-multiaddr"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(clientCmd)
	clientCmd.AddCommand(clientHelloCmd())
}

var (
	flagPeerAddress string
	flagListenPort  int
)

var clientCmd = &cobra.Command{
	Use:   "client",
	Short: "Connect to remote nodes",
	Run: func(cmd *cobra.Command, args []string) {
		// Do nothing...
	},
}

func clientHelloCmd() *cobra.Command {
	var cmd = &cobra.Command{
		Use:   "hello",
		Short: "Sends a hello message to peer, use this to test connection",
		Run: func(cmd *cobra.Command, args []string) {
			host, err := libp2p.New(
				libp2p.ListenAddrStrings(fmt.Sprintf("/ip4/127.0.0.1/tcp/%d", flagListenPort)),
				// libp2p.Identity(priv),
			)
			if err != nil {
				panic(err)
			}
			defer host.Close()

			rw, err := startPeerAndConnect(context.Background(), host, flagPeerAddress)
			if err != nil {
				log.Println(err)
				return
			}

			// Create a thread to read and write data.
			go writeData(rw)
			go readData(rw)

			// Wait forever
			select {}
		},
	}

	// cmd.Flags().StringVar(&flagDataDir, "datadir", "./data", "Absolute path to your node's data dir where the DB will be/is stored")
	// // cmd.MarkFlagRequired("datadir")
	cmd.Flags().StringVar(&flagPeerAddress, "peer-address", "", "The address of the remote node")
	cmd.MarkFlagRequired("peer-address")
	cmd.Flags().IntVar(&flagListenPort, "listen-port", 9001, "The address of the remote node")
	cmd.MarkFlagRequired("listen-port")

	return cmd
}

func startPeerAndConnect(ctx context.Context, h host.Host, destination string) (*bufio.ReadWriter, error) {
	log.Println("This node's multiaddresses:")
	for _, la := range h.Addrs() {
		log.Printf(" - %v\n", la)
	}
	log.Println()

	// Turn the destination into a multiaddr.
	maddr, err := multiaddr.NewMultiaddr(destination)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Extract the peer ID from the multiaddr.
	info, err := peer.AddrInfoFromP2pAddr(maddr)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	// Add the destination's peer multiaddress in the peerstore.
	// This will be used during connection and stream creation by libp2p.
	h.Peerstore().AddAddrs(info.ID, info.Addrs, peerstore.PermanentAddrTTL)

	// Start a stream with the destination.
	// Multiaddress of the destination peer is fetched from the peerstore using 'peerId'.
	s, err := h.NewStream(context.Background(), info.ID, "/comic-/1.0.0")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Println("Established connection to destination")

	// Create a buffered stream so that read and writes are non-blocking.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	return rw, nil
}

func handleStream(s network.Stream) {
	log.Println("Got a new stream!")

	// Create a buffer stream for non-blocking read and write.
	rw := bufio.NewReadWriter(bufio.NewReader(s), bufio.NewWriter(s))

	go readData(rw)
	go writeData(rw)

	// stream 's' will stay open until you close it (or the other side closes it).
}

func readData(rw *bufio.ReadWriter) {
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

func writeData(rw *bufio.ReadWriter) {
	stdReader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("> ")
		sendData, err := stdReader.ReadString('\n')
		if err != nil {
			log.Println(err)
			return
		}

		rw.WriteString(fmt.Sprintf("%s\n", sendData))
		rw.Flush()
	}
}
