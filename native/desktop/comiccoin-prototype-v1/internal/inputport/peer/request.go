package peer

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
)

// requestPeer listens for new blocks from our local blockchain and if new blocks come in then we broadcast them to the P2P network
func (peer *peerInputPort) requestHandler(rw *bufio.ReadWriter) {
	log.Println("Running request handler...")

	//
	// STEP 1:
	//

	req := &StreamMessageIDO{
		Type: StreamMessageTypeRequestLatestBlock,
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

	// if peer.blockchain.IsBlockchainEmpty {
	// 	log.Println("Fetching request to download entire blockchain....")
	// 	req := &StreamMessageIDO{
	// 		Type: StreamMessageTypeRequestBlockchain,
	// 	}
	// 	bin, _ := json.Marshal(req)
	// 	str := string(bin)
	// 	str = fmt.Sprintf("%s\n", str)
	// 	if _, err := rw.WriteString(str); err != nil {
	// 		log.Printf("Finished running request handler with write error: %v\n", err)
	// 		return
	// 	}
	// 	if err := rw.Flush(); err != nil {
	// 		log.Fatal(err)
	// 	}
	//
	// 	log.Println("Fetch request submitted successfully.")
	// }
	// if peer.blockchain.IsBlockchainInSynch == false {
	// 	//TODO:
	// }

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
