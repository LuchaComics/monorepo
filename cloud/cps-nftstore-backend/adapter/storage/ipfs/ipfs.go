package ipfs

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"mime/multipart"

	ipfscliwrapper "github.com/bartmika/ipfs-cli-wrapper"

	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
)

type IPFSStorager interface {
	AddFileContentFromMulipartFile(ctx context.Context, file multipart.File) (string, error)
	AddFileContent(ctx context.Context, fileContent []byte) (string, error)
	AddFileContentAndPin(ctx context.Context, fileContent []byte) (string, error)
	GetContent(ctx context.Context, cidString string) ([]byte, error)
	PinContent(ctx context.Context, cidString string) error
	ListPins(ctx context.Context) ([]string, error)
	UnpinContent(ctx context.Context, cidString string) error
	DeleteContent(ctx context.Context, cidString string) error
	Id(ctx context.Context) (*IpfsNodeInfo, error)
	Shutdown()
}

type ipfsStorager struct {
	ipfsBinFilepath string
	ipfsCliWrapper  ipfscliwrapper.IpfsCliWrapper
	logger          *slog.Logger
}

type IpfsNodeInfo struct {
	ID              string   `json:"ID"`
	PublicKey       string   `json:"PublicKey"`
	Addresses       []string `json:"Addresses"`
	AgentVersion    string   `json:"AgentVersion"`
	ProtocolVersion string   `json:"ProtocolVersion"`
}

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	logger.Debug("ipfs storage adapter initializing...",
		slog.String("os", appConf.IPFSNode.BinaryOperatingSystem),
		slog.String("arch", appConf.IPFSNode.BinaryCPUArchitecture),
	)

	launcher, initErr := ipfscliwrapper.NewWrapper(
		ipfscliwrapper.WithOverrideDaemonInitialWarmupDuration(25), // Wait 25 seconds for IPFS to startup for the first time. This is dependent on your machine.
		ipfscliwrapper.WithContinousOperation(),
		ipfscliwrapper.WithOverrideBinaryOsAndArch(appConf.IPFSNode.BinaryOperatingSystem, appConf.IPFSNode.BinaryCPUArchitecture),
		ipfscliwrapper.WithDenylist("badbits.deny", "https://badbits.dwebops.pub/badbits.deny"), // Taken from https://github.com/ipfs/kubo/blob/master/docs/content-blocking.md#denylist-file-format
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

	// Create our storage handler for IPFS.
	ipfsStorage := &ipfsStorager{
		ipfsBinFilepath: "./bin/kubo/ipfs",
		ipfsCliWrapper:  launcher,
		logger:          logger,
	}

	// DEVELOPERS NOTE:
	// The following code will create a sample file so we can verify our
	// ipfs adapter is confirmed working.
	sampleBytes := "Hello world via `Collectibles Protective Services` NFT store!"
	sampleCid, addFileContentErr := ipfsStorage.AddFileContentAndPin(context.Background(), []byte(sampleBytes))
	if addFileContentErr != nil {
		log.Fatalf("failed adding sample file with error: %v", addFileContentErr)
	}
	logger.Debug("ipfs storage successfully saved sample file",
		slog.String("sample_cid", sampleCid))

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (impl *ipfsStorager) AddFileContentFromMulipartFile(ctx context.Context, file multipart.File) (string, error) {
	fileContent, err := convertFileToBytes(file)
	if err != nil {
		return "", fmt.Errorf("failed convert file to bytes array: %w", err)
	}
	return impl.AddFileContent(ctx, fileContent)
}

func (impl *ipfsStorager) AddFileContent(ctx context.Context, fileContent []byte) (string, error) {
	cid, addFileErr := impl.ipfsCliWrapper.AddFileContent(ctx, fileContent)
	if addFileErr != nil {
		impl.logger.Error("failed to save file locally",
			slog.Any("error", addFileErr))
		return "", fmt.Errorf("failed to save file locally: %v", addFileErr)
	}
	return cid, nil
}

func (impl *ipfsStorager) AddFileContentAndPin(ctx context.Context, fileContent []byte) (string, error) {
	impl.logger.Debug("adding file content to ipfs...")
	cid, addFileErr := impl.ipfsCliWrapper.AddFileContent(ctx, fileContent)
	if addFileErr != nil {
		impl.logger.Error("failed to save file locally",
			slog.Any("error", addFileErr))
		return "", fmt.Errorf("failed to save file locally: %v", addFileErr)
	}
	impl.logger.Debug("successfullly adding file content to ipfs")
	impl.logger.Debug("pinning file content to ipfs...")
	if pinErr := impl.ipfsCliWrapper.Pin(ctx, cid); pinErr != nil {
		impl.logger.Error("failed to pin local file content",
			slog.String("cid", cid),
			slog.Any("error", pinErr))
		return "", fmt.Errorf("failed to pin local file content: %v", pinErr)
	}
	impl.logger.Debug("successfullly pinned file content to ipfs")

	return cid, nil
}

func (impl *ipfsStorager) GetContent(ctx context.Context, cidString string) ([]byte, error) {
	impl.logger.Debug("fetching content from IPFS", slog.String("cid", cidString))
	content, catErr := impl.ipfsCliWrapper.Cat(ctx, cidString)
	if catErr != nil {
		impl.logger.Error("failed fetching content",
			slog.Any("error", catErr))
		return []byte{}, fmt.Errorf("failed fetching content: %v", catErr)
	}

	impl.logger.Debug("successfully fetched content from IPFS", slog.String("cid", cidString))
	return content, nil
}

func (impl *ipfsStorager) PinContent(ctx context.Context, cidString string) error {
	impl.logger.Debug("pinning content to IPFS", slog.String("cid", cidString))
	if pinErr := impl.ipfsCliWrapper.Pin(ctx, cidString); pinErr != nil {
		impl.logger.Error("failed to pin locally",
			slog.String("cid", cidString),
			slog.Any("error", pinErr))
		return fmt.Errorf("failed to pin locally: %v", pinErr)
	}
	impl.logger.Debug("successfully pinned content to IPFS", slog.String("cid", cidString))
	return nil
}

func (impl *ipfsStorager) ListPins(ctx context.Context) ([]string, error) {
	cids, err := impl.ipfsCliWrapper.ListPins(ctx)
	if err != nil {
		impl.logger.Error("failed listing pins",
			slog.Any("error", err))
		return []string{}, fmt.Errorf("failed listing pins: %v", err)
	}
	return cids, nil
}

func (impl *ipfsStorager) UnpinContent(ctx context.Context, cidString string) error {
	impl.logger.Debug("unpinning content from IPFS", slog.String("cid", cidString))
	unpinErr := impl.ipfsCliWrapper.Unpin(ctx, cidString)
	if unpinErr != nil {
		impl.logger.Error("failed to unpin content",
			slog.String("cid", cidString),
			slog.Any("error", unpinErr))
		return fmt.Errorf("failed to unpin content: %v", unpinErr)
	}
	impl.logger.Debug("successfully unpinned content from IPFS", slog.String("cid", cidString))
	return nil
}

func (s *ipfsStorager) DeleteContent(ctx context.Context, cidString string) error {
	// To delete content from an IPFS node, you generally need to unpin the content first, and then run the garbage collector to remove unpinned data. However, directly controlling garbage collection isn't typically exposed through the HTTP API, so simply unpinning the content is the standard way to "delete" it from the node.
	return s.UnpinContent(ctx, cidString)
}

func convertFileToBytes(file multipart.File) ([]byte, error) {
	// Use io.ReadAll to read the entire content of the file into a byte slice
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return fileBytes, nil
}

func (impl *ipfsStorager) Id(ctx context.Context) (*IpfsNodeInfo, error) {
	info, err := impl.ipfsCliWrapper.Id(ctx)
	if err != nil {
		impl.logger.Error("failed getting ipfds node info",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting ipfds node info: %v", err)
	}
	return &IpfsNodeInfo{
		ID:              info.ID,
		PublicKey:       info.PublicKey,
		Addresses:       info.Addresses,
		AgentVersion:    info.AgentVersion,
		ProtocolVersion: info.ProtocolVersion,
	}, nil
}

func (impl *ipfsStorager) Shutdown() {
	impl.ipfsCliWrapper.ShutdownDaemon()
}
