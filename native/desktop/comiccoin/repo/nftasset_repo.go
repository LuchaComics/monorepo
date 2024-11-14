// The NFTAssetRepo package provides a set of functions for interacting with a remote IPFS repository custom built to exclusively host NFT assets for the ComicCoin blockchain. It allows users to retrieve the version of the NFT store service, pin files, and retrieve files.
package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	pkgfilepath "path/filepath"
	"strings"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

// Constants defining API endpoint paths
const (
	versionURL = "/version"      // Path for checking IPFS service version
	pinAddURL  = "/ipfs/pin-add" // Path for adding a pin to IPFS
	gatewayURL = "/ipfs/${CID}"  // Path for fetching content from IPFS by CID
)

// NFTAssetRepoConfigurationProvider is an interface for configuration providers
// that provide remote address and API key needed for IPFS interaction.
type NFTAssetRepoConfigurationProvider interface {
	GetNFTAssetNodeAddress() string // Retrieves the remote IPFS service address
	GetNFTAssetNodeAPIKey() string  // Retrieves the API key for authentication
}

// NFTAssetRepo handles interactions with a remote IPFS service.
type NFTAssetRepo struct {
	config NFTAssetRepoConfigurationProvider // Holds IPFS connection configuration
	logger *slog.Logger                      // Logger instance for logging debug and error messages
}

// NFTAssetRepoConfigurationProviderImpl is a struct that implements
// NFTAssetRepoConfigurationProvider for storing IPFS connection details.
type NFTAssetRepoConfigurationProviderImpl struct {
	remoteAddress string // Address of the IPFS service
	apiKey        string // API key for accessing IPFS service
}

// NewNFTAssetRepoConfigurationProvider constructs a new configuration provider
// for IPFS connections with the specified remote address and API key.
func NewNFTAssetRepoConfigurationProvider(remoteAddress string, apiKey string) NFTAssetRepoConfigurationProvider {
	// Defensive code: Enforce `remoteAddress` is set at minimum.
	if remoteAddress == "" {
		log.Fatal("Missing `remoteAddress` parameter.")
	}
	return &NFTAssetRepoConfigurationProviderImpl{
		remoteAddress: remoteAddress,
		apiKey:        apiKey,
	}
}

// GetNFTAssetNodeAddress retrieves the remote IPFS service address.
func (impl *NFTAssetRepoConfigurationProviderImpl) GetNFTAssetNodeAddress() string {
	return impl.remoteAddress
}

// GetNFTAssetNodeAPIKey retrieves the API key for IPFS service authentication.
func (impl *NFTAssetRepoConfigurationProviderImpl) GetNFTAssetNodeAPIKey() string {
	return impl.apiKey
}

// NewNFTAssetRepo initializes a new NFTAssetRepo instance with the provided configuration and logger.
func NewNFTAssetRepo(cfg NFTAssetRepoConfigurationProvider, logger *slog.Logger) domain.NFTAssetRepository {
	return &NFTAssetRepo{
		config: cfg,
		logger: logger,
	}
}

// NewNFTAssetRepoWithConfiguration initializes a new NFTAssetRepo instance directly with address and API key.
func NewNFTAssetRepoWithConfiguration(logger *slog.Logger, remoteAddress string, apiKey string) domain.NFTAssetRepository {
	return &NFTAssetRepo{
		config: NewNFTAssetRepoConfigurationProvider(remoteAddress, apiKey),
		logger: logger,
	}
}

// Version fetches the version of the remote IPFS service.
// It makes a GET request to the version endpoint and parses the JSON response.
func (r *NFTAssetRepo) Version(ctx context.Context) (string, error) {
	//
	// STEP 1:
	// Make `GET` request to HTTP JSON API.
	//

	httpEndpoint := fmt.Sprintf("%s%s", r.config.GetNFTAssetNodeAddress(), versionURL)

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

	// Copy and pasted from "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/interface/http/handler".
	type VersionResponseIDO struct {
		Version string `json:"version"`
	}

	respContent := &VersionResponseIDO{}
	if err := json.NewDecoder(teeReader).Decode(&respContent); err != nil {
		r.logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return "", err
	}

	return respContent.Version, nil
}

