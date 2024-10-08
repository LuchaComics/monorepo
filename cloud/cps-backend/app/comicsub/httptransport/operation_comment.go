package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ComicSubmissionOperationCreateCommentRequest struct {
	SubmissionID primitive.ObjectID `bson:"submission_id" json:"submission_id"`
	Content      string             `bson:"content" json:"content"`
}

func UnmarshalOperationCreateCommentRequest(ctx context.Context, r *http.Request) (*ComicSubmissionOperationCreateCommentRequest, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData ComicSubmissionOperationCreateCommentRequest

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println("UnmarshalOperationCreateCommentRequest | NewDecoder/Decode | err:", err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateOperationCreateCommentRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateOperationCreateCommentRequest(dirtyData *ComicSubmissionOperationCreateCommentRequest) error {
	e := make(map[string]string)

	if dirtyData.SubmissionID.Hex() == "" {
		e["submission_id"] = "missing value"
	}

	if dirtyData.Content == "" {
		e["content"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) OperationCreateComment(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqData, err := UnmarshalOperationCreateCommentRequest(ctx, r)
	if err != nil {
		log.Println("OperationCreateComment | UnmarshalOperationCreateCommentRequest | err:", err)
		httperror.ResponseError(w, err)
		return
	}
	data, err := h.Controller.CreateComment(ctx, reqData.SubmissionID, reqData.Content)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalOperationCreateCommentResponse(data, w)
}

func MarshalOperationCreateCommentResponse(res *sub_s.ComicSubmission, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
