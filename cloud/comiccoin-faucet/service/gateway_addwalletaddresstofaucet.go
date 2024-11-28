package service

import (
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayAddWalletAddressToFaucetService struct {
	logger               *slog.Logger
	tenantGetByIDUseCase *usecase.TenantGetByIDUseCase
	userGetByIDUseCase   *usecase.UserGetByIDUseCase
	userUpdateUseCase    *usecase.UserUpdateUseCase
}

func NewGatewayAddWalletAddressToFaucetService(
	logger *slog.Logger,
	uc1 *usecase.TenantGetByIDUseCase,
	uc2 *usecase.UserGetByIDUseCase,
	uc3 *usecase.UserUpdateUseCase,
) *GatewayAddWalletAddressToFaucetService {
	return &GatewayAddWalletAddressToFaucetService{logger, uc1, uc2, uc3}
}

type GatewayProfileWalletAddressRequestIDO struct {
	WalletAddress string `bson:"wallet_address" json:"wallet_address"`
}

func (s *GatewayAddWalletAddressToFaucetService) Execute(sessCtx mongo.SessionContext, req *GatewayProfileWalletAddressRequestIDO) (*domain.User, error) {
	//
	// STEP 1: Get from session and related records.
	//

	// Extract from our session the following data.
	userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := s.userGetByIDUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if ou == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	t, err := s.tenantGetByIDUseCase.Execute(sessCtx, ou.TenantID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if t == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("tenant_id", "does not exist")
	}

	//
	// STEP 2: Validation of input.
	//

	e := make(map[string]string)
	if req.WalletAddress == "" {
		e["wallet_address"] = "missing value"
	}
	if ou.WalletAddress != nil {
		// If user has wallet address, make sure it's not the nil address.
		if ou.WalletAddress.Hex() != "0x0000000000000000000000000000000000000000" {
			e["wallet_address"] = fmt.Sprintf("already set: %v", ou.WalletAddress.Hex())
		}
		//TODO: LastCoinsDepositAt time.Time
	}
	if t.Account.Balance == 0 {
		e["message"] = "Faucet has no coins in wallet"

		s.logger.Error("Wallet has empty balance",
			slog.String("address", t.Account.Address.Hex()))
	}

	if len(e) != 0 {
		s.logger.Warn("Failed validation login",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3: Set wallet address.
	//

	// walletAddress := common.HexToAddress(strings.ToLower(req.WalletAddress))
	// ou.WalletAddress = &walletAddress
	//
	// if err := s.userUpdateUseCase.Execute(sessCtx, ou); err != nil {
	// 	s.logger.Error("user update by id error", slog.Any("error", err))
	// 	return nil, err
	// }

	//
	// STEP 4: Transfer coins into new wallet address.
	//

	//TODO: IMPL.

	return ou, nil
}
