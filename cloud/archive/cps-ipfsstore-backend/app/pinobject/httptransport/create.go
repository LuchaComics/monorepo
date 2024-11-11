package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"go.mongodb.org/mongo-driver/bson/primitive"

	a_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/controller"
	a_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

func UnmarshalCreateRequest(ctx context.Context, r *http.Request) (*a_c.PinObjectCreateRequestIDO, error) {
	defer r.Body.Close()

	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Limit the maximum memory used for parsing to 32MB
	if err != nil {
		log.Println("UnmarshalCreateRequest:ParseMultipartForm:err:", err)
		return nil, err
	}

	// Get the values of form fields
	name := r.FormValue("name")
	originsStr := r.FormValue("origins")
	metaStr := r.FormValue("meta")
	projectID := r.FormValue("project_id")

	// Get the uploaded file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("UnmarshalCmsImageCreateRequest:FormFile:err:", err)
		// return nil, err, http.StatusInternalServerError
	}

	pid, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		log.Println("UnmarshalCmsImageCreateRequest: primitive.ObjectIDFromHex:err:", err)
	}

	// Parse `originsStr` into a slice of strings
	origins := strings.Split(originsStr, ",")

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
	requestData := &a_c.PinObjectCreateRequestIDO{
		Name:      name,
		ProjectID: pid,
		Origins:   origins,
		Meta:      meta,
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

	data, err := UnmarshalCreateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	pinobject, err := h.Controller.Create(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(pinobject, w)
}

func MarshalCreateResponse(res *a_s.PinObject, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
