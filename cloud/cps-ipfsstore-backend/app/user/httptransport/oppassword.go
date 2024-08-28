package httptransport

import (
	"context"
	"encoding/json"
	"log"
	"net/http"

	user_c "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/controller"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

func UnmarshalOperationChangePasswordRequest(ctx context.Context, r *http.Request) (*user_c.UserOperationChangePasswordRequest, error) {
	// Initialize our array which will store all the results from the remote server.
	var requestData user_c.UserOperationChangePasswordRequest

	defer r.Body.Close()

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(r.Body).Decode(&requestData) // [1]
	if err != nil {
		log.Println("UnmarshalOperationChangePasswordRequest | NewDecoder/Decode | err:", err)
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *Handler) OperationChangePassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	reqData, err := UnmarshalOperationChangePasswordRequest(ctx, r)
	if err != nil {
		log.Println("OperationChangePassword | UnmarshalOperationChangePasswordRequest | err:", err)
		httperror.ResponseError(w, err)
		return
	}

	if err := h.Controller.ChangePassword(ctx, reqData); err != nil {
		httperror.ResponseError(w, err)
		return
	}

}
