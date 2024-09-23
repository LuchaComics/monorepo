package http

import (
	"log/slog"

	keypair_ds "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/keypair/datastore"
	ledger_c "github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/app/ledger/controller"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/inputport"
)

type httpInputPort struct {
	cfg              *config.Config
	logger           *slog.Logger
	keypairStorer    keypair_ds.KeypairStorer
	ledgerController ledger_c.LedgerController
}

func NewInputPort(
	cfg *config.Config,
	logger *slog.Logger,
	kp keypair_ds.KeypairStorer,
	bc ledger_c.LedgerController,
) inputport.InputPortServer {
	// ctx := context.Background()

	// Begin our function by initializing the defaults for our peer-to-peer (p2p)
	// node and then the rest of this function pertains to setting up a p2p
	// network to utilize in our app.
	node := &httpInputPort{
		cfg:              cfg,
		logger:           logger,
		keypairStorer:    kp,
		ledgerController: bc,
	}

	return node
}

func (port *httpInputPort) Run() {
	// ctx := context.Background()
	port.logger.Info("Running HTTP JSON API")
}

func (port *httpInputPort) Shutdown() {
	port.logger.Info("Gracefully shutting down HTTP JSON API")
}
