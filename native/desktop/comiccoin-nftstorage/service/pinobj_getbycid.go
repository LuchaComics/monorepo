package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin-nftstorage/usecase"
)

type PinObjectGetByCIDService struct {
	logger                   *slog.Logger
	pinObjectGetByCIDUseCase *usecase.PinObjectGetByCIDUseCase
	ipfsGetUseCase           *usecase.IPFSGetUseCase
}

func NewPinObjectGetByCIDService(
	logger *slog.Logger,
	uc1 *usecase.PinObjectGetByCIDUseCase,
	uc2 *usecase.IPFSGetUseCase,
) *PinObjectGetByCIDService {
	return &PinObjectGetByCIDService{logger, uc1, uc2}
}

func (s *PinObjectGetByCIDService) Execute(ctx context.Context, cid string) (*domain.PinObject, error) {
	content, err := s.ipfsGetUseCase.Execute(ctx, cid)
	if err != nil {
		s.logger.Error("Failed getting from ipfs",
			slog.Any("cid", cid),
			slog.Any("error", err))
		return nil, err
	}
	pinobj, err := s.pinObjectGetByCIDUseCase.Execute(cid)
	if err != nil {
		s.logger.Error("Failed getting pinobject locally",
			slog.Any("cid", cid),
			slog.Any("error", err))
		return nil, err
	}
	if pinobj == nil {
		err := fmt.Errorf("Does not exist for CID: %v", cid)
		s.logger.Error("Failed getting pinobject locally",
			slog.Any("cid", cid),
			slog.Any("error", err))
		return nil, err
	}

	pinobj.Content = content

	return pinobj, nil
}
