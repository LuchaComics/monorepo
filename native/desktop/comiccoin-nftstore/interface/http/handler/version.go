package handler

import (
	"log/slog"
	"net/http"
)

type GetVersionHTTPHandler struct {
	logger *slog.Logger
}

func NewGetVersionHTTPHandler(
	logger *slog.Logger,
) *GetVersionHTTPHandler {
	return &GetVersionHTTPHandler{logger}
}

func (h *GetVersionHTTPHandler) Execute(w http.ResponseWriter, r *http.Request) {

}
