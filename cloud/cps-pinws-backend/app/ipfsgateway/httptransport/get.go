package httptransport

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

// IPFS Gateway Spec via https://specs.ipfs.tech/http-gateways/path-gateway/

func (h *Handler) GetByContentID(w http.ResponseWriter, r *http.Request, cid string) {
	ctx := r.Context()

	// Extract url parameters.
	query := r.URL.Query()

	// Get the IPFS Gateway spec parameters.
	downloadQuery := query.Get("download")

	m, err := h.Controller.GetByContentID(ctx, cid)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	if m == nil {
		httperror.ResponseError(w, httperror.NewForNotFoundWithSingleField("cid", "does not exist"))
		return
	}

	filename := m.Meta["filename"]
	contentType := m.Meta["content_type"]

	// Set Content-Disposition header
	var attch string
	if downloadQuery == "true" {
		attch = fmt.Sprintf("attachment;filename*=\"%v\"", filename)
		w.Header().Set("Content-Disposition", attch)
	} else {
		attch = fmt.Sprintf("inline;filename*=\"%v\"", filename)
		w.Header().Set("Content-Disposition", attch)
	}

	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(m.Content)))

	// Add the X-Content-Type-Options header to prevent MIME type sniffing
	w.Header().Set("X-Content-Type-Options", "nosniff")

	// Convert []byte to io.Reader using bytes.NewReader
	reader := bytes.NewReader(m.Content)

	// Stream the content directly to the HTTP response
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to write content", http.StatusInternalServerError)
		return
	}
}
