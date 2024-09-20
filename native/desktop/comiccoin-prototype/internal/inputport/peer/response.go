package peer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"strings"
)

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
		req := &StreamMessageIDO{}
		if err := json.Unmarshal([]byte(requestStr), &req); err != nil {
			fmt.Printf("Error unmarshaling JSON in stream message: %v\n", err)
			break
		}

		// Log the deserialized request for debugging
		log.Printf("responseHandler | Request Type: %d, Hash: %s\n", req.Type, req.Hash)

		switch req.Type {
		case StreamMessageTypeRequestLatestBlock:
			{
				lastHash := peer.blockchain.LastHash
				log.Println("Sending latest hash:", lastHash)
				block, err := peer.blockchain.GetBlock(lastHash)
				if err != nil {
					log.Printf("failed getting block: %v\n", err)
					log.Printf("StreamMessageTypeRequestLatestBlock: req: %v\n", requestStr)
					break
				}
				if block == nil {
					log.Printf("no block found")
					break
				}

				res := &StreamMessageIDO{
					Type:  StreamMessageTypeRespondLatestBlock,
					Block: block,
				}
				bin, _ := json.Marshal(res)
				str := string(bin)
				str = fmt.Sprintf("%s\n", str)
				if _, err := rw.WriteString(str); err != nil {
					log.Printf("Finished running response handler with write error: %v\n", err)
					break
				}
				if err := rw.Flush(); err != nil {
					log.Fatal(err)
				}
				log.Println("Sent latest hash:", lastHash)
			}
		case StreamMessageTypeRespondLatestBlock:
			{
				log.Printf("Saving latest block: %v\n", req.Block.Hash)
				if err := peer.blockchain.AddBlock(req.Block); err != nil {
					fmt.Printf("failed adding block to blockchain: %v\n", err)
					break
				}
				log.Printf("Saved latest block: %v\n", req.Block.Hash)

				// Do we have the previous block in our database? If so then
				// skip else make another fetch request. Also if we get to the
				// genesis block then stop.
				if req.Block.PreviousHash != "" {
					block, err := peer.blockchain.GetBlock(req.Block.PreviousHash)
					if err != nil || block == nil {
						log.Printf("Fetching previous block...")

						req := &StreamMessageIDO{
							Type: StreamMessageTypeRequestBlock,
							Hash: req.Block.PreviousHash,
						}
						bin, _ := json.Marshal(req)
						str := string(bin)
						str = fmt.Sprintf("%s\n", str)
						if _, err := rw.WriteString(str); err != nil {
							log.Printf("Finished running request handler with write error: %v\n", err)
							return
						}
						if err := rw.Flush(); err != nil {
							log.Fatal(err)
						}
						log.Printf("Submited fetching previous block...")
					}
				}
			}
		case StreamMessageTypeRequestBlock:
			{
				log.Println("Sending block at hash:", req.Hash)
				block, err := peer.blockchain.GetBlock(req.Hash)
				if err != nil {
					log.Printf("failed getting block: %v\n", err)
					log.Printf("StreamMessageTypeRequestBlock: req: %v\n", requestStr)
					break
				}
				if block == nil {
					log.Printf("no block found")
					break
				}

				res := &StreamMessageIDO{
					Type:  StreamMessageTypeRespondBlock,
					Block: block,
				}
				bin, _ := json.Marshal(res)
				str := string(bin)
				str = fmt.Sprintf("%s\n", str)
				if _, err := rw.WriteString(str); err != nil {
					log.Printf("Finished running response handler with write error: %v\n", err)
					break
				}
				if err := rw.Flush(); err != nil {
					log.Fatal(err)
				}

				log.Println("Sent block at hash:", req.Hash)
			}
		case StreamMessageTypeRespondBlock:
			{

				log.Printf("Saving a block: %v\n", req.Block.Hash)
				if err := peer.blockchain.AddBlock(req.Block); err != nil {
					fmt.Printf("failed adding block to blockchain: %v\n", err)
					break
				}
				log.Printf("Saved a block: %v\n", req.Block.Hash)

				// Do we have the previous block in our database? If so then
				// skip else make another fetch request. Also if we get to the
				// genesis block then stop.
				if req.Block.PreviousHash != "" {
					log.Printf("Has previous has, fetching it now...")
					block, err := peer.blockchain.GetBlock(req.Block.PreviousHash)
					if err != nil {
						log.Printf("failed getting block: %v\n", err)
						log.Printf("StreamMessageTypeRespondBlock: req: %v\n", requestStr)
						break
					}
					if block == nil {
						log.Printf("Fetching previous block...")

						req := &StreamMessageIDO{
							Type: StreamMessageTypeRequestBlock,
							Hash: block.PreviousHash,
						}
						bin, _ := json.Marshal(req)
						str := string(bin)
						str = fmt.Sprintf("%s\n", str)
						if _, err := rw.WriteString(str); err != nil {
							log.Printf("Finished running request handler with write error: %v\n", err)
							return
						}
						if err := rw.Flush(); err != nil {
							log.Fatal(err)
						}
						log.Printf("Submited fetching previous block...")

					}
				}
			}
		}
	}
	log.Println("Finished running response handler")
}
