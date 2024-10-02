package domain

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type StreamMessageDTO struct {
	FunctionID string
	Type       int
	Content    []byte
}

const (
	StreamMessageDTOTypeRequest  = 1
	StreamMessageDTOTypeResponse = 2
)

func (b *StreamMessageDTO) Serialize() ([]byte, error) {
	var result bytes.Buffer
	encoder := gob.NewEncoder(&result)
	err := encoder.Encode(b)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize stream message dto: %v", err)
	}
	return result.Bytes(), nil
}

func NewMessageDTOFromDeserialize(data []byte) (*StreamMessageDTO, error) {
	// Variable we will use to return.
	dto := &StreamMessageDTO{}

	// Defensive code: If programmer entered empty bytes then we will
	// return nil deserialization result.
	if data == nil {
		return nil, nil
	}

	decoder := gob.NewDecoder(bytes.NewReader(data))
	err := decoder.Decode(&dto)
	if err != nil {
		return nil, fmt.Errorf("failed to deserialize stream message dto: %v", err)
	}
	return dto, nil
}
