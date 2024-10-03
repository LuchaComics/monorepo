package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/usecase"
)

type BlockDataDTOServerService struct {
	config                             *config.Config
	logger                             *slog.Logger
	listAllBlockDataUseCase            *usecase.ListAllBlockDataUseCase
	uploadToNetworkBlockDataDTOUseCase *usecase.UploadToNetworkBlockDataDTOUseCase
}

func NewBlockDataDTOServerService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.ListAllBlockDataUseCase,
	uc2 *usecase.UploadToNetworkBlockDataDTOUseCase,
) *BlockDataDTOServerService {
	return &BlockDataDTOServerService{cfg, logger, uc1, uc2}
}

func (s *BlockDataDTOServerService) Execute(ctx context.Context) error {
	blockDataList, err := s.listAllBlockDataUseCase.Execute()
	if err != nil {
		return err
	}
	if blockDataList == nil {
		err := fmt.Errorf("block data results: %v", "does not exist")
		return err
	}
	for _, blockData := range blockDataList {
		blockDataDTO := &domain.BlockDataDTO{
			Hash:   blockData.Hash,
			Header: blockData.Header,
			Trans:  blockData.Trans,
		}
		if err := s.uploadToNetworkBlockDataDTOUseCase.Execute(ctx, blockDataDTO); err != nil {
			return err
		}
	}

	return nil
}
