package handler

import (
	"log/slog"
	"time"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/account/service"
)

type AccountServer struct {
	cfg           *config.Config
	logger        *slog.Logger
	getKeyService *service.GetKeyService
}

func NewAccountServer(
	cfg *config.Config,
	logger *slog.Logger,
	getKeyService *service.GetKeyService,
) *AccountServer {
	accountServer := new(AccountServer)

	// Attach our dependencies
	accountServer.cfg = cfg
	accountServer.logger = logger
	accountServer.getKeyService = getKeyService

	return accountServer
}

type Args struct{}

func (t *AccountServer) GiveServerTime(args *Args, reply *int64) error {
	// Fill reply pointer to send the data back
	*reply = time.Now().Unix()
	return nil
}
