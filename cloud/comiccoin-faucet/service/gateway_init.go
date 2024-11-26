package service

import (
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayInitService struct {
	config                  *config.Configuration
	logger                  *slog.Logger
	passwordProvider        password.Provider
	tenantGetByNameUseCase  *usecase.TenantGetByNameUseCase
	tenantCreate            *usecase.TenantCreateUseCase
	walletEncryptKeyUseCase *usecase.WalletEncryptKeyUseCase
	walletDecryptKeyUseCase *usecase.WalletDecryptKeyUseCase
	createWalletUseCase     *usecase.CreateWalletUseCase
	userGet                 *usecase.UserGetByEmailUseCase
	userCreate              *usecase.UserCreateUseCase
	createAccountUseCase    *usecase.CreateAccountUseCase
	getAccountUseCase       *usecase.GetAccountUseCase
}

func NewGatewayInitService(
	config *config.Configuration,
	logger *slog.Logger,
	pp password.Provider,
	uc1 *usecase.TenantGetByNameUseCase,
	uc2 *usecase.TenantCreateUseCase,
	uc3 *usecase.WalletEncryptKeyUseCase,
	uc4 *usecase.WalletDecryptKeyUseCase,
	uc5 *usecase.CreateWalletUseCase,
	uc6 *usecase.UserGetByEmailUseCase,
	uc7 *usecase.UserCreateUseCase,
	uc8 *usecase.CreateAccountUseCase,
	uc9 *usecase.GetAccountUseCase,
) *GatewayInitService {
	return &GatewayInitService{config, logger, pp, uc1, uc2, uc3, uc4, uc5, uc6, uc7, uc8, uc9}
}

func (s *GatewayInitService) Execute(
	sessCtx mongo.SessionContext,
	tenantName string,
	chainID uint16,
	email string,
	pass *sstring.SecureString,
) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if tenantName == "" {
		e["tenant_name"] = "missing value"
	} else {
		tenant, err := s.tenantGetByNameUseCase.Execute(sessCtx, tenantName)
		if err != nil {
			s.logger.Debug("Failed to get tenant by name")
			return err
		}
		if tenant != nil {
			err := fmt.Errorf("Tenant already exists with name: %v", tenantName)
			s.logger.Error("Failed because tenant exists", slog.Any("error", err))
			return err
		}
	}
	if chainID == 0 {
		e["chain_id"] = "missing value"
	}
	if email == "" {
		e["email"] = "missing value"
	} else {

	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Create our tenant
	//

	tenant := &domain.Tenant{
		ID:         primitive.NewObjectID(),
		Name:       tenantName,
		ChainID:    chainID,
		Status:     domain.TenantActiveStatus,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
	}

	err := s.tenantCreate.Execute(sessCtx, tenant)
	if err != nil {
		s.logger.Error("Failed creating tenant",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 3: Create our administrator
	//

	user := &domain.User{
		ID:          primitive.NewObjectID(),
		TenantID:    tenant.ID,
		TenantName:  tenantName,
		FirstName:   "System",
		LastName:    "Administrator",
		Name:        "System Administrator",
		LexicalName: "Administrator, System",
		Email:       email,
	}

	if err := s.userCreate.Execute(sessCtx, user); err != nil {
		s.logger.Error("Failed creating user",
			slog.Any("error", err))
		return err
	}

	// return uc.repo.Create(sessCtx, user)
	return nil
}
