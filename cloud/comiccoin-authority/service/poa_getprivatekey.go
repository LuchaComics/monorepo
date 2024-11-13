package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/ethereum/go-ethereum/accounts/keystore"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-authority/usecase"
)

type GetProofOfAuthorityPrivateKeyService struct {
	config                  *config.Configuration
	logger                  *slog.Logger
	getWalletUseCase        *usecase.GetWalletUseCase
	walletDecryptKeyUseCase *usecase.WalletDecryptKeyUseCase
}

func NewGetProofOfAuthorityPrivateKeyService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc1 *usecase.GetWalletUseCase,
	uc2 *usecase.WalletDecryptKeyUseCase,
) *GetProofOfAuthorityPrivateKeyService {
	return &GetProofOfAuthorityPrivateKeyService{cfg, logger, uc1, uc2}
}

func (s *GetProofOfAuthorityPrivateKeyService) Execute(ctx context.Context) (*keystore.Key, error) {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if s.config.Blockchain.ProofOfAuthorityAccountAddress == nil {
		e["COMICCOIN_AUTHORITY_BLOCKCHAIN_PROOF_OF_AUTHORITY_WALLET_PASSWORD"] = "missing environment value"
	}
	if s.config.Blockchain.ProofOfAuthorityWalletPassword == "" {
		e["COMICCOIN_AUTHORITY_BLOCKCHAIN_PROOF_OF_AUTHORITY_WALLET_PASSWORD"] = "missing environment value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating get key parameters",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Return the account.
	//

	wallet, err := s.getWalletUseCase.Execute(ctx, s.config.Blockchain.ProofOfAuthorityAccountAddress)
	if err != nil {
		s.logger.Error("failed getting wallet from database",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting wallet from database: %s", err)
	}
	if wallet == nil {
		return nil, fmt.Errorf("failed getting wallet from database: %s", "d.n.e.")
	}

	key, err := s.walletDecryptKeyUseCase.Execute(ctx, wallet.KeystoreBytes, s.config.Blockchain.ProofOfAuthorityWalletPassword)
	if err != nil {
		s.logger.Error("failed getting key",
			slog.Any("error", err))
		return nil, fmt.Errorf("failed getting key: %s", err)
	}
	if key == nil {
		return nil, fmt.Errorf("failed getting key: %s", "d.n.e.")
	}
	return key, nil
}
