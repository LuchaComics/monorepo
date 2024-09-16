package ipfs

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	ipfsFiles "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
	ma "github.com/multiformats/go-multiaddr"

	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
)

type IPFSStorager interface {
	// Uploads content to IPFS from different sources
	UploadString(ctx context.Context, content string) (cid string, err error)
	UploadBytes(ctx context.Context, content []byte) (cid string, err error)
	UploadMultipart(ctx context.Context, file multipart.File) (cid string, err error)

	// Uploads content to IPFS within a specified directory
	UploadStringToDir(ctx context.Context, content, fileName, dirName string) (dirCid, fileCid string, err error)
	UploadBytesToDir(ctx context.Context, content []byte, fileName, dirName string) (dirCid, fileCid string, err error)
	UploadMultipartToDir(ctx context.Context, file multipart.File, fileName, dirName string) (dirCid, fileCid string, err error)

	// Retrieves content from IPFS
	Get(ctx context.Context, cid string) ([]byte, error)

	// Manages pinning of content in IPFS
	Pin(ctx context.Context, cid string) error
	Unpin(ctx context.Context, cid string) error
	ListPins(ctx context.Context) ([]string, error)

	// IPNS-related methods
	GenerateKey(ctx context.Context, keyName string) (ipnsName string, err error)
	PublishToIPNS(ctx context.Context, keyName, dirCid string) (ipnsName string, err error)

	CheckIfKeyNameExists(ctx context.Context, keyName string) (bool, error)

	// Shutdown the IPFS service
	Shutdown()
}

type ipfsStorager struct {
	api    *rpc.HttpApi
	logger *slog.Logger
	apiUrl string
}

// DEVELOPERS NOTE:
// Useful links as follows:
// - https://github.com/ipfs/kubo/tree/master/client/rpc
// - https://pkg.go.dev/github.com/ipfs/kubo/client/rpc
//

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	logger.Debug("connecting to ipfs node...")

	// The following block of code will be used to resolve the dns of our
	// other docker container to get the `ipfs-node` ip address.
	var ipfsIP string
	ips, err := net.LookupIP("ipfs-node")
	if err != nil {
		log.Fatalf("failed to lookup dns record: %v", err)
	}
	for _, ip := range ips {
		ipfsIP = ip.String()
		break
	}

	// Step 1: Define the remote IPFS server address (replace with your remote IPFS server address)
	ipfsAddress := fmt.Sprintf("/ip4/%s/tcp/5001", ipfsIP) // Example: Replace with your remote IPFS server address

	// Step 2: Create a Multiaddr using the remote IPFS address
	multiaddr, err := ma.NewMultiaddr(ipfsAddress)
	if err != nil {
		log.Fatalf("failed to create multiaddr: %v", err)
	}

	// Step 3: Create a new IPFS HTTP API client using the remote server address
	api, err := rpc.NewApi(multiaddr)
	if err != nil {
		log.Fatalf("failed to create IPFS HTTP API client: %v", err)
	}

	// Create our storage handler for IPFS.
	ipfsStorage := &ipfsStorager{
		apiUrl: ipfsIP,
		logger: logger,
		api:    api,
	}

	logger.Debug("connected to ipfs node")

	// Try uploading a sample file to verify our ipfs adapter works.
	sampleDirCid, sampleFileCid, sampleErr := ipfsStorage.UploadStringToDir(context.Background(), "Hello world via `Collectibles Protective Services`!", "sample.txt", "sampledir")
	if sampleErr != nil {
		log.Fatalf("failed uploading sample: %v\n", sampleErr)
	}
	logger.Debug("ipfs storage adapter successfully uploaded sample file",
		slog.String("dir_cid", sampleDirCid),
		slog.String("file_cid", sampleFileCid))

	cids, listErr := ipfsStorage.ListPins(context.Background())
	if listErr != nil {
		log.Fatalf("failed listing: %v\n", sampleErr)
	}
	logger.Debug("ipfs storage adapter listed successfully",
		slog.Any("cids", cids),
		slog.Int("cids_len", len(cids)))

	// Return our ipfs storage handler.
	return ipfsStorage
}

type ipfsApiAddResponse struct {
	Name string `json:"Name"`
	Hash string `json:"Hash"`
	Size string `json:"Size"`
}

func (s *ipfsStorager) UploadString(ctx context.Context, fileContent string) (string, error) {
	// Create a reader for the file content
	content := strings.NewReader(fileContent)

	// Upload the content to IPFS
	cid, err := s.api.Unixfs().Add(ctx, ipfsFiles.NewReaderFile(content))
	if err != nil {
		return "", fmt.Errorf("failed to upload content from string: %v", err)
	}

	// Remove "/ipfs/" prefix if present
	cidString := strings.TrimPrefix(cid.String(), "/ipfs/")

	return cidString, nil
}

