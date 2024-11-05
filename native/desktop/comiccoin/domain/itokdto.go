package domain

import (
	"context"
	"crypto/ecdsa"
	"fmt"

	"github.com/fxamacker/cbor/v2"
)

type IssuedToken struct {
	ID          uint64 `json:"id"`
	MetadataURI string `json:"metadata_uri"` // ComicCoin: URI pointing to Token metadata file (if this transaciton is an Token).
}

// Sign function signs the  transaction using the user's private key
// and returns a signed version of that transaction.
func (itok *IssuedToken) Sign(privateKey *ecdsa.PrivateKey) ([]byte, error) {
	// // Break the signature into the 3 parts: R, S, and V.
	// v, r, s, err := signature.Sign(itok, privateKey)
	// if err != nil {
	// 	return nil, err
	// }

	return nil, nil
}

// IssuedTokenDTO represents what can be serialized to disk and over the network
// for a token which was newly minted by the Proof of Authority and now needs
// to be issued by that authority to the network.
type IssuedTokenDTO struct {
	Token               *IssuedToken `json:"token"`
	TokenSignatureBytes []byte       `json:"token_signature_bytes"`
	Validator           *Validator   `json:"validator"`
}

type IssuedTokenDTORepository interface {
	BroadcastToP2PNetwork(ctx context.Context, dto *IssuedTokenDTO) error
	ReceiveFromP2PNetwork(ctx context.Context) (*IssuedTokenDTO, error)
}

func (b *IssuedTokenDTO) Serialize() ([]byte, error) {
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize issued token dto: %v", err)
	}
	return dataBytes, nil
}

func NewIssuedTokenDTOFromDeserialize(data []byte) (*IssuedTokenDTO, error) {
	// Variable we will use to return.
	issuedTok := &IssuedTokenDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	if err := cbor.Unmarshal(data, &issuedTok); err != nil {
		return nil, fmt.Errorf("failed to deserialize issued token dto: %v", err)
	}
	return issuedTok, nil
}