// PinAddViaFilepath uploads a file to IPFS via the specified file path and pins it.
// The function prepares a multipart form request with the file content and metadata, and sends it.
func (r *NFTAssetRepo) PinAddViaFilepath(ctx context.Context, fullFilePath string) (string, error) {
	// Defensive Code: Enforce access based on if API key was set.
	if r.config.GetNFTAssetNodeAPIKey() == "" {
		r.logger.Warn("Cannot pin and add because no API key set.")
		return "", errors.New("Cannot pin and add because no API key set.")
	}

	//
	// STEP 1:
	// Open the file and extract the file details.
	//

	file, err := os.Open(fullFilePath)
	if err != nil {
		r.logger.Error("Failed opening file.",
			slog.Any("err", err),
			slog.String("fullFilePath", fullFilePath))
		return "", err
	}
	defer file.Close()

	// Detect content type of the file
	buffer := make([]byte, 512) // 512 bytes are sufficient for content detection
	_, err = file.Read(buffer)
	if err != nil {
		log.Fatalf("failed to read file for content detection: %v", err)
	}
	file.Seek(0, 0) // Reset file pointer after reading for detection
	contentType := http.DetectContentType(buffer)

	// Get the filename from the filepath.
	fileName := pkgfilepath.Base(fullFilePath)

	//
	// STEP 2:
	// Add the file to the form.
	//

	// Create a buffer to write the multipart form data
	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	// Create a form field writer for the file field
	fileField, err := writer.CreateFormFile("data", fileName)
	if err != nil {
		return "", fmt.Errorf("failed to create form file: %v", err)
	}

	// Copy the contents of the *os.File to the multipart form field
	if _, err := io.Copy(fileField, file); err != nil {
		return "", fmt.Errorf("failed to copy file to form field: %v", err)
	}

	//
	// STEP 3:
	// Close the form
	//

	// Close the multipart writer to finalize the form data
	if err := writer.Close(); err != nil {
		return "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Send HTTP request with the multipart form data
	req, err := http.NewRequest("POST", fmt.Sprintf("%v%v", r.config.GetNFTAssetNodeAddress(), pinAddURL), &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v\n", err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "JWT " + string(r.config.GetNFTAssetNodeAPIKey())

	// Add headers
	req.Header.Add("Authorization", bearer)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	req.Header.Set("X-File-Content-Type", contentType) // Custom header to carry content type

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %v\n", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		e := make(map[string]string)
		var rawJSON bytes.Buffer
		teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

		// Try to decode the response as a string first
		var jsonStr string
		err := json.NewDecoder(teeReader).Decode(&jsonStr)
		if err != nil {
			r.logger.Error("decoding string error",
				slog.Any("err", err),
				slog.String("json", rawJSON.String()),
			)
			return "", err
		}

		// Now try to decode the string into a map
		err = json.Unmarshal([]byte(jsonStr), &e)
		if err != nil {
			r.logger.Error("decoding map error",
				slog.Any("err", err),
				slog.String("json", jsonStr),
			)
			return "", err
		}

		r.logger.Debug("Parsed error response",
			slog.Any("errors", e))
		return "", err
	}

	//
	// STEP 5:
	// Print the success message.
	//

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(resp.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	// Copied from `"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/service/ipfs_pinadd.go"`.
	type IPFSPinAddResponseIDO struct {
		RequestID uint64            `bson:"requestid" json:"requestid"`
		Status    string            `bson:"status" json:"status"`
		Created   time.Time         `bson:"created,omitempty" json:"created,omitempty"`
		Delegates []string          `bson:"delegates" json:"delegates"`
		Info      map[string]string `bson:"info" json:"info"`
		CID       string            `bson:"cid" json:"cid"`
		Name      string            `bson:"name" json:"name"`
		Origins   []string          `bson:"origins" json:"origins"`
		Meta      map[string]string `bson:"meta" json:"meta"`
	}

	post := &IPFSPinAddResponseIDO{}
	if err := json.NewDecoder(teeReader).Decode(&post); err != nil {
		r.logger.Error("decoding string error",
			slog.Any("err", err),
			slog.String("json", rawJSON.String()),
		)
		return "", err
	}

	r.logger.Debug("Submitted successfully",
		slog.Any("RequestID", post.RequestID),
		slog.Any("Status", post.Status),
		slog.Any("Created", post.Created),
		slog.Any("Delegates", post.Delegates),
		slog.Any("Info", post.Info),
		slog.Any("CID", post.CID),
		slog.Any("Name", post.Name),
		slog.Any("Origins", post.Origins),
		slog.Any("Meta", post.Meta))

	return post.CID, nil
}

// Get retrieves content from the IPFS service by a given CID.
// This method constructs the URL for the IPFS gateway and fetches the content.
func (r *NFTAssetRepo) Get(ctx context.Context, cidString string) (*domain.NFTAsset, error) {
	modifiedGatewayURL := strings.ReplaceAll(gatewayURL, "${CID}", cidString)
	httpEndpoint := fmt.Sprintf("%s%s", r.config.GetNFTAssetNodeAddress(), modifiedGatewayURL)

	httpClient, err := http.NewRequest("GET", httpEndpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to setup get request: %w", err)
	}
	httpClient.Header.Add("Content-Type", "application/json")

	r.logger.Debug("Submitting to HTTP JSON API",
		slog.String("url", httpEndpoint),
		slog.String("method", "GET"))

	client := &http.Client{}
	res, err := client.Do(httpClient)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected response code: %v", res.StatusCode)
	}

	// Prepare the response struct
	resp := &domain.NFTAsset{}

	// Determine the content type and filename
	resp.ContentType = res.Header.Get("Content-Type")
	contentDisposition := res.Header.Get("Content-Disposition")

	if strings.Contains(contentDisposition, "filename*=") {
		filenameParts := strings.Split(contentDisposition, "filename*=")
		if len(filenameParts) > 1 {
			resp.Filename = strings.Trim(filenameParts[1], "\"")
		}
	} else {
		resp.Filename = "default-filename"
	}

	// Read response body into a byte slice
	bodyContent, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Assign content and length
	resp.Content = bodyContent
	resp.ContentLength = uint64(len(bodyContent))

	r.logger.Debug("Fetched file content",
		slog.String("filename", resp.Filename),
		slog.String("content_type", resp.ContentType),
		slog.Uint64("content_length", resp.ContentLength),
	)

	return resp, nil
}
