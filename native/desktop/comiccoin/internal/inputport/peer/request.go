package peer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
)

const (
	StreamRequestTypeGetBlock = iota
)

type StreamRequestIDO struct {
	Type int    `json:"type"`
	Hash string `json:"hash"`
}

// requestPeer listens for new blocks from our local blockchain and if new blocks come in then we broadcast them to the P2P network
func (peer *peerInputPort) requestHandler(rw *bufio.ReadWriter) {
	log.Println("Running request handler...")

	//
	// STEP 1:
	// On startup of application, if our blockchain is empty then we need to
	// send a request to our peer to download the entire blockchain locally.
	//

	if peer.blockchain.IsBlockchainEmpty {
		req := &StreamCommunicationIDO{
			Type: StreamCommunicationTypeRequestBlockchain,
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

		log.Println("Blockchain download request submitted successfully.")
	}
	if peer.blockchain.IsBlockchainInSynch == false {
		//TODO:
	}

	//
	// STEP 2:
	// We need to monitor any changes (e.i. block additions into our blockchain)
	// that occured locally in our blockchain and then push our new record to
	// the network.
	//

	// for newBlock := range peer.blockchain.Subscribe() {
	// 	fmt.Printf("New local block received: %v\n", newBlock)
	//
	// 	sendData := newBlock.Serialize()
	//
	// 	_, err := rw.WriteString(fmt.Sprintf("%s\n", sendData))
	// 	if err != nil {
	// 		break
	// 	}
	//
	// 	if err := rw.Flush(); err != nil {
	// 		break
	// 	}
	//
	// 	fmt.Println("New local block sent to network")
	// }

	log.Println("Finished running request handler...")
}

func (peer *peerInputPort) runFetchLatestOperation() {

}
