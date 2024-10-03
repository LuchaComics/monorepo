package simple

import (
	"encoding/json"

	"github.com/libp2p/go-libp2p/core/peer"
)

func copyMapDeepResponses(m map[peer.ID][]*SimpleMessageResponse) map[peer.ID][]*SimpleMessageResponse {
	data, err := json.Marshal(m)
	if err != nil {
		// handle error
	}
	var newMap map[peer.ID][]*SimpleMessageResponse
	err = json.Unmarshal(data, &newMap)
	if err != nil {
		// handle error
	}
	return newMap
}

func copyMapDeepRequests(m map[peer.ID][]*SimpleMessageRequest) map[peer.ID][]*SimpleMessageRequest {
	data, err := json.Marshal(m)
	if err != nil {
		// handle error
	}
	var newMap map[peer.ID][]*SimpleMessageRequest
	err = json.Unmarshal(data, &newMap)
	if err != nil {
		// handle error
	}
	return newMap
}
