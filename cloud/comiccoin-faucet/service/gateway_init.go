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
	config                 *config.Configuration
	logger                 *slog.Logger
	passwordProvider       password.Provider
	createAccountService   *CreateAccountService
	tenantGetByNameUseCase *usecase.TenantGetByNameUseCase
	tenantCreate           *usecase.TenantCreateUseCase
	userGet                *usecase.UserGetByEmailUseCase
	userCreate             *usecase.UserCreateUseCase
}

func NewGatewayInitService(
	config *config.Configuration,
	logger *slog.Logger,
	pp password.Provider,
	s1 *CreateAccountService,
	uc1 *usecase.TenantGetByNameUseCase,
	uc2 *usecase.TenantCreateUseCase,
	uc3 *usecase.UserGetByEmailUseCase,
	uc4 *usecase.UserCreateUseCase,
) *GatewayInitService {
	return &GatewayInitService{config, logger, pp, s1, uc1, uc2, uc3, uc4}
}

func (s *GatewayInitService) Execute(
	sessCtx mongo.SessionContext,
	tenantName string,
	chainID uint16,
	email string,
	walletPassword *sstring.SecureString,
	walletPasswordRepeated *sstring.SecureString,
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
		//TODO
	}
	if walletPasswordRepeated == nil {
		e["wallet_password_repeated"] = "missing value"
	}
	if walletPassword.String() != walletPasswordRepeated.String() {
		e["wallet_password"] = "do not match"
		e["wallet_password_repeated"] = "do not match"
	}
	if len(e) != 0 {
		s.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Setup our unique identifiers
	//

	tenantID := primitive.NewObjectID()
	userID := primitive.NewObjectID()

	//
	// STEP 3:
	// Create the encryted physical wallet on file and our account.
	//

	account, err := s.createAccountService.Execute(sessCtx, walletPassword, walletPasswordRepeated, tenantID.Hex())
	if err != nil {
		s.logger.Error("Failed creating account",
			slog.Any("error", err))
		return err
	}
	if account == nil {
		err := fmt.Errorf("Account was not returned for tenant: %v", tenantID.Hex())
		s.logger.Error("Failed creating account",
			slog.Any("error", err))
		return err
	}

	//
	// STEP 4: Create our tenant.
	//

	tenant := &domain.Tenant{
		ID:         tenantID,
		Name:       tenantName,
		ChainID:    chainID,
		Status:     domain.TenantActiveStatus,
		CreatedAt:  time.Now(),
		ModifiedAt: time.Now(),
		Account:    account,
	}

	if createTenantErr := s.tenantCreate.Execute(sessCtx, tenant); createTenantErr != nil {
		s.logger.Error("Failed creating tenant",
			slog.Any("error", createTenantErr))
		return createTenantErr
	}

	//
	// STEP 5: Create our administrator user account.
	//

	passwordHash, err := s.passwordProvider.GenerateHashFromPassword(walletPassword)
	if err != nil {
		s.logger.Error("Failed hashing password",
			slog.Any("error", err))
		return err
	}

	user := &domain.User{
		ID:                    userID,
		TenantID:              tenant.ID,
		TenantName:            tenantName,
		FirstName:             "System",
		LastName:              "Administrator",
		Name:                  "System Administrator",
		LexicalName:           "Administrator, System",
		Email:                 email,
		Status:                domain.UserStatusActive,
		PasswordHash:          passwordHash,
		PasswordHashAlgorithm: s.passwordProvider.AlgorithmName(),
		Role:                  domain.UserRoleRoot,
		WasEmailVerified:      true,
		CreatedAt:             time.Now(),
		ModifiedAt:            time.Now(),
	}

	if createUserErr := s.userCreate.Execute(sessCtx, user); err != nil {
		s.logger.Error("Failed creating user",
			slog.Any("error", createUserErr))
		return createUserErr
	}

	//
	// Step 6: For debugging puprposes only.
	//
	s.logger.Info("Gateway initialized",
		slog.Any("tenant_id", tenantID.Hex()),
		slog.Any("user_id", userID.Hex()),
		slog.Any("account_address", account.Address))

	return nil
}
