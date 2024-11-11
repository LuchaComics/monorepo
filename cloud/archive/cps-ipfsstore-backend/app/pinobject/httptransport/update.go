package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	a_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/controller"
	sub_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/controller"
	sub_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/pinobject/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func UnmarshalUpdateRequest(ctx context.Context, r *http.Request) (*sub_c.PinObjectUpdateRequestIDO, error) {
	defer r.Body.Close()

	// Parse the multipart form data
	err := r.ParseMultipartForm(32 << 20) // Limit the maximum memory used for parsing to 32MB
	if err != nil {
		log.Println("UnmarshalUpdateRequest:ParseMultipartForm:err:", err)
		return nil, err
	}

	// Get the values of form fields
	requestID := r.FormValue("requestid")
	name := r.FormValue("name")
	projectID := r.FormValue("project_id")

	// Get the uploaded file from the request
	file, header, err := r.FormFile("file")
	if err != nil {
		log.Println("UnmarshalUpdateRequest:FormFile:err:", err)
		// return nil, err, http.StatusInternalServerError
	}

	pid, err := primitive.ObjectIDFromHex(projectID)
	if err != nil {
		log.Println("UnmarshalUpdateRequest: primitive.ObjectIDFromHex:err:", err)
	}

	rid, err := primitive.ObjectIDFromHex(requestID)
	if err != nil {
		log.Println("UnmarshalUpdateRequest: primitive.ObjectIDFromHex:err:", err)
	}

	// Initialize our array which will store all the results from the remote server.
	requestData := &a_c.PinObjectUpdateRequestIDO{
		RequestID: rid,
		Name:      name,
		ProjectID: pid,
	}

	if header != nil {
		// Extract filename and filetype from the file header
		requestData.FileName = header.Filename
		requestData.FileType = header.Header.Get("Content-Type")
		requestData.File = file
	}
	return requestData, nil
}

func (h *Handler) UpdateByRequestID(w http.ResponseWriter, r *http.Request, id string) {
	ctx := r.Context()

	data, err := UnmarshalUpdateRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}
	pinobject, err := h.Controller.UpdateByRequestID(ctx, data)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalUpdateResponse(pinobject, w)
}

func MarshalUpdateResponse(res *sub_s.PinObject, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
