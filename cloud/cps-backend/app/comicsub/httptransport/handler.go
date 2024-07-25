package httptransport

import (
	"log/slog"

	comicsub_c "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller comicsub_c.ComicSubmissionController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c comicsub_c.ComicSubmissionController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
