package handler

import (
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"strings"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/common/httperror"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/service"
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

	// Extract the content-type from the request header
	contentType := r.Header.Get("X-File-Content-Type")
	if contentType == "" {
		h.logger.Error("missing `X-File-Content-Type` header in your request")
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Type` header in your request"))
	}

	// Parse the multipart form data
	mr, err := r.MultipartReader()
	if err != nil {
		h.logger.Error("Failed to create multipart reader", slog.Any("error", err))
		http.Error(w, "Failed to read multipart data", http.StatusBadRequest)
		return
	}

	var filename string
	var data []byte

	// Iterate over each part in the multipart form
	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break // End of parts
		}
		if err != nil {
			h.logger.Error("Error reading next part", slog.Any("error", err))
			http.Error(w, "Error reading multipart data", http.StatusInternalServerError)
			return
		}

		// Check if the part has a filename (i.e., it's a file upload)
		if part.FileName() != "" {
			filename = part.FileName()
			data, err = io.ReadAll(part)
			if err != nil {
				h.logger.Error("Failed to read file part", slog.Any("error", err))
				http.Error(w, "Failed to read file data", http.StatusInternalServerError)
				return
			}
			break // Process only the first file part
		}
	}

	if filename == "" {
		h.logger.Error("No file found in the multipart data")
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "No file found in the multipart data"))
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
		h.logger.Error("Failed executing ipfs pin-add",
			// slog.Any("apiKey", apiKey), // For debugging purposes only.
			slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}

	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		h.logger.Error("Failed encoding response", slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
