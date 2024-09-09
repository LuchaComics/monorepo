package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"

	a_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/controller"
	a_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftasset/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) unmarshalCreateRequest(ctx context.Context, r *http.Request) (*a_c.NFTAssetCreateRequestIDO, error) {
	defer r.Body.Close()

	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Limit the maximum memory used for parsing to 32MB
	if err != nil {
		h.Logger.Error("failed parising multipart form", slog.Any("error", err))
		return nil, err
	}

	// Get the values of form fields
	name := r.FormValue("name")
	metaStr := r.FormValue("meta")

	// Get the uploaded file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		h.Logger.Error("failed getting form file", slog.Any("error", err))
		return nil, err
	}

	/// Parse `metaStr` into a map of string keys and values
	var meta map[string]string
	if err := json.Unmarshal([]byte(metaStr), &meta); err != nil {
		log.Printf("Failed to parse meta: %v", err)
		meta = make(map[string]string) // Initialize an empty map if parsing fails
	}

	// Add default values.
	meta["filename"] = ""
	meta["content_type"] = ""

	// Initialize our array which will store all the results from the remote server.
	requestData := &a_c.NFTAssetCreateRequestIDO{
		Name: name,
		Meta: meta,
	}

	if header != nil {
		// Extract filename and filetype from the file header
		requestData.File = file

		// Handle meta. We will attach meta along with some custom fields.
		meta["filename"] = header.Filename
		meta["content_type"] = header.Header.Get("Content-Type")
		requestData.Meta = meta
	}
	return requestData, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := h.unmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	nftasset, err := h.Controller.Create(ctx, data)
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
