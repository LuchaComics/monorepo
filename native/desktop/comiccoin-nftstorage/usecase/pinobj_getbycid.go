package usecase

import (
	"log/slog"

	domain "github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
)

type PinObjectGetByCIDUseCase struct {
	logger *slog.Logger
	repo   domain.PinObjectRepository
}

func NewPinObjectGetByCIDUseCase(logger *slog.Logger, r1 domain.PinObjectRepository) *PinObjectGetByCIDUseCase {
	return &PinObjectGetByCIDUseCase{logger, r1}
}

func (uc *PinObjectGetByCIDUseCase) Execute(cid string) (*domain.PinObject, error) {
	return uc.repo.GetByCID(cid)
}
