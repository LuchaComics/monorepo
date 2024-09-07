package ipfs

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net"
	"strings"
	"time"

	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"

	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	ipfsFiles "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/kubo/client/rpc"
	ma "github.com/multiformats/go-multiaddr"
)

type IPFSStorager interface {
	UploadContentFromString(ctx context.Context, fileContent string) (string, error)
	UploadContentFromMulipart(ctx context.Context, file multipart.File) (string, error)
	UploadContentFromBytes(ctx context.Context, fileContent []byte) (string, error)
	GetContent(ctx context.Context, cidString string) ([]byte, error)
	PinContent(ctx context.Context, cidString string) error
	ListPins(ctx context.Context) ([]string, error)
	UnpinContent(ctx context.Context, cidString string) error
	Shutdown()
}

type ipfsStorager struct {
	api    *rpc.HttpApi
	logger *slog.Logger
}

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
		logger: logger,
		api:    api,
	}

	logger.Debug("connected to ipfs node")

	// Try uploading a sample file to verify our ipfs adapter works.
	sampleCid, sampleErr := ipfsStorage.UploadContentFromString(context.Background(), "Hello world via `Collectibles Protective Services`!")
	if sampleErr != nil {
		log.Fatalf("failed uploading sample: %v\n", sampleErr)
	}
	logger.Debug("ipfs storage adapter successfully uploaded sample file", slog.String("cid", sampleCid))

	cids, listErr := ipfsStorage.ListPins(context.Background())
	if listErr != nil {
		log.Fatalf("failed listing: %v\n", sampleErr)
	}
	logger.Debug("ipfs storage adapter listed successfully", slog.Any("cids", cids))

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (s *ipfsStorager) UploadContentFromString(ctx context.Context, fileContent string) (string, error) {
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

func (s *ipfsStorager) UploadContentFromMulipart(ctx context.Context, file multipart.File) (string, error) {
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

func (s *ipfsStorager) UploadContentFromBytes(ctx context.Context, fileContent []byte) (string, error) {
	content := bytes.NewReader(fileContent) // THIS IS WRONG PLEASE REPAIR
	cid, err := s.api.Unixfs().Add(context.Background(), ipfsFiles.NewReaderFile(content))
	if err != nil {
		return "", fmt.Errorf("failed to upload content from string: %v", err)
	}

	// Remove "/ipfs/" prefix if present
	cidString := strings.TrimPrefix(cid.String(), "/ipfs/")

	return cidString, nil
}

func (s *ipfsStorager) GetContent(ctx context.Context, cidString string) ([]byte, error) {
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

func (impl *ipfsStorager) PinContent(ctx context.Context, cidString string) error {
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

func (s *ipfsStorager) UnpinContent(ctx context.Context, cidString string) error {
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

func (impl *ipfsStorager) Shutdown() {
}
