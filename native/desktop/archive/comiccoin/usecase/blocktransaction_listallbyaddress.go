package usecase

import (
	"log/slog"

	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/domain"
)

type ListAllBlockTransactionByAddressUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.BlockDataRepository
}

func NewListAllBlockTransactionByAddressUseCase(config *config.Config, logger *slog.Logger, repo domain.BlockDataRepository) *ListAllBlockTransactionByAddressUseCase {
	return &ListAllBlockTransactionByAddressUseCase{config, logger, repo}
}

func (uc *ListAllBlockTransactionByAddressUseCase) Execute(addr *common.Address) ([]*domain.BlockTransaction, error) {
	data, err := uc.repo.ListAllBlockTransactionsByAddress(addr)
	if err != nil {
		uc.logger.Error("failed listing all block data", slog.Any("error", err))
		return nil, err
	}
	return data, nil
}
