package repo

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/ipfs/boxo/files"
	"github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/core/coreiface/options"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"
)

type IPFSRepo struct {
	logger              *slog.Logger
	api                 *rpc.HttpApi
	publicGatewayDomain string
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
func NewIPFSRepo(cfg *config.Config, logger *slog.Logger) domain.IPFSRepository {

	// Step 1: Define the remote IPFS server address (replace with your remote IPFS server address)
	ipfsAddress := fmt.Sprintf("/ip4/%s/tcp/%s", cfg.IPFS.LocalIP, cfg.IPFS.LocalPort) // Example: Replace with your remote IPFS server address

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

	return &IPFSRepo{
		logger:              logger,
		api:                 api,
		publicGatewayDomain: cfg.IPFS.PublicGatewayDomain,
	}
}

// ID returns the IPFS node's identity information
func (r *IPFSRepo) ID() (peer.ID, error) {
	keyAPI := r.api.Key()
	if keyAPI == nil {
		return "", fmt.Errorf("Failed getting key: %v", "does not exist")
	}
	selfKey, err := keyAPI.Self(context.Background())
	if err != nil {
		return "", fmt.Errorf("Failed getting self: %v", err)
	}
	if selfKey == nil {
		return "", fmt.Errorf("Failed getting self: %v", "does not exist")
	}
	return selfKey.ID(), nil
}

func (r *IPFSRepo) Add(fullFilePath string, shouldPin bool) (string, error) {
	unixfs := r.api.Unixfs()
	if unixfs == nil {
		return "", fmt.Errorf("Failed getting unix fs: %v", "does not exist")
	}

	// Open the file
	file, err := os.Open(fullFilePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Get the file stat
	stat, err := file.Stat()
	if err != nil {
		return "", err
	}

	// Create a reader file node
	node, err := files.NewReaderPathFile(fullFilePath, file, stat)
	if err != nil {
		return "", err
	}

	// We want to use the newest `CidVersion` in our update.
	opts := func(settings *options.UnixfsAddSettings) error {
		settings.CidVersion = 1
		settings.Pin = shouldPin
		return nil
	}

	pathRes, err := unixfs.Add(context.Background(), node, opts)
	if err != nil {
		return "", err
	}

	return strings.Replace(pathRes.String(), "/ipfs/", "", -1), nil
}

func (impl *IPFSRepo) Pin(cidString string) error {
	impl.logger.Debug("pinning content to IPFS", slog.String("cid", cidString))

	cid, err := cid.Decode(cidString)
	if err != nil {
		impl.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to decode CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.FromCid(cid)

	// Attempt to pin the content to the IPFS node using the CID
	if err := impl.api.Pin().Add(context.Background(), ipfsPath); err != nil {
		impl.logger.Error("failed to pin content to IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to pin content to IPFS: %v", err)
	}
	return nil
}

func (r *IPFSRepo) PinAdd(fullFilePath string) (string, error) {
	fileCID, err := r.Add(fullFilePath, false)
	if err != nil {
		return "", err
	}

	if err := r.Pin(fileCID); err != nil {
		return "", err
	}

	return fileCID, nil
}

// Cat retrieves the contents of a file from IPFS
func (s *IPFSRepo) Get(ctx context.Context, cidString string) ([]byte, string, error) {
	s.logger.Debug("fetching content from IPFS", slog.String("cid", cidString))

	cid, err := cid.Decode(cidString)
	if err != nil {
		s.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
		return nil, "", fmt.Errorf("failed to decode CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.FromCid(cid)

	// Add a timeout to prevent hanging requests.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Attempt to get the file from IPFS using the path
	fileNode, err := s.api.Unixfs().Get(ctx, ipfsPath)
	if err != nil {
		s.logger.Debug("Failed fetching from remote IPFS...",
			slog.String("cid", cidString),
			slog.Any("error", err))
		return s.getViaHTTPPublicGateway(ctx, cidString)

	}

	// Convert the file node to a reader
	fileReader := files.ToFile(fileNode)
	if fileReader == nil {
		s.logger.Error("failed to convert IPFS node to file reader", slog.String("cid", cidString))
		return nil, "", fmt.Errorf("failed to convert IPFS node to file reader")
	}

	// Read the content from the file reader
	content, err := io.ReadAll(fileReader)
	if err != nil {
		s.logger.Error("failed to read content from IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return nil, "", fmt.Errorf("failed to read content from IPFS: %v", err)
	}

	return content, http.DetectContentType(content), nil
}

func (s *IPFSRepo) getViaHTTPPublicGateway(ctx context.Context, cidString string) ([]byte, string, error) {
	uri := fmt.Sprintf("%v/ipfs/%v", s.publicGatewayDomain, cidString)

	s.logger.Debug("Fetching from public IPFS gateway... Please wait...",
		slog.String("cid", cidString))

	resp, err := http.Get(uri)
	if err != nil {
		s.logger.Error("Failed fetching metadata uri via http.",
			slog.Any("error", err))
		return nil, "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code error: %d", resp.StatusCode)
		s.logger.Error("Status code error",
			slog.Any("error", err))
		return nil, "", err
	}

	// Get the content type from the response header
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		err := fmt.Errorf("Content type not specified in response header")
		s.logger.Error("Content-Type error",
			slog.Any("error", err))
		return nil, "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.logger.Error("Failed read all.",
			slog.Any("error", err))
		return nil, "", err
	}

	s.logger.Debug("Successfully fetching from public IPFS gateway.",
		slog.String("cid", cidString))

	return body, contentType, nil
}
