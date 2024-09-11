package httptransport

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

// IPFS Gateway Spec via https://specs.ipfs.tech/http-gateways/path-gateway/

// GetByCID functions provides HTTP interface for requesting content-addressed data at specified content path from IPFS network.
func (h *Handler) GetByCID(w http.ResponseWriter, r *http.Request, cid string) {
	ctx := r.Context()

	// Extract url parameters.
	query := r.URL.Query()

	// Get the IPFS Gateway spec parameters.
	filenameQuery := query.Get("filename")
	downloadQuery := query.Get("download")

	res, err := h.Controller.GetByCID(ctx, cid)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if res == nil {
		httperror.ResponseError(w, httperror.NewForNotFoundWithSingleField("cid", "does not exist"))
		return
	}
	// 2.2.1 filename (request query parameter) via
	// https://specs.ipfs.tech/http-gateways/path-gateway/#filename-request-query-parameter
	var filename string
	if filenameQuery != "" {
		filename = filenameQuery
	} else {
		filename = res.Filename
	}

	contentType := res.ContentType

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
	w.Header().Set("Content-Length", strconv.Itoa(len(res.Content)))

	// 3.2.3 Last-Modified (response header)
	// https://specs.ipfs.tech/http-gateways/path-gateway/#last-modified-response-header
	// Format the time in the correct format for the Last-Modified header
	w.Header().Set("Last-Modified", res.ModifiedAt.UTC().Format(http.TimeFormat))

	// Add the X-Content-Type-Options header to prevent MIME type sniffing
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Convert []byte to io.Reader using bytes.NewReader
	reader := bytes.NewReader(res.Content)

	// 3.3 Response Payload
	// https://specs.ipfs.tech/http-gateways/path-gateway/#response-payload
	// Stream the content directly to the HTTP response
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to write content", http.StatusInternalServerError)
		return
	}
}
