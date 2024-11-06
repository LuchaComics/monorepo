package handler

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstore/service"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/common/httperror"
)

type IPFSPinAddHTTPHandler struct {
	logger  *slog.Logger
	service *service.IPFSPinAddService
}

func NewIPFSPinAddHTTPHandler(
	logger *slog.Logger,
	service *service.IPFSPinAddService,
) *IPFSPinAddHTTPHandler {
	return &IPFSPinAddHTTPHandler{logger, service}
}

func (h *IPFSPinAddHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		h.logger.Error("Authorization header is missing")
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
		h.logger.Error("missing `Content-Disposition` header in your request")
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Disposition` header in your request"))
	}
	filename := getFilenameFromContentDispositionText(contentDisposition)
	if filename == "" {
		h.logger.Error("missing or corrupt `filename` from your requests `Content-Disposition` text")
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing or corrupt `filename` from your requests `Content-Disposition` text"))
	}

	// Extract the content-type from the request header
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		h.logger.Error("missing `Content-Type` header in your request")
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Type` header in your request"))
	}

	// Read the binary data from the request body into a byte slice
	data, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.Error("Failed to read request body")
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	req := &service.IPFSPinAddRequestIDO{
		ApiKey:      apiKey,
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}
	resp, err := h.service.Execute(context.Background(), req)
	if err != nil {
		h.logger.Error("Failed executing ipfs pin-add", slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		h.logger.Error("Failed encoding response", slog.Any("error", err))
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
