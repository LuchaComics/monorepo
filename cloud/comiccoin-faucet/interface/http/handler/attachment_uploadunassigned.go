package handler

import (
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	_ "time/tzdata"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/service"
)

type AttachmentCreateHTTPHandler struct {
	logger   *slog.Logger
	dbClient *mongo.Client
	service  *service.AttachmentCreateService
}

func NewAttachmentCreateHTTPHandler(
	logger *slog.Logger,
	dbClient *mongo.Client,
	service *service.AttachmentCreateService,
) *AttachmentCreateHTTPHandler {
	return &AttachmentCreateHTTPHandler{
		logger:   logger,
		dbClient: dbClient,
		service:  service,
	}
}

func (h *AttachmentCreateHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	////
	//// Unmarshal request
	////

	// Set the maximum upload size (100 MB in this example)
	r.Body = http.MaxBytesReader(w, r.Body, 100<<20) // 100 MB

	// Extract the filename from the "Content-Disposition" header, if provided
	contentDisposition := r.Header.Get("Content-Disposition")
	if contentDisposition == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Disposition` header in your request"))
		return
	}
	h.logger.Debug("AttachmentCreateHTTPHandler: Create", slog.Any("Content-Disposition", contentDisposition))
	filename := getFilenameFromContentDispositionText(contentDisposition)
	if filename == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing or corrupt `filename` from your requests `Content-Disposition` text"))
		return
	}

	// Extract the content-type from the request header
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		httperror.ResponseError(w, httperror.NewForBadRequestWithSingleField("error", "missing `Content-Type` header in your request"))
		return
	}

	// Read the binary data from the request body into a byte slice
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	// Initialize our array which will store all the results from the remote server.
	requestData := &service.AttachmentCreateRequestIDO{
		Filename:    filename,
		ContentType: contentType,
		Data:        data,
	}

	////
	//// Start the transaction.
	////

	session, err := h.dbClient.StartSession()
	if err != nil {
		h.logger.Error("start session error",
			slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		resp, err := h.service.Execute(sessCtx, requestData)
		if err != nil {
			h.logger.Error("service execution failure",
				slog.Any("error", err))
			return nil, err
		}
		return resp, err
	}

	// Start a transaction
	result, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		h.logger.Error("session failed error",
			slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}

	resp := result.(*service.AttachmentCreateResponseIDO)

	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&resp); err != nil {
		h.logger.Error("Encoding failed",
			slog.Any("error", err))
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func getFilenameFromContentDispositionText(contentDispositionText string) string {
	// Define the regular expression pattern to extract the filename
	pattern := `filename="?([^"]+)"?`
	re := regexp.MustCompile(pattern)

	// Find the first match
	matches := re.FindStringSubmatch(contentDispositionText)
	if len(matches) > 1 {
		// Clean the filename:
		// 1. Trim any leading or trailing whitespace
		// 2. Remove any quotes
		// 3. Replace any potentially problematic characters
		filename := strings.TrimSpace(matches[1])
		filename = strings.Trim(filename, `"`)           // Remove any remaining quotes
		filename = strings.ReplaceAll(filename, `\`, "") // Remove any backslashes

		return filename
	}

	log.Println("contentDispositionText:", contentDispositionText)
	return ""
}
