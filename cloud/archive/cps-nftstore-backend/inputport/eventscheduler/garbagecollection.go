package eventscheduler

import (
	"context"
	"log/slog"
)

func (port *crontabInputPort) runGarbageCollection() {
	port.Logger.Debug("starting garbage collection...")
	if err := port.NFTAsset.DeleteByExecutingGarbageCollection(context.Background()); err != nil {
		port.Logger.Error("failed garbage collection",
			slog.Any("error", err))
	}
	port.Logger.Debug("finished garbage collection")
}
