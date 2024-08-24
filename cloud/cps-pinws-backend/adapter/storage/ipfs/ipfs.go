package ipfs

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"net"
	"os"
	"strings"

	ipfsFiles "github.com/ipfs/go-ipfs-files"
	"github.com/ipfs/kubo/client/rpc"
	ma "github.com/multiformats/go-multiaddr"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

type IPFSStorager interface {
	UploadContentFromFilepath(ctx context.Context, filepath string) (string, error)
	GetContentByCID(ctx context.Context, cidString string) ([]byte, error)
	PinContent(ctx context.Context, cidString string) error
}

type ipfsStorager struct {
	httpApi *rpc.HttpApi
	logger  *slog.Logger
}

func NewStorage(appConf *c.Conf, logger *slog.Logger) IPFSStorager {
	logger.Debug("ipfs storage adapter initializing...")

	var ipv4Str string
	ips, _ := net.LookupIP(appConf.IPFSNode.DomainOrIPAddress)
	for _, ip := range ips {
		if ipv4 := ip.To4(); ipv4 != nil {
			ipv4Str = fmt.Sprintf("%v", ipv4)
			logger.Debug("ipfs storage adapter will connect to active running ipfs node instance",
				slog.String("domain", appConf.IPFSNode.DomainOrIPAddress),
				slog.String("ip_address", ipv4Str))
		}
	}
	if ipv4Str == "" {
		ipv4Str = appConf.IPFSNode.DomainOrIPAddress
		logger.Debug("ipfs storage adapter will connect to active running node instance", slog.String("ip_address", ipv4Str))
	}

	addr, err := ma.NewMultiaddr("/ip4/" + ipv4Str + "/tcp/5001")
	if err != nil {
		log.Fatalf("failed make address: %s", err)
	}

	// "Connect" to local node
	httpApi, err := rpc.NewApi(addr)
	if err != nil {
		log.Fatalf("failed to connect to node: %s", err)
	}

	logger.Debug("ipfs storage adapter rpc connected successfully")

	content := strings.NewReader("Hello world via `Collectibles Protective Services`!")
	p, err := httpApi.Unixfs().Add(context.Background(), ipfsFiles.NewReaderFile(content))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("Data successfully stored in IPFS: %v\n", p)

	// // Pin a given file by its CID
	// ctx := context.Background()
	// cid := "bafkreidtuosuw37f5xmn65b3ksdiikajy7pwjjslzj2lxxz2vc4wdy3zku"
	// p := path.New(cid)
	// err = node.Pin().Add(ctx, p)
	// if err != nil {
	// 	log.Fatalf("failed to pin: %s", err)
	// }

	// Create our storage handler for IPFS.
	ipfsStorage := &ipfsStorager{
		httpApi: httpApi,
		logger:  logger,
	}

	// Return our ipfs storage handler.
	return ipfsStorage
}

func (impl *ipfsStorager) UploadContentFromFilepath(ctx context.Context, filepath string) (string, error) {
	return "", nil
}

func (impl *ipfsStorager) GetContentByCID(ctx context.Context, cidString string) ([]byte, error) {
	return []byte{}, nil
}

func (impl *ipfsStorager) PinContent(ctx context.Context, cidString string) error {
	return nil
}
