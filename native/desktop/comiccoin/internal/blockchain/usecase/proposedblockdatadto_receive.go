package usecase

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/blockchain/domain"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type ReceiveProposedBlockDataDTOUseCase struct {
	config *config.Config
	logger *slog.Logger
	repo   domain.ProposedBlockDataDTORepository
}

func NewReceiveProposedBlockDataDTOUseCase(config *config.Config, logger *slog.Logger, repo domain.ProposedBlockDataDTORepository) *ReceiveProposedBlockDataDTOUseCase {
	return &ReceiveProposedBlockDataDTOUseCase{config, logger, repo}
}

func (uc *ReceiveProposedBlockDataDTOUseCase) Execute(ctx context.Context) (*domain.ProposedBlockData, error) {
	//
	// STEP 1:
	// Wait to receive from the P2P Network. It just takes one node to publish
	// the data and then we will receive it here.
	//

	dto, err := uc.repo.ReceiveFromP2PNetwork(ctx)
	if err != nil {
		uc.logger.Error("failed receiving signed transaction dto from network",
			slog.Any("error", err))
		return nil, err
	}
	if dto == nil {
		// Developer Note:
		// If we haven't received anything, that means we haven't connected to
		// the distributed / P2P network, so all we can do is return nil.
		return nil, nil
	}

	//
	// STEP 2:
	// Convert back to our signed trnsaction data-type and then perform simple
	// validation before returning it for this function.
	//

	ido := &domain.ProposedBlockData{
		Hash:   dto.Hash,
		Header: dto.Header,
		Trans:  dto.Trans,
	}

	e := make(map[string]string)
	if ido.Hash == "" {
		e["hash"] = "missing value"
	}
	if ido.Header == nil {
		e["header"] = "missing value"
	} else {

	}
	if ido.Trans == nil {
		e["trans"] = "missing value"
	} else {
		if len(ido.Trans) == 0 {
			e["trans"] = "empty transactions"
		}
		// TODO: Validation for the following/
		// for _, tx := ido.Trans {
		// TODO:
		// SignedTransaction
		// TimeStamp uint64 `json:"timestamp"` // Ethereum: The time the transaction was received.
		// GasPrice  uint64 `json:"gas_price"` // Ethereum: The price of one unit of gas to be paid for fees.
		// GasUnits  uint64 `json:"gas_units"` // Ethereum: The number of units of gas used for this transaction.
		// }
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for received",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	uc.logger.Debug("Received proposed block data dto from network",
		slog.Any("hash", ido.Hash))

	return ido, nil
}
