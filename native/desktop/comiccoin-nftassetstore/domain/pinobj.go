package domain

import (
	"fmt"
	"time"

	"github.com/fxamacker/cbor/v2"
)

// "queued" "pinning" "pinned" "failed"

const (
	StatusQueued  = "queued"
	StatusPinning = "pinning"
	StatusPinned  = "pinned"
	StatusFailed  = "failed"

	ContentTypeFile  = 1
	ContentTypeImage = 2
)

// PinObject is a representation of a pin request. It means it is the IPFS content which we are saving in our system and sharing to the IPFS network, also know as "pinning". This structure has the core variables required to work with IPFS as per their documentation https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers, in additon we also have our applications specific varaibles.
type PinObject struct {
	// RequestID variable is the public viewable unique identifier of this pin.
	RequestID uint64 `json:"requestid"`

	// Status represents operational state of this pinned object in IPFS.
	Status string `json:"status"`

	// CID variable is the unique identifier of our content on the IPFS network. The official definition is: Content Identifier (CID) points at the root of a DAG that is pinned recursively.
	CID string `json:"cid"`

	// Name variable used to provide human readable description for the content. This is optional.
	Name string `json:"name"`

	// The date/time this content was pinned in IPFS network. Developers note: Normally we write it as `CreatedAt`, but IPFS specs require us to write it this way.
	Created time.Time `json:"created,omitempty"`

	// Addresses provided in origins list are relevant only during the initial pinning, and don't need to be persisted by the pinning service
	Origins []string `json:"origins"`

	// Any additional vendor-specific information is returned in optional info.
	Meta map[string]string `json:"meta"`

	// Addresses in the delegates array are peers designated by pinning service that will receive the pin data over bitswap
	Delegates []string `json:"delegates"`

	// Any additional vendor-specific information is returned in optional info.
	Info map[string]string `json:"info"`

	// CODE BELOW IS EXTENSION TO THE IPFS SPEC.
	//------------------------------------------//

	// FileContent variable holds all the content of this pin. Variable will not be saved to database, only returned in API endpoint.
	Content []byte `json:"content,omitempty"`

	Filename    string `json:"filename,omitempty"`
	ContentType int8   `json:"content_type,omitempty"`

	// ID variable is the unique identifier we use internally in our system.
	CreatedFromIPAddress  string    `json:"created_from_ip_address,omitempty"`
	ModifiedAt            time.Time `json:"modified_at,omitempty"`
	ModifiedFromIPAddress string    `json:"modified_from_ip_address,omitempty"`
}

type PinObjectAsSelectOption struct {
	Value uint64 `json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string `json:"label"`
}

// PinObjectRepository Interface for pinobject.
type PinObjectRepository interface {
	Upsert(pinobj *PinObject) error
	// GetByID(ctx context.Context, id uint64) (*PinObject, error)
	GetByCID(cid string) (*PinObject, error)
	GetByRequestID(requestID uint64) (*PinObject, error)
	ListAll() ([]*PinObject, error)
	DeleteByCID(cid string) error
	DeleteByRequestID(requestID uint64) error
	OpenTransaction() error
	CommitTransaction() error
	DiscardTransaction()
}

// Serialize serializes a pin object into a byte array.
// It returns the serialized byte array and an error if one occurs.
func (b *PinObject) Serialize() ([]byte, error) {
	// Marshal the block data DTO into a byte array using CBOR.
	dataBytes, err := cbor.Marshal(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize pin object: %v", err)
	}
	return dataBytes, nil
}

// NewPinObjectFromDeserialize deserializes a pin object from a byte array.
// It returns the deserialized block data DTO and an error if one occurs.
func NewPinObjectFromDeserialize(data []byte) (*PinObject, error) {
	// Variable we will use to return.
	pinobj := &PinObject{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if pinobj == nil {
		return nil, nil
	}

	// Unmarshal the byte array into a block data DTO using CBOR.
	if err := cbor.Unmarshal(data, &pinobj); err != nil {
		return nil, fmt.Errorf("failed to deserialize pin object: %v", err)
	}
	return pinobj, nil
}
