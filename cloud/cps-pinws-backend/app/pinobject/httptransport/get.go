package httptransport

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"go.mongodb.org/mongo-driver/bson/primitive"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/utils/httperror"
)

func (h *Handler) GetByRequestID(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()

	objectRequestID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	m, err := h.Controller.GetByRequestID(ctx, objectRequestID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalDetailResponse(m, w)
}

func MarshalDetailResponse(res *sub_s.PinObject, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) GetContentByRequestID(w http.ResponseWriter, r *http.Request, requestID string) {
	ctx := r.Context()
	objectRequestID, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	h.Logger.Debug("get by requestid", slog.String("requestid", requestID))

	m, err := h.Controller.GetWithContentByRequestID(ctx, objectRequestID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	h.Logger.Debug("retrieved pin", slog.String("requestid", requestID))

	filename := m.Meta["filename"]
	contentType := m.Meta["content_type"]

	if filename == "" {
		filename = "default_filename.txt"
	}
	if contentType == "" {
		contentType = "application/octet-stream"
	}

	h.Logger.Debug("content",
		slog.String("filename", filename),
		slog.String("contentType", contentType))

	// Note: https://stackoverflow.com/a/24116517
	attch := fmt.Sprintf("attachment;filename*=\"%v\"", filename)
	w.Header().Set("Content-Disposition", attch)
	w.Header().Set("Content-Type", contentType)
	w.Header().Set("Content-Length", strconv.Itoa(len(m.Content)))

	h.Logger.Debug("content",
		slog.String("attch", attch),
		slog.String("filename", filename),
		slog.String("contentType", contentType))

	// Convert []byte to io.Reader using bytes.NewReader
	reader := bytes.NewReader(m.Content)

	// Stream the content directly to the HTTP response without fully loading it into memory (for big files this is important) - simply copy the body reader to the response writer:
	_, err = io.Copy(w, reader)
	if err != nil {
		http.Error(w, "Failed to write content", http.StatusInternalServerError)
		return
	}
}
