package usecase

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
)

type TenantCreateUseCase struct {
	config *config.Configuration
	logger *slog.Logger
	repo   domain.TenantRepository
}

func NewTenantCreateUseCase(config *config.Configuration, logger *slog.Logger, repo domain.TenantRepository) *TenantCreateUseCase {
	return &TenantCreateUseCase{config, logger, repo}
}

func (uc *TenantCreateUseCase) Execute(ctx context.Context, tenant *domain.Tenant) error {
	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if tenant == nil {
		e["tenant"] = "missing value"
	} else {
		//TODO: IMPL.
	}
	if len(e) != 0 {
		uc.logger.Warn("Validation failed for upsert",
			slog.Any("error", e))
		return httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Check our strucutre.
	//

	if checkExists, err := uc.repo.CheckIfExistsByID(ctx, tenant.ID); err != nil || checkExists {
		if err != nil {
			uc.logger.Error("Failed checking if tenant exists",
				slog.Any("tenant_id", tenant.ID),
				slog.Any("error", err),
			)
			return err
		}
		exErr := fmt.Errorf("Tenant already exists for id: %v\n", tenant.ID)
		uc.logger.Error("Failed checking if tenant exists",
			slog.Any("tenant_id", tenant.ID),
			slog.Any("error", exErr))
		return exErr
	}

	//
	// STEP 3: Insert into database.
	//

	return uc.repo.Create(ctx, tenant)
}
