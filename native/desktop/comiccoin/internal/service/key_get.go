package service

import (
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/common"

	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/config"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/internal/usecase"
	"github.com/LuchaComics/monorepo/native/desktop/comiccoin/pkg/httperror"
)

type GetKeyService struct {
	config                  *config.Config
	logger                  *slog.Logger
	getWalletUseCase        *usecase.GetWalletUseCase
	walletDecryptKeyUseCase *usecase.WalletDecryptKeyUseCase
}

func NewGetKeyService(
	cfg *config.Config,
	logger *slog.Logger,
	uc1 *usecase.GetWalletUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
) *GetKeyService {
	return &GetKeyService{cfg, logger, uc1, uc2}
}

func (s *GetKeyService) Execute(walletAddress *common.Address, walletPassword string) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if walletAddress == nil {
		e["wallet_address"] = "missing value"
	}
	if walletPassword == "" {
		e["wallet_password"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating get key parameters",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	wallet, err := s.getWalletUseCase.Execute(walletAddress)
	if err != nil {
		s.logger.Error("failed getting from database",
			slog.Any("wallet_address", walletAddress),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting from database: %s", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("failed getting from database: %s", "d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(wallet.Filepath, walletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("wallet_address", walletAddress),
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return nil, fmt.Errorf("failed getting key: %s", "d.n.e.")
	}
	return key, nil
}
