package ipfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"strings"

	ipfslauncher "github.com/bartmika/ipfs-daemon-launcher"
	path "github.com/ipfs/boxo/path"
	"github.com/ipfs/go-cid"
	ipfsFiles "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/kubo/client/rpc"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

type IPFSStorager interface {
	UploadContentFromString(ctx context.Context, fileContent string) (string, error)
	UploadContentFromMulipart(ctx context.Context, file multipart.File) (string, error)
	GetContent(ctx context.Context, cidString string) ([]byte, error)
	PinContent(ctx context.Context, cidString string) error
	UnpinContent(ctx context.Context, cidString string) error
	DeleteContent(ctx context.Context, cidString string) error
	Shutdown()
}

type ipfsStorager struct {
	ipfsDaemonLauncher *ipfslauncher.IpfsDaemonLauncher
	httpApi            *rpc.HttpApi
	logger             *slog.Logger
}

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	logger.Debug("ipfs storage adapter initializing...", appConf.IPFSNode.BinaryOperatingSystem, appConf.IPFSNode.BinaryCPUArchitecture)

	launcher, initErr := ipfslauncher.NewDaemonLauncher(
		ipfslauncher.WithOverrideDaemonWarmupDuration(3),
		ipfslauncher.WithContinousOperation(),
		ipfslauncher.WithOverrideBinaryOsAndArch(appConf.IPFSNode.BinaryOperatingSystem, appConf.IPFSNode.BinaryCPUArchitecture),
	)
	if initErr != nil {
		log.Fatalf("failed creating ipfs-launcher: %v", initErr)
	}
	if launcher == nil {
		log.Fatal("cannot return nil launcher")
	}

	if startErr := launcher.StartDaemonInBackground(); startErr != nil {
		log.Fatal(startErr)
	}

	logger.Debug("ipfs storage adapter setup ipfs node")

	httpClient := &http.Client{}
	httpApi, httpErr := rpc.NewURLApiWithClient(appConf.IPFSNode.GatewayRPCURL, httpClient)
	if httpErr != nil {
		log.Fatalf("failed loading ipfs daemon: %v\n", httpErr)
	}

	// Create our storage handler for IPFS.
	ipfsStorage := &ipfsStorager{
		ipfsDaemonLauncher: launcher,
		httpApi:            httpApi,
		logger:             logger,
	}

	// Try uploading a sample file to verify our ipfs adapter works.
	sampleCid, sampleErr := ipfsStorage.UploadContentFromString(context.Background(), "Hello world via `Collectibles Protective Services`!")
	if sampleErr != nil {
		log.Fatalf("failed loading ipfs daemon: %v\n", sampleErr)
	}
	logger.Debug("ipfs storage adapter successfully uploaded sample file", slog.String("cid", sampleCid))

	// DEVELOPERS NOTE:
	// You can startup another `ipfs` deamon outside this project's docker and
	// run the following code to verify you get our sample contents:
	//
	//     ipfs cat /ipfs/QmNU6Q311PUfFuczUTYWkB7vnd4fD41dMxHKUEEGUYKLce

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (s *ipfsStorager) UploadContentFromString(ctx context.Context, fileContent string) (string, error) {
	content := strings.NewReader(fileContent)
	cid, err := s.httpApi.Unixfs().Add(context.Background(), ipfsFiles.NewReaderFile(content))
	if err != nil {
		return "", fmt.Errorf("failed to upload content from string: %v", err)
	}

	// Remove "/ipfs/" prefix if present
	cidString := strings.TrimPrefix(cid.String(), "/ipfs/")

	s.logger.Debug("uploaded content to IPFS via string",
		slog.String("cid", cidString))

	return cidString, nil
}

func (s *ipfsStorager) UploadContentFromMulipart(ctx context.Context, file multipart.File) (string, error) {
	// Debug log the start of the upload process
	s.logger.Debug("starting to upload file to IPFS")

	// Ensure the file gets closed when the function ends
	defer file.Close()

	// Upload the file to IPFS
	res, err := s.httpApi.Unixfs().Add(ctx, ipfsFiles.NewReaderFile(file))
	if err != nil {
		return "", fmt.Errorf("failed to add file to IPFS: %v", err)
	}

	// Retrieve the CID (Content Identifier) for the uploaded file and
	// remove "/ipfs/" prefix if present
	cidString := strings.TrimPrefix(res.String(), "/ipfs/")

	// Debug log the CID of the uploaded file
	s.logger.Debug("file successfully uploaded to IPFS", slog.String("cid", cidString))
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

	// Attempt to get the file from IPFS using the path
	fileNode, err := s.httpApi.Unixfs().Get(ctx, ipfsPath)
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

	s.logger.Debug("successfully fetched content from IPFS", slog.String("cid", cidString))
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
	if err := impl.httpApi.Pin().Add(ctx, ipfsPath); err != nil {
		impl.logger.Error("failed to pin content to IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to pin content to IPFS: %v", err)
	}

	impl.logger.Debug("successfully pinned content to IPFS", slog.String("cid", cidString))
	return nil
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
	err = s.httpApi.Pin().Rm(ctx, ipfsPath)
	if err != nil {
		s.logger.Error("failed to unpin content from IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return fmt.Errorf("failed to unpin content from IPFS: %v", err)
	}

	s.logger.Debug("successfully unpinned content from IPFS", slog.String("cid", cidString))
	return nil
}

func (s *ipfsStorager) DeleteContent(ctx context.Context, cidString string) error {
	// To delete content from an IPFS node, you generally need to unpin the content first, and then run the garbage collector to remove unpinned data. However, directly controlling garbage collection isn't typically exposed through the HTTP API, so simply unpinning the content is the standard way to "delete" it from the node.
	return s.UnpinContent(ctx, cidString)
}

func (impl *ipfsStorager) Shutdown() {
	impl.ipfsDaemonLauncher.ShutdownDaemon()
}
