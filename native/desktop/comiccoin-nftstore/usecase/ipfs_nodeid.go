package usecase

import (
	"log/slog"

	pkg_domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
	"github.com/libp2p/go-libp2p/core/peer"
)

type IPFSGetNodeIDUseCase struct {
	logger   *slog.Logger
	ipfsRepo pkg_domain.IPFSRepository
}

func NewIPFSGetNodeIDUseCase(logger *slog.Logger, r1 pkg_domain.IPFSRepository) *IPFSGetNodeIDUseCase {
	return &IPFSGetNodeIDUseCase{logger, r1}
}

func (uc *IPFSGetNodeIDUseCase) Execute() (peer.ID, error) {
	return uc.ipfsRepo.ID()
}
