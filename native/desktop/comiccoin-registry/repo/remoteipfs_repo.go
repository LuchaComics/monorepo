package repo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	pkgfilepath "path/filepath"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-registry/domain"
)

type RemoteIPFSRepo struct {
	logger        *slog.Logger
	remoteAddress string
	apiKey        string
}

const (
	versionURL = "/version"
	pinAddURL  = "/ipfs/pin-add"
	gatewayURL = "/ipfs/${CID}"
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

	httpEndpoint := fmt.Sprintf("%s%s", r.remoteAddress, versionURL)

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

func (r *RemoteIPFSRepo) PinAddViaFilepath(ctx context.Context, fullFilePath string) (string, error) {

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
	req, err := http.NewRequest("POST", fmt.Sprintf("%v%v", r.remoteAddress, pinAddURL), &b)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %v\n", err)
	}

	// Create a Bearer string by appending string access token
	var bearer = "JWT " + string(r.apiKey)

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

func (r *RemoteIPFSRepo) Get(ctx context.Context, cidString string) ([]byte, string, error) {
	log.Fatal("TODO: IMPL.")
	return nil, "", nil
}
