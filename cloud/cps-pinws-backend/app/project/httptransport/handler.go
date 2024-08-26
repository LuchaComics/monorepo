package httptransport

import (
	"log/slog"

	project_c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/project/controller"
)

// Handler Creates http request handler
type Handler struct {
	Logger     *slog.Logger
	Controller project_c.ProjectController
}

// NewHandler Constructor
func NewHandler(loggerp *slog.Logger, c project_c.ProjectController) *Handler {
	return &Handler{
		Logger:     loggerp,
		Controller: c,
	}
}
