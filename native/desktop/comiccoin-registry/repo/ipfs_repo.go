package repo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
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

// AddResponse represents the response from the /api/v0/add endpoint
type AddResponse struct {
	Bytes      int64  `json:"Bytes"`
	Hash       string `json:"Hash"`
	Mode       string `json:"Mode"`
	Mtime      int64  `json:"Mtime"`
	MtimeNsecs int    `json:"MtimeNsecs"`
	Name       string `json:"Name"`
	Size       string `json:"Size"`
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

// // Version returns the IPFS node's version information
// func (n *IPFSNode) Version() (*Version, error) {
//     resp, err := http.Get(n.APIEndpoint + "/api/v0/version")
//     if err != nil {
//         return nil, err
//     }
//     defer resp.Body.Close()
//
//     var version Version
//     err = json.NewDecoder(resp.Body).Decode(&version)
//     if err != nil {
//         return nil, err
//     }
//
//     return &version, nil
// }

func (r *IPFSRepo) AddAndPinSingleFileFromLocalFileSystem(fullFilePath string) (*AddResponse, error) {
	res, err := r.AddAndPinFromLocalFileSystem(fullFilePath)
	if err != nil {
		return nil, err
	}
	return res[0], nil
}

// Add adds a new file or directory to IPFS
func (r *IPFSRepo) AddAndPinFromLocalFileSystem(fullFilePath string) ([]*AddResponse, error) {

	// Open the file
	file, err := os.Open(fullFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Get the base name of the file
	filename := filepath.Base(fullFilePath)

	// Create a multipart/form-data request body
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	// Copy the file contents to the multipart request body
	_, err = io.Copy(part, file)
	if err != nil {
		return nil, err
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	// Please see: https://docs.ipfs.tech/reference/kubo/rpc/#api-v0-add
	// Summary:
	// - pin: We want to automatically pin along with our save.
	// - cid-version=1 --> We want to use the latest cid system.
	// - wrap-with-directory --> Include filename! Important!
	params := "?pin=true&cid-version=1&wrap-with-directory=false"

	// Set the request headers
	req, err := http.NewRequest("POST", r.APIEndpoint+"/api/v0/add"+params, buf)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Make the Add API call
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Create a JSON decoder
	decoder := json.NewDecoder(resp.Body)

	// Create a slice to store the AddResponse objects
	var addResponses []*AddResponse

	// Decode the JSON response
	for decoder.More() {
		var addResponse AddResponse
		err = decoder.Decode(&addResponse)
		if err != nil {
			return nil, err
		}
		addResponses = append(addResponses, &addResponse)
	}

	// Print the Add API response
	for _, addResponse := range addResponses {
		fmt.Printf("Added file with CID: %s\n", addResponse.Hash)
		fmt.Printf("File name: %s\n", addResponse.Name)
		fmt.Printf("File size: %s\n", addResponse.Size)
	}

	return addResponses, nil
}

// Cat retrieves the contents of a file from IPFS
func (r *IPFSRepo) Cat(cid string) ([]byte, string, uint64, error) {
	url := fmt.Sprintf("%v/api/v0/get?arg=%s", r.APIEndpoint, cid)

	// Create a new HTTP request
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		r.logger.Error("error new request",
			slog.Any("error", err))
		return nil, "", 0, err
	}

	// Set the User-Agent header
	req.Header.Set("User-Agent", "My IPFS Client")

	// Send the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		r.logger.Error("failed executing http request",
			slog.Any("error", err))
		return nil, "", 0, err
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != 200 {
		return nil, "", 0, fmt.Errorf("invalid response code: %d", resp.StatusCode)
	}

	// Get the content type
	contentType := resp.Header.Get("Content-Type")

	// If the content type is not specified, try to determine it based on the file extension
	if contentType == "" {
		// Try to determine the content type based on the Content-Transfer-Encoding header
		encoding := resp.Header.Get("Content-Transfer-Encoding")
		if encoding == "binary" {
			// If the encoding is binary, try to determine the content type based on the file extension
			// ...
		}
		log.Println("------->", encoding)

		// Parse the CID and extract the file extension
		parts := strings.Split(cid, "/")
		filename := parts[len(parts)-1]
		extension := strings.Split(filename, ".")[1]

		// Determine the content type based on the file extension
		switch extension {
		case "png":
			contentType = "image/png"
		case "jpg":
			contentType = "image/jpeg"
		case "gif":
			contentType = "image/gif"
			// ...
		}
	}

	// Get the content length
	contentLength, err := strconv.ParseUint(resp.Header.Get("X-Content-Length"), 10, 64)
	if err != nil {
		r.logger.Error("failed parsing to int",
			slog.Any("header", resp.Header),
			slog.Any("error", err))
		return nil, "", 0, err
	}

	// Read the response body
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		r.logger.Error("error reading all data",
			slog.Any("error", err))
		return nil, "", 0, err
	}

	r.logger.Debug("Get done",
		slog.Any("header", resp.Header))

	return data, contentType, contentLength, nil
}
