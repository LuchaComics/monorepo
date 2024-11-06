package handler

import (
	"fmt"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type NFTAssetPinAddHTTPHandler struct {
	logger  *slog.Logger
	service *service.NFTAssetPinAddService
}

func NewNFTAssetPinAddHTTPHandler(
	logger *slog.Logger,
	service *service.NFTAssetPinAddService,
) *NFTAssetPinAddHTTPHandler {
	return &NFTAssetPinAddHTTPHandler{logger, service}
}

func (h *NFTAssetPinAddHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
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

	req := &IpfsAddRequestIDO{
		ApiKey:      apiKey,
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}
	fmt.Println(req) //TODO: THIS CODE WORKS, CONTINUE DEV. HERE!
}

type IpfsAddRequestIDO struct {
	ApiKey      string `bson:"api_key" json:"api_key"`
	Filename    string `bson:"filename" json:"filename"`
	ContentType string `bson:"content_type" json:"content_type"`
	Data        []byte `bson:"data" json:"data"`
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
