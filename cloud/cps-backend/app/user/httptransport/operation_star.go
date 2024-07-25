package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	sub_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UserOperationStarRequest struct {
	UserID primitive.ObjectID `bson:"user_id" json:"user_id"`
}

func UnmarshalOperationStarRequest(ctx context.Context, r *http.Request) (*UserOperationStarRequest, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData UserOperationStarRequest

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println("UnmarshalOperationStarRequest | NewDecoder/Decode | err:", err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	// Perform our validation and return validation error on any issues detected.
	if err := ValidateOperationStarRequest(&requestData); err != nil {
		return nil, err
	}
	return &requestData, nil
}

func ValidateOperationStarRequest(dirtyData *UserOperationStarRequest) error {
	e := make(map[string]string)

	if dirtyData.UserID.IsZero() {
		e["user_id"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (h *Handler) OperationStar(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqData, err := UnmarshalOperationStarRequest(ctx, r)
	if err != nil {
		log.Println("OperationStar | UnmarshalOperationStarRequest | err:", err)
		httperror.ResponseError(w, err)
		return
	}
	data, err := h.Controller.Star(ctx, reqData.UserID)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalOperationStarResponse(data, w)
}

func MarshalOperationStarResponse(res *sub_s.User, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
