package handler

import (
	"bytes"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/httperror"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/service"
)

type IPFSGatewayHTTPHandler struct {
	logger                   *slog.Logger
	pinObjectGetByCIDService *service.PinObjectGetByCIDService
}

func NewIPFSGatewayHTTPHandler(
	logger *slog.Logger,
	s *service.PinObjectGetByCIDService,
) *IPFSGatewayHTTPHandler {
	return &IPFSGatewayHTTPHandler{logger, s}
}

func (h *IPFSGatewayHTTPHandler) Execute(w http.ResponseWriter, r *http.Request, cid string) {
	ctx := r.Context()

	// Extract url parameters.
	query := r.URL.Query()

	// Get the IPFS Gateway spec parameters.
	filenameQuery := query.Get("filename")
	downloadQuery := query.Get("download")

	pinobj, err := h.pinObjectGetByCIDService.Execute(ctx, cid)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if pinobj == nil {
		httperror.ResponseError(w, httperror.NewForNotFoundWithSingleField("cid", "does not exist"))
		return
	}
	h.logger.Debug("Fetched the local pinned-object",
		slog.String("content_type", pinobj.Meta["content_type"]),
		slog.String("filename", pinobj.Meta["filename"]),
		slog.String("cid", cid))

	// 2.2.1 filename (request query parameter) via
	// https://specs.ipfs.tech/http-gateways/path-gateway/#filename-request-query-parameter
	var filename string
	if filenameQuery != "" {
		filename = filenameQuery
	} else {
		filename = pinobj.Meta["filename"]
	}

	contentType := pinobj.Meta["content_type"]

	// Set Content-Disposition header
	var attch string

	// 2.2.2 download (request query parameter) via
	// https://specs.ipfs.tech/http-gateways/path-gateway/#download-request-query-parameter
	if downloadQuery == "true" {
		attch = fmt.Sprintf("attachment;filename*=\"%v\"", filename)

		// 3.2.5 Content-Disposition (response header)
		// https://specs.ipfs.tech/http-gateways/path-gateway/#content-disposition-response-header
		w.Header().Set("Content-Disposition", attch)
	} else {
		attch = fmt.Sprintf("inline;filename*=\"%v\"", filename)

		// 3.2.5 Content-Disposition (response header)
		// https://specs.ipfs.tech/http-gateways/path-gateway/#content-disposition-response-header
		w.Header().Set("Content-Disposition", attch)
	}

	// 3.2.4 Content-Type (response header)
	// https://specs.ipfs.tech/http-gateways/path-gateway/#content-type-response-header
	w.Header().Set("Content-Type", contentType)

	// 3.2.7 Content-Length (response header)
	// https://specs.ipfs.tech/http-gateways/path-gateway/#content-length-response-header
	w.Header().Set("Content-Length", strconv.Itoa(len(pinobj.Content)))

	// 3.2.3 Last-Modified (response header)
	// https://specs.ipfs.tech/http-gateways/path-gateway/#last-modified-response-header
	// Format the time in the correct format for the Last-Modified header
	w.Header().Set("Last-Modified", pinobj.ModifiedAt.UTC().Format(http.TimeFormat))

	// Add the X-Content-Type-Options header to prevent MIME type sniffing
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Convert []byte to io.Reader using bytes.NewReader
	reader := bytes.NewReader(pinobj.Content)

	// 3.3 Response Payload
	// https://specs.ipfs.tech/http-gateways/path-gateway/#response-payload
	// Stream the content directly to the HTTP response
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to write content", http.StatusInternalServerError)
		return
	}
}
