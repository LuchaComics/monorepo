package repo

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

type IPFSRepo struct {
	logger      *slog.Logger
	APIEndpoint string
}

// Identity represents the identity of an IPFS node
type Identity struct {
	Addresses    []string `json:"Addresses"`
	AgentVersion string   `json:"AgentVersion"`
	ID           string   `json:"ID"`
	Protocols    []string `json:"Protocols"`
	PublicKey    string   `json:"PublicKey"`
}

// NewIPFSRepo returns a new IPFSNode instance
func NewIPFSRepo(logger *slog.Logger, apiEndpoint string) *IPFSRepo {
	return &IPFSRepo{logger: logger, APIEndpoint: apiEndpoint}
}

// ID returns the IPFS node's identity information
func (r *IPFSRepo) ID() (*Identity, error) {
	req, err := http.NewRequest("POST", r.APIEndpoint+"/api/v0/id", nil)
	if err != nil {
		r.logger.Debug("failed to create request",
			slog.Any("error", err))
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.logger.Debug("failed to post",
			slog.Any("error", err))
		return nil, err
	}
	defer resp.Body.Close()

	var identity Identity
	err = json.NewDecoder(resp.Body).Decode(&identity)
	if err != nil {
		r.logger.Debug("failed to decode",
			slog.Any("resp", resp),
			slog.Any("error", err))
		return nil, err
	}

	return &identity, nil
}
