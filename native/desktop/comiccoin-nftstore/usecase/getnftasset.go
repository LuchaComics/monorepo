package usecase

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/constants"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type GetNFTAssetUseCase struct {
	logger *slog.Logger
}

func NewGetNFTAssetUseCase(logger *slog.Logger) *GetNFTAssetUseCase {
	return &GetNFTAssetUseCase{logger}
}

type NFTAssetIDO struct {
	Filename      string
	FileExtension string
	ContentType   string
	ContentLength int64
	Content       *bytes.Reader
}

func (uc *GetNFTAssetUseCase) Execute(uri string) (*NFTAssetIDO, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if uri == "" {
		e["asset_uri"] = "missing value"
	}
	if len(e) != 0 {
		uc.logger.Warn("Failed validating get remote file uri.",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Check protocol and fetch using the particular protocol.
	//

	if strings.Contains(uri, "https://") {
		return uc.fetchViaHTTPS(uri)
	} else if strings.Contains(uri, "ipfs://") {
		return uc.fetchViaIPFS(uri)
	}

	return nil, fmt.Errorf("Unsupported protocol in URI, only accepted options are: %v", "https and ipfs")
}

func (uc *GetNFTAssetUseCase) fetchViaHTTPS(uri string) (*NFTAssetIDO, error) {
	uc.logger.Debug("fetching from uri via HTTPS protocol.",
		slog.String("uri", uri))

	resp, err := http.Get(uri)
	if err != nil {
		uc.logger.Error("Failed fetching metadata uri via http.",
			slog.Any("error", err))
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err := fmt.Errorf("status code error: %d", resp.StatusCode)
		uc.logger.Error("Status code error",
			slog.Any("error", err))
		return nil, err
	}

	// Get the content type from the response header
	contentType := resp.Header.Get("Content-Type")
	if contentType == "" {
		err := fmt.Errorf("Content type not specified in response header")
		uc.logger.Error("Content-Type error",
			slog.Any("error", err))
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		uc.logger.Error("Failed read all.",
			slog.Any("error", err))
		return nil, err
	}

	// Create a file name based on the content type
	ext := ""
	switch contentType {
	case "image/jpeg":
		ext = ".jpg"
	case "image/png":
		ext = ".png"
	case "video/mp4":
		ext = ".mp4"
	default:
		log.Fatalf("unsupported content type: %s", contentType)
	}
	filename := filepath.Base(uri) + ext

	// Get length of data. Note that resp.ContentLength will return -1 if
	// the Content-Length header is not present or if it's not a valid integer.
	contentLength := resp.ContentLength

	// Convert the response bytes into reader.
	content := bytes.NewReader(body)

	// Create the response IDO.
	payload := &NFTAssetIDO{
		Filename:      filename,
		FileExtension: ext,
		ContentType:   contentType,
		ContentLength: contentLength,
		Content:       content,
	}

	return payload, nil
}

func (uc *GetNFTAssetUseCase) fetchViaIPFS(uri string) (*NFTAssetIDO, error) {
	cid := strings.Replace(uri, "ipfs://", "", -1)
	url := strings.Replace(constants.IPFSPublicHTTPSGatewayBaseURL, "{CID}", cid, -1)
	uc.logger.Debug("fetching from uri via IPFS protocol.",
		slog.String("ipfs_cid", cid),
		slog.String("url", url))

	return uc.fetchViaHTTPS(url)
}