func (s *ipfsStorager) UploadStringToDir(ctx context.Context, fileContent string, fileName string, directoryName string) (string, string, error) {
	// Debug log the start of the upload process
	s.logger.Debug("starting to upload file to IPFS")

	// Create an in-memory reader for the file content
	fileReader := strings.NewReader(fileContent)

	// Prepare the multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form file field in the writer and include the directory structure
	fileField, err := writer.CreateFormFile(directoryName+"/"+fileName, fileName)
	if err != nil {
		s.logger.Error("x1")
		return "", "", fmt.Errorf("failed to create form file field: %v", err)
	}

	// Copy the file content into the multipart form file field
	_, err = io.Copy(fileField, fileReader)
	if err != nil {
		s.logger.Error("x2")
		return "", "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to complete the form
	err = writer.Close()
	if err != nil {
		s.logger.Error("x3")
		return "", "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Make the request to IPFS API to add files wrapped within a directory
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:5001/api/v0/add?wrap-with-directory=true&cid-version=1", s.apiUrl), body)
	if err != nil {
		s.logger.Error("x4")
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type header for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("x5")
		return "", "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle error
	}

	// Remove the escape characters
	jsonString := strings.ReplaceAll(string(jsonData), "\\", "")

	// Split the string by newlines
	jsonObjects := strings.Split(strings.TrimSpace(jsonString), "\n")

	// Slice to store the parsed structs
	var responses []ipfsApiAddResponse

	// Loop through the JSON objects and unmarshal them
	for _, jsonObject := range jsonObjects {
		if len(strings.TrimSpace(jsonObject)) == 0 {
			continue // skip empty lines
		}

		var resp ipfsApiAddResponse
		err := json.Unmarshal([]byte(jsonObject), &resp)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON object: %v\n", err)
			continue
		}

		// Append to the slice
		responses = append(responses, resp)
	}

	var fileCid string
	var dirCid string

	// Print the parsed structs
	for _, r := range responses {
		if r.Name == fileName {
			fileCid = r.Hash
		} else {
			dirCid = r.Hash

		}
	}

	return dirCid, fileCid, nil
}

func (s *ipfsStorager) UploadMultipart(ctx context.Context, file multipart.File) (string, error) {
	// Debug log the start of the upload process
	s.logger.Debug("starting to upload file to IPFS")

	// Ensure the file gets closed when the function ends
	defer file.Close()

	// Upload the file to IPFS
	res, err := s.api.Unixfs().Add(ctx, ipfsFiles.NewReaderFile(file))
	if err != nil {
		return "", fmt.Errorf("failed to add file to IPFS: %v", err)
	}

	// Retrieve the CID (Content Identifier) for the uploaded file and
	// remove "/ipfs/" prefix if present
	cidString := strings.TrimPrefix(res.String(), "/ipfs/")

	return cidString, nil
}

func (s *ipfsStorager) UploadMultipartToDir(ctx context.Context, file multipart.File, fileName string, directoryName string) (string, string, error) {
	// Debug log the start of the upload process
	s.logger.Debug("starting to upload file to IPFS")

	// Create an in-memory reader for the file content
	fileReader := ipfsFiles.NewReaderFile(file)

	// Prepare the multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form file field in the writer and include the directory structure
	fileField, err := writer.CreateFormFile(directoryName+"/"+fileName, fileName)
	if err != nil {
		s.logger.Error("x1")
		return "", "", fmt.Errorf("failed to create form file field: %v", err)
	}

	// Copy the file content into the multipart form file field
	_, err = io.Copy(fileField, fileReader)
	if err != nil {
		s.logger.Error("x2")
		return "", "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to complete the form
	err = writer.Close()
	if err != nil {
		s.logger.Error("x3")
		return "", "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Make the request to IPFS API to add files wrapped within a directory
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:5001/api/v0/add?wrap-with-directory=true&cid-version=1", s.apiUrl), body)
	if err != nil {
		s.logger.Error("x4")
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type header for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("x5")
		return "", "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle error
	}

	// Remove the escape characters
	jsonString := strings.ReplaceAll(string(jsonData), "\\", "")

	// Split the string by newlines
	jsonObjects := strings.Split(strings.TrimSpace(jsonString), "\n")

	// Slice to store the parsed structs
	var responses []ipfsApiAddResponse

	// Loop through the JSON objects and unmarshal them
	for _, jsonObject := range jsonObjects {
		if len(strings.TrimSpace(jsonObject)) == 0 {
			continue // skip empty lines
		}

		var resp ipfsApiAddResponse
		err := json.Unmarshal([]byte(jsonObject), &resp)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON object: %v\n", err)
			continue
		}

		// Append to the slice
		responses = append(responses, resp)
	}

	var fileCid string
	var dirCid string

	// Print the parsed structs
	for _, r := range responses {
		if r.Name == fileName {
			fileCid = r.Hash
		} else {
			dirCid = r.Hash

		}
	}

	return dirCid, fileCid, nil
}

func (s *ipfsStorager) UploadBytes(ctx context.Context, fileContent []byte) (string, error) {
	content := bytes.NewReader(fileContent) // THIS IS WRONG PLEASE REPAIR
	cid, err := s.api.Unixfs().Add(context.Background(), ipfsFiles.NewReaderFile(content))
	if err != nil {
		return "", fmt.Errorf("failed to upload content from string: %v", err)
	}

	// Remove "/ipfs/" prefix if present
	cidString := strings.TrimPrefix(cid.String(), "/ipfs/")

	return cidString, nil
}

func (s *ipfsStorager) UploadBytesToDir(ctx context.Context, fileContent []byte, fileName string, directoryName string) (string, string, error) {
	// Debug log the start of the upload process
	s.logger.Debug("starting to upload file to IPFS")

	// Create an in-memory reader for the file content
	content := bytes.NewReader(fileContent)

	// Prepare the multipart form data
	body := new(bytes.Buffer)
	writer := multipart.NewWriter(body)

	// Create a form file field in the writer and include the directory structure
	fileField, err := writer.CreateFormFile(directoryName+"/"+fileName, fileName)
	if err != nil {
		s.logger.Error("x1")
		return "", "", fmt.Errorf("failed to create form file field: %v", err)
	}

	// Copy the file content into the multipart form file field
	_, err = io.Copy(fileField, content)
	if err != nil {
		s.logger.Error("x2")
		return "", "", fmt.Errorf("failed to copy file content: %v", err)
	}

	// Close the writer to complete the form
	err = writer.Close()
	if err != nil {
		s.logger.Error("x3")
		return "", "", fmt.Errorf("failed to close writer: %v", err)
	}

	// Make the request to IPFS API to add files wrapped within a directory
	req, err := http.NewRequest("POST", fmt.Sprintf("http://%s:5001/api/v0/add?wrap-with-directory=true&cid-version=1", s.apiUrl), body)
	if err != nil {
		s.logger.Error("x4")
		return "", "", fmt.Errorf("failed to create request: %v", err)
	}

	// Set the content type header for multipart form data
	req.Header.Set("Content-Type", writer.FormDataContentType())

	// Perform the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		s.logger.Error("x5")
		return "", "", fmt.Errorf("failed to perform request: %v", err)
	}
	defer resp.Body.Close()

	jsonData, err := io.ReadAll(resp.Body)
	if err != nil {
		// Handle error
	}

	// Remove the escape characters
	jsonString := strings.ReplaceAll(string(jsonData), "\\", "")

	// Split the string by newlines
	jsonObjects := strings.Split(strings.TrimSpace(jsonString), "\n")

	// Slice to store the parsed structs
	var responses []ipfsApiAddResponse

	// Loop through the JSON objects and unmarshal them
	for _, jsonObject := range jsonObjects {
		if len(strings.TrimSpace(jsonObject)) == 0 {
			continue // skip empty lines
		}

		var resp ipfsApiAddResponse
		err := json.Unmarshal([]byte(jsonObject), &resp)
		if err != nil {
			fmt.Printf("Failed to unmarshal JSON object: %v\n", err)
			continue
		}

		// Append to the slice
		responses = append(responses, resp)
	}

	var fileCid string
	var dirCid string

	// Print the parsed structs
	for _, r := range responses {
		if r.Name == fileName {
			fileCid = r.Hash
		} else {
			dirCid = r.Hash

		}
	}

	return dirCid, fileCid, nil
}

func (s *ipfsStorager) Get(ctx context.Context, cidString string) ([]byte, error) {
	s.logger.Debug("fetching content from IPFS", slog.String("cid", cidString))

	cid, err := cid.Decode(cidString)
	if err != nil {
		s.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
		return nil, fmt.Errorf("failed to decode CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.FromCid(cid)

	// Add a timeout to prevent hanging requests.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to get the file from IPFS using the path
	fileNode, err := s.api.Unixfs().Get(ctx, ipfsPath)
	if err != nil {
		s.logger.Error("failed to fetch content from IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return nil, fmt.Errorf("failed to fetch content from IPFS: %v", err)
	}

	// Convert the file node to a reader
	fileReader := ipfsFiles.ToFile(fileNode)
	if fileReader == nil {
		s.logger.Error("failed to convert IPFS node to file reader", slog.String("cid", cidString))
		return nil, fmt.Errorf("failed to convert IPFS node to file reader")
	}

	// Read the content from the file reader
	content, err := io.ReadAll(fileReader)
	if err != nil {
		s.logger.Error("failed to read content from IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return nil, fmt.Errorf("failed to read content from IPFS: %v", err)
	}

	return content, nil
}

func (impl *ipfsStorager) Pin(ctx context.Context, cidString string) error {
	impl.logger.Debug("pinning content to IPFS", slog.String("cid", cidString))

	cid, err := cid.Decode(cidString)
	if err != nil {
		impl.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to decode CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.FromCid(cid)

	// Attempt to pin the content to the IPFS node using the CID
	if err := impl.api.Pin().Add(ctx, ipfsPath); err != nil {
		impl.logger.Error("failed to pin content to IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to pin content to IPFS: %v", err)
	}
	return nil
}

func (impl *ipfsStorager) ListPins(ctx context.Context) ([]string, error) {
	// Fetch the pinned items channel
	pinCh, err := impl.api.Pin().Ls(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to list pinned items: %v", err)
	}

	// Prepare a slice to hold the pinned CIDs
	pinnedCIDs := make([]string, 0)

	// Read from the channel until it is closed
	for pin := range pinCh {

		pinnedCIDs = append(pinnedCIDs, pin.Path().RootCid().String())
	}

	return pinnedCIDs, nil
}

func (s *ipfsStorager) Unpin(ctx context.Context, cidString string) error {
	s.logger.Debug("unpinning content from IPFS", slog.String("cid", cidString))

	// Decode the CID string into a CID object
	c, err := cid.Decode(cidString)
	if err != nil {
		s.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to decode CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.FromCid(c)

	// Use the IPFS HTTP API to unpin the content
	err = s.api.Pin().Rm(ctx, ipfsPath)
	if err != nil {
		s.logger.Error("failed to unpin content from IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to unpin content from IPFS: %v", err)
	}

	return nil
}

func (s *ipfsStorager) GenerateKey(ctx context.Context, keyName string) (string, error) {
	key, err := s.api.Key().Generate(ctx, keyName)
	if err != nil {
		return "", fmt.Errorf("failed to list pinned items: %v", err)
	}
	if key == nil {
		return "", fmt.Errorf("nil key")
	}
	// fmt.Println("name:", key.Name())
	// fmt.Println("path:", key.Path())
	// fmt.Println("id:", key.ID())
	return key.Path().String(), nil
}

func (s *ipfsStorager) PublishToIPNS(ctx context.Context, keyName string, directoryCidString string) (string, error) {
	directoryCid, err := cid.Decode(directoryCidString)
	if err != nil {
		s.logger.Error("failed to decode", slog.String("dir_cid", directoryCidString), slog.Any("error", err))
		return "", fmt.Errorf("failed to decode directory CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.FromCid(directoryCid)

	// Define the option functions
	optKey := func(settings *options.NamePublishSettings) error {
		settings.Key = keyName
		return nil
	}
	optValidTime := func(settings *options.NamePublishSettings) error {
		settings.ValidTime = 24 * time.Hour
		return nil
	}
	optCompatibleWithV1 := func(settings *options.NamePublishSettings) error {
		settings.CompatibleWithV1 = true
		return nil
	}
	optAllowOffline := func(settings *options.NamePublishSettings) error {
		settings.AllowOffline = true
		return nil
	}

	// Call the IPFS API with the option functions directly
	res, err := s.api.Name().Publish(ctx, ipfsPath, optKey, optValidTime, optCompatibleWithV1, optAllowOffline)
	if err != nil {
		return "", fmt.Errorf("failed to publish to IPNS: %v", err)
	}

	// // Log and return the result
	// log.Println(res)
	// log.Println(res.RoutingKey())
	// log.Println(res.Cid())
	// log.Println(res.Peer())
	// log.Println(res.String())
	// log.Println(res.AsPath())
	// log.Println(res.String())

	return res.String(), nil
}

func (s *ipfsStorager) CheckIfKeyNameExists(ctx context.Context, keyName string) (bool, error) {
	keyAPI := s.api.Key()

	keys, err := keyAPI.List(ctx)
	if err != nil {
		return false, err
	}
	for _, key := range keys {
		if key.Name() == keyName {
			return true, nil
		}
	}

	// return res.String(), nil
	return false, nil
}

func (impl *ipfsStorager) Shutdown() {
}
