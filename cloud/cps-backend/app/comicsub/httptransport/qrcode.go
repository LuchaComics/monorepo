package httptransport

import (
	"bytes"
	"net/http"
	"time"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (h *Handler) GetQRCodePNGImageOfRegisteryURLByCPSRN(w http.ResponseWriter, r *http.Request, cpsrn string) {
	ctx := r.Context()

	pngImage, err := h.Controller.GetQRCodePNGImageOfRegisteryURLByCPSRN(ctx, cpsrn)
	if err != nil {
		httperror.ResponseError(w, err)
		return
	}

	// Set the Content-Type header
	w.Header().Set("Content-Type", "image/png")

	// Serve the content
	http.ServeContent(w, r, "qrcode.png", time.Now(), bytes.NewReader(pngImage))
}
