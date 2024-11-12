package service

import (
	"context"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

// ProofOfAuthorityConsensusMechanismService represents the service which
// delivers comparatively fast transactions using identity as a stake.
//
// Would you like to know more?
// https://coinmarketcap.com/academy/glossary/proof-of-authority-poa
type ProofOfAuthorityConsensusMechanismService struct {
	config                                     *config.Configuration
	logger                                     *slog.Logger
	kmutex                                     kmutexutil.KMutexProvider
	getProofOfAuthorityPrivateKeyService       *GetProofOfAuthorityPrivateKeyService
	mempoolTransactionInsertionDetectorUseCase *usecase.MempoolTransactionInsertionDetectorUseCase
	mempoolTransactionListByChainIDUseCase     *usecase.MempoolTransactionListByChainIDUseCase
	mempoolTransactionDeleteByChainIDUseCase   *usecase.MempoolTransactionDeleteByChainIDUseCase
	getBlockchainStateUseCase                  *usecase.GetBlockchainStateUseCase
	upsertBlockchainStateUseCase               *usecase.UpsertBlockchainStateUseCase
	getGenesisBlockDataUseCase                 *usecase.GetGenesisBlockDataUseCase
}

func NewProofOfAuthorityConsensusMechanismService(
	config *config.Configuration,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	s1 *GetProofOfAuthorityPrivateKeyService,
	uc1 *usecase.MempoolTransactionInsertionDetectorUseCase,
	uc2 *usecase.MempoolTransactionListByChainIDUseCase,
	uc3 *usecase.MempoolTransactionDeleteByChainIDUseCase,
	uc4 *usecase.GetBlockchainStateUseCase,
	uc5 *usecase.UpsertBlockchainStateUseCase,
	uc6 *usecase.GetGenesisBlockDataUseCase,
) *ProofOfAuthorityConsensusMechanismService {
	return &ProofOfAuthorityConsensusMechanismService{config, logger, kmutex, s1, uc1, uc2, uc3, uc4, uc5, uc6}
}

func (s *ProofOfAuthorityConsensusMechanismService) Execute(ctx context.Context) error {
	mempoolTx, err := s.mempoolTransactionInsertionDetectorUseCase.Execute(ctx)
	if err != nil {
		return err
	}
	s.logger.Debug("New memory pool transaction detected!", slog.Any("record", mempoolTx))
	return nil
}
