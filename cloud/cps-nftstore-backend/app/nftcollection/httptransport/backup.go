package httptransport

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	_ "time/tzdata"

	sub_c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nftcollection/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func UnmarshalBackupRequest(ctx context.Context, r *http.Request) (*sub_c.NFTCollectionBackupOperationRequestIDO, error) {
	// Initialize our array which will tenant all the results from the remote server.
	var requestData sub_c.NFTCollectionBackupOperationRequestIDO

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println(err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *Handler) OperationBackupInJSON(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := UnmarshalBackupRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.OperationBackup(ctx, payload)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Convert the result to JSON (assuming it's JSON data to be downloaded)
	fileContent, err := json.Marshal(&res)
	if err != nil {
		http.Error(w, "Error marshaling data", http.StatusInternalServerError)
		return
	}

	// Set headers for the downloadable file
	w.Header().Set("Content-Disposition", "attachment; filename=backup.json")
	w.Header().Set("Content-Type", "application/json") // You can adjust this if you're downloading another file type.
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(fileContent)))

	// Write the file content to the response
	_, err = w.Write(fileContent)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}
}

func (h *Handler) OperationBackupInXML(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := UnmarshalBackupRequest(ctx, r)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.OperationBackup(ctx, payload)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Convert the result to XML
	fileContent, err := xml.MarshalIndent(&res, "", "  ") // Pretty-printed XML output
	if err != nil {
		http.Error(w, "Error marshaling data to XML", http.StatusInternalServerError)
		return
	}

	// Set headers for the downloadable file
	w.Header().Set("Content-Disposition", "attachment; filename=backup.xml")
	w.Header().Set("Content-Type", "application/xml") // Set the correct Content-Type for XML
	w.Header().Set("Content-Length", fmt.Sprintf("%v", len(fileContent)))

	// Write the file content to the response
	_, err = w.Write(fileContent)
	if err != nil {
		http.Error(w, "Error writing file", http.StatusInternalServerError)
		return
	}
}
