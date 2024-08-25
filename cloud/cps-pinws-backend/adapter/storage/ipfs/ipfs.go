package ipfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strings"
	"time"

	ipfswrapper "github.com/bartmika/ipfs-wrapper"
	"github.com/ipfs/go-cid"
	ipfsFiles "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/kubo/client/rpc"
	"github.com/ipfs/kubo/path"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

type IPFSStorager interface {
	UploadContentFromMulipart(ctx context.Context, file multipart.File) (string, error)
	GetContentByCID(ctx context.Context, cidString string) ([]byte, error)
	PinContent(ctx context.Context, cidString string) error
	Shutdown()
}

type ipfsStorager struct {
	ipfsWrapper *ipfswrapper.IpfsWrapper
	httpApi     *rpc.HttpApi
	logger      *slog.Logger
}

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	logger.Debug("ipfs storage adapter initializing...", appConf.IPFSNode.BinaryOperatingSystem, appConf.IPFSNode.BinaryCPUArchitecture)
	return &ipfsStorager{} //TODO: REMOVE WHEN READY.

	wrapper, initErr := ipfswrapper.NewWrapper(
		// ipfswrapper.WithOverrideDaemonWarmupDuration(10),
		ipfswrapper.WithContinousOperation(),
		ipfswrapper.WithOverrideBinaryOsAndArch(appConf.IPFSNode.BinaryOperatingSystem, appConf.IPFSNode.BinaryCPUArchitecture),
	)
	if initErr != nil {
		// log.Fatalf("failed creating ipfs-wrapper: %v", initErr)

		logger.Error("ipfs storage adapter failed creating ipfs-wrapper",
			slog.Any("err", initErr))
		return &ipfsStorager{}
	}
	if wrapper == nil {
		log.Fatal("cannot return nil wrapper")
	}

	if startErr := wrapper.StartDaemonInBackground(); startErr != nil {
		log.Fatal(startErr)
	}

	// Set an artificial delay to give time for the `ipfs` binary to load up.
	// This is dependent on your machine.
	time.Sleep(8 * time.Second)

	logger.Debug("ipfs storage adapter setup ipfs node")

	httpClient := &http.Client{}
	httpApi, err := rpc.NewURLApiWithClient(appConf.IPFSNode.GatewayRPCURL, httpClient)
	if err != nil {
		log.Fatalf("failed loading ipfs daemon: %v\n", err)
	}

	logger.Debug("ipfs storage adapter rpc connected successfully")

	content := strings.NewReader("Hello world via `Collectibles Protective Services`!")
	p, err := httpApi.Unixfs().Add(context.Background(), ipfsFiles.NewReaderFile(content))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Data successfully stored in IPFS: %v\n", p)

	logger.Debug("ipfs storage adapter confirmed working locally")

	// Create our storage handler for IPFS.
	ipfsStorage := &ipfsStorager{
		ipfsWrapper: wrapper,
		httpApi:     httpApi,
		logger:      logger,
	}

	// Return our ipfs storage handler.
	return ipfsStorage
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

	// Retrieve the CID (Content Identifier) for the uploaded file
	cid := res.String()

	// Debug log the CID of the uploaded file
	s.logger.Debug("file successfully uploaded to IPFS", slog.String("cid", cid))

	return cid, nil
}

func (s *ipfsStorager) GetContentByCID(ctx context.Context, cidString string) ([]byte, error) {
	s.logger.Debug("fetching content from IPFS", slog.String("cid", cidString))

	c, err := cid.Decode(cidString)
	if err != nil {
		s.logger.Error("failed to decode CID", slog.String("cid", cidString), slog.Any("error", err))
		return nil, fmt.Errorf("failed to decode CID: %v", err)
	}

	// Convert the CID to a path.Path
	ipfsPath := path.IpfsPath(c)

	// Attempt to get the file from IPFS using the path
	fileNode, err := s.httpApi.Unixfs().Get(ctx, ipfsPath)
	if err != nil {
		s.logger.Error("failed to fetch content from IPFS", slog.String("cid", cidString), slog.Any("error", err))
		return nil, fmt.Errorf("failed to fetch content from IPFS: %v", err)
	}

	// Convert the file node to a reader
	fileReader := ipfsFiles.ToFile(fileNode) // TODO: FIX THIS BUG
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
	return nil // TODO
}

func (impl *ipfsStorager) Shutdown() {
	impl.ipfsWrapper.ShutdownDaemon()
}
