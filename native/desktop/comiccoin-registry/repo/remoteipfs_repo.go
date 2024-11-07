package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"

	httphandler "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http/handler"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

type RemoteIPFSRepo struct {
	logger        *slog.Logger
	remoteAddress string
	apiKey        string
}

const (
	versionURL = "/version"
)

// NewRemoteIPFSRepo returns a new RemoteIPFSRepo instance
func NewRemoteIPFSRepo(logger *slog.Logger, remoteAddress string, apiKey string) domain.RemoteIPFSRepository {
	return &RemoteIPFSRepo{
		logger:        logger,
		remoteAddress: remoteAddress,
		apiKey:        apiKey,
	}
}

func (r *RemoteIPFSRepo) Version(ctx context.Context) (string, error) {
	//
	// STEP 1:
	// Make `GET` request to HTTP JSON API.
	//

	httpEndpoint := fmt.Sprintf("http://%s%s", r.remoteAddress, versionURL)

	httpClient, err := http.NewRequest("GET", httpEndpoint, nil)
	if err != nil {
		log.Fatalf("failed to setup get request: %v", err)
	}
	httpClient.Header.Add("Content-Type", "application/json")

	r.logger.Debug("Get version from remote HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "GET"))

	client := &http.Client{}
	resp, err := client.Do(httpClient)
	if err != nil {
		log.Fatalf("failed to do get request: %v", err)
	}

	defer resp.Body.Close()

	//
	// STEP 2:
	// Handle response.
	//

	if resp.StatusCode == http.StatusNotFound {
		return "", fmt.Errorf("http endpoint does not exist for: %v", httpEndpoint)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Failed to access API: %v", "-")
	}

	//
	// STEP 3:
	// Return the response to the app.
	//

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	respContent := &httphandler.VersionResponseIDO{}
	if err := json.NewDecoder(teeReader).Decode(&respContent); err != nil {
		r.logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return "", err
	}

	return respContent.Version, nil
}

func (r *RemoteIPFSRepo) PinAddViaFilepath(ctx context.Context, fullFilePath string) (string, error) {
	return "", nil
}

func (r *RemoteIPFSRepo) Get(ctx context.Context, cidString string) ([]byte, string, error) {
	return nil, "", nil
}
