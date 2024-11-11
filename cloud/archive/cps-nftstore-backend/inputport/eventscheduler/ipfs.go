package eventscheduler

import (
	"context"
	"log/slog"
)

func (port *crontabInputPort) runReprovideIPNS() {
	// Notes:
	// How do I make my IPNS records live longer?
	// https://discuss.ipfs.tech/t/how-do-i-make-my-ipns-records-live-longer/14768

	port.Logger.Debug("starting to reprovide ipns")
	if err := port.NFTCollection.ReprovidehCollectionsInIPNS(context.Background()); err != nil {
		port.Logger.Error("failed to reprovide collection to ipns",
			slog.Any("error", err))
	}
	port.Logger.Debug("finished to reprovided ipns")
}
