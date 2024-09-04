package httptransport

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"

	a_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/pinobject/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) IpfsAdd(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		http.Error(w, "Authorization header is missing", http.StatusUnauthorized)
		return
	}

	// Parse the JWT token
	apiKey := strings.TrimPrefix(authHeader, "JWT ")

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

	// // Optionally, you can log the received information
	// fmt.Printf("Received file: %s\n", filename)
	// fmt.Printf("Content-Type: %s\n", contentType)
	// fmt.Printf("File size: %d bytes\n", len(data))
	// fmt.Printf("apiKey: %s\n", apiKey)

	req := &a_c.IpfsAddRequestIDO{
		ApiKey:      apiKey,
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}

	ctx := r.Context()
	cid, err := h.Controller.IpfsAdd(ctx, req)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Respond back to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(cid))
}

func MarshalIpfsAddResponse(res *a_c.IpfsAddResponseIDO, w http.ResponseWriter) {
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
