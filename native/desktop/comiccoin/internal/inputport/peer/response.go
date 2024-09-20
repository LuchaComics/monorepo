package peer

import (
	"bufio"
	"fmt"
	"log"
	"strings"
)

func (peer *peerInputPort) responseHandler(rw *bufio.ReadWriter) {
	log.Println("Running response handler...")
	// for {
	// 	bin, _ := rw.ReadBytes('\n')
	//
	// 	if bin == nil {
	// 		break
	// 	}
	// 	// Green console colour: 	\x1b[32m
	// 	// Reset console colour: 	\x1b[0m
	// 	newBlock, err := blockchain.DeserializeBlock(bin)
	// 	if err != nil {
	// 		break
	// 	}
	// 	if err := peer.blockchain.AddBlock(newBlock); err != nil {
	// 		log.Printf("popBlockFromPeer err: %v\n", err)
	// 		break
	// 	}
	//
	// }

	for {
		log.Println("Waiting for a request...")
		requestStr, err := rw.ReadString('\n')
		if err != nil {
			fmt.Printf("Error reading from buffer: %v\n", err)
			break
		}

		if requestStr == "" {
			log.Printf("Request had empty data")
			break
		}

		requestStr = strings.ReplaceAll(requestStr, "\n", "")

		req := strings.Split(requestStr, "|")
		fmt.Printf("req: %v\n", req)
		fmt.Printf("req[0]: %v\n", req[0])
		fmt.Printf("req[1]: %v\n", req[1])
		if req[0] == "FETCH_BLOCK" && req[1] == "NONE" {
			err := peer.kvs.ViewFromFirst(func(key, value []byte) error {
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
	}
	log.Println("Finished running response handler")
}

func processFunc(key, value []byte) error {
	return nil
}
