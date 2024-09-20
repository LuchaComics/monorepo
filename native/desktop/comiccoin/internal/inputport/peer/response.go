package peer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/blockchain"
)

const (
	StreamCommunicationTypeRequestBlockchain = iota
	StreamCommunicationTypeRequestBlock      = iota
	StreamCommunicationTypeRespondBlockchain = iota
	StreamCommunicationTypeRespondBlock      = iota
)

type StreamCommunicationIDO struct {
	Type int `json:"type"`

	Hash string `json:"hash"`

	Block []byte `json:"block"`
}

func (peer *peerInputPort) responseHandler(rw *bufio.ReadWriter) {
	log.Println("Running response handler...")
	for {
		log.Println("Waiting for a request...")
		requestStr, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from buffer: %v\n", err)
			break
		}

		// Trim the newline characters
		requestStr = strings.TrimSpace(requestStr)

		if requestStr == "" {
			log.Printf("Request had empty data")
			break
		}

		// Deserialize the JSON string into StreamRequestIDO
		var req StreamCommunicationIDO
		if err := json.Unmarshal([]byte(requestStr), &req); err != nil {
			fmt.Printf("Error unmarshaling JSON: %v\n", err)
			break
		}

		// Log the deserialized request for debugging
		log.Printf("Request Type: %d, Hash: %s\n", req.Type, req.Hash)

		switch req.Type {
		case StreamCommunicationTypeRequestBlock:
			{
				//TODO: IMPL.
				fmt.Println("Received a request to fetch the latest block")
				err := peer.kvs.View([]byte(fmt.Sprintf("BLOCK_%v", req.Hash)), func(key, value []byte) error {
					// Do something with key and value
					fmt.Printf("Uploading block with key: %s, value: %s\n", key, value)

					// Return nil to indicate success
					return nil
				})
				if err != nil {
					fmt.Printf("Error sending data: %v\n", err)
					break
				}
			}
		case StreamCommunicationTypeRequestBlockchain:
			{
				fmt.Println("Received a request to fetch the entire blockchain")
				err := peer.kvs.ViewFromFirst(func(key, value []byte) error {
					res := &StreamCommunicationIDO{
						Type:  StreamCommunicationTypeRespondBlockchain,
						Block: value,
					}
					bin, _ := json.Marshal(res)
					str := string(bin)
					str = fmt.Sprintf("%s\n", str)
					if _, err := rw.WriteString(str); err != nil {
						log.Printf("Finished running response handler with write error: %v\n", err)
						return err
					}
					if err := rw.Flush(); err != nil {
						log.Fatal(err)
					}

					// Return nil to indicate success
					return nil
				})
				if err != nil {
					fmt.Printf("Error sending response data: %v\n", err)
					break
				}
			}
		case StreamCommunicationTypeRespondBlockchain:
			{
				block, err := blockchain.DeserializeBlock(req.Block)
				if err != nil {
					fmt.Printf("Error unmarshalling response data: %v\n", err)
					break
				}
				fmt.Printf("Saving blockchain data: %v\n", block)
				if err := peer.blockchain.AddBlock(block); err != nil {
					fmt.Printf("failed adding block to blockchain: %v\n", err)
					break
				}

			}
		}
	}
	log.Println("Finished running response handler")
}
