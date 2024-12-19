package service

import (
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/kmutexutil"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type FaucetBalanceService struct {
	config            *config.Configuration
	logger            *slog.Logger
	kmutex            kmutexutil.KMutexProvider
	getAccountUseCase *usecase.GetAccountUseCase
}

func NewFaucetBalanceService(
	cfg *config.Configuration,
	logger *slog.Logger,
	kmutex kmutexutil.KMutexProvider,
	uc1 *usecase.GetAccountUseCase,
) *FaucetBalanceService {
	return &FaucetBalanceService{cfg, logger, kmutex, uc1}
}

type FaucetBalanceResponseIDO struct {
	Count uint64 `bson:"count" json:"count"`
}

func (s *FaucetBalanceService) Execute(sessCtx mongo.SessionContext) (*FaucetBalanceResponseIDO, error) {
	account, err := s.getAccountUseCase.Execute(sessCtx, s.config.App.WalletAddress)
	if err != nil {
		s.logger.Error("failed getting account",
			slog.Any("wallet_address", s.config.App.WalletAddress),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting account: %s", err)
	}
	if account == nil {
		return nil, fmt.Errorf("failed getting account: %s", "d.n.e.")
	}
	return &FaucetBalanceResponseIDO{Count: account.Balance}, nil
}
