package httptransport

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	a_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/controller"
	a_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	// Set the maximum upload size (100 MB in this example)
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100 MB

	// Extract the filename from the "Content-Disposition" header, if provided
	contentDisposition := r.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Disposition` header in your request"))
	}
	filename := getFilenameFromContentDispositionText(contentDisposition)
	if filename == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing or corrupt `filename` from your requests `Content-Disposition` text"))
	}

	// Extract the content-type from the request header
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Type` header in your request"))
	}

	// Read the binary data from the request body into a byte slice
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Initialize our array which will store all the results from the remote server.
	requestData := &a_c.NFTAssetCreateRequestIDO{
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}

	nftasset, err := h.Controller.Create(ctx, requestData)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(nftasset, w)
}

func MarshalCreateResponse(res *a_s.NFTAsset, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getFilenameFromContentDispositionText(contentDispositionText string) string {
	// Define the regular expression pattern to extract the filename
	pattern := `filename=([^;]+)`
	re := regexp.MustCompile(pattern)

	// Find the first match
	matches := re.FindStringSubmatch(contentDispositionText)
	if len(matches) > 1 {
		// Trim any leading or trailing whitespace around the filename
		return strings.TrimSpace(matches[1])
	} else {
		log.Println("contentDispositionText:", contentDispositionText)
		return ""
	}
}
