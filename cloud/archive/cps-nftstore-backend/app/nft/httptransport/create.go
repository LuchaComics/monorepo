package httptransport

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	_ "time/tzdata"

	control "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/controller"
	sub_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/nft/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

func (h *Handler) unmarshalCreateRequest(ctx context.Context, r *http.Request) (*control.NFTCreateRequestIDO, error) {
	// Initialize our array which will tenant all the results from the remote server.
	var requestData control.NFTCreateRequestIDO

	defer r.Body.Close()

	var rawJSON bytes.Buffer
	teeReader := io.TeeReader(r.Body, &rawJSON) // TeeReader allows you to read the JSON and capture it

	// Read the JSON string and convert it into our golang stuct else we need
	// to send a `400 Bad Request` errror message back to the client,
	err := json.NewDecoder(teeReader).Decode(&requestData)
	if err != nil {
		h.Logger.Error("failed decoding",
			slog.Any("error", err))
		return nil, httperror.NewForSingleField(http.StatusBadRequest, "non_field_error", "payload structure is wrong")
	}

	return &requestData, nil
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	payload, err := h.unmarshalCreateRequest(ctx, r)
	if err != nil {
		h.Logger.Error("unmarshal create request",
			slog.Any("error", err))
		httperror.ResponseError(w, err)
		return
	}

	res, err := h.Controller.Create(ctx, payload)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	MarshalCreateResponse(res, w)
}

func MarshalCreateResponse(res *sub_s.NFT, w http.ResponseWriter) {
	if err := json.NewEncoder(w).Encode(&res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}
