package handler

import (
	"crypto/ecdsa"

	"github.com/ethereum/go-ethereum/common"
	"github.com/google/uuid"
)

type GetKeyRequest struct {
	ID             string
	WalletPassword string
}

type GetKeyResponse struct {
	ID         uuid.UUID         `json:"id"`
	Address    common.Address    `json:"address"`
	PrivateKey *ecdsa.PrivateKey `json:"private_key"`
}

func (h *AccountServer) Execute(args *GetKeyRequest, reply *GetKeyResponse) error {
	// ctx := r.Context()

	key, serviceErr := h.getKeyService.Execute("", "")
	if serviceErr != nil {
		return serviceErr
	}
	*reply = GetKeyResponse{
		ID:         key.Id,
		Address:    key.Address,
		PrivateKey: key.PrivateKey,
	}
	return nil
}
