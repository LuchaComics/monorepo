package usecase

import (
	"log/slog"

	"github.com/libp2p/go-libp2p/core/peer"

	domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
)

type IPFSGetNodeIDUseCase struct {
	logger   *slog.Logger
	ipfsRepo domain.IPFSRepository
}

func NewIPFSGetNodeIDUseCase(logger *slog.Logger, r1 domain.IPFSRepository) *IPFSGetNodeIDUseCase {
	return &IPFSGetNodeIDUseCase{logger, r1}
}

func (uc *IPFSGetNodeIDUseCase) Execute() (peer.ID, error) {
	return uc.ipfsRepo.ID()
}
