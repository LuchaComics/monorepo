package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOperationArchiveRequest struct {
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}

func UnmarshalOperationArchiveRequest(ctx context.Context, r *http.Request) (*UserOperationArchiveRequest, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData UserOperationArchiveRequest

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println("UnmarshalOperationArchiveRequest | NewDecoder/Decode | err:", err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateOperationArchiveRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateOperationArchiveRequest(dirtyData *UserOperationArchiveRequest) error {
	e := make(map[string]string)

	if dirtyData.UserID.IsZero() {
		e["user_id"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) OperationArchive(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqData, err := UnmarshalOperationArchiveRequest(ctx, r)
	if err != nil {
		log.Println("OperationArchive | UnmarshalOperationArchiveRequest | err:", err)
		httperror.ResponseError(w, err)
		return
	}
	data, err := h.Controller.ArchiveByID(ctx, reqData.UserID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalOperationArchiveResponse(data, w)
}

func MarshalOperationArchiveResponse(res *sub_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
