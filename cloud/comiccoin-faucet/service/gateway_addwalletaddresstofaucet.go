package service

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayAddWalletAddressToFaucetService struct {
	config                    *config.Configuration
	logger                    *slog.Logger
	tenantGetByIDUseCase      *usecase.TenantGetByIDUseCase
	userGetByIDUseCase        *usecase.UserGetByIDUseCase
	userUpdateUseCase         *usecase.UserUpdateUseCase
	faucetCoinTransferService *FaucetCoinTransferService
}

func NewGatewayAddWalletAddressToFaucetService(
	cfg *config.Configuration,
	logger *slog.Logger,
	uc1 *usecase.TenantGetByIDUseCase,
	uc2 *usecase.UserGetByIDUseCase,
	uc3 *usecase.UserUpdateUseCase,
	s1 *FaucetCoinTransferService,
) *GatewayAddWalletAddressToFaucetService {
	return &GatewayAddWalletAddressToFaucetService{cfg, logger, uc1, uc2, uc3, s1}
}

type GatewayProfileWalletAddressRequestIDO struct {
	WalletAddress string `bson:"wallet_address" json:"wallet_address"`
}

func (s *GatewayAddWalletAddressToFaucetService) Execute(
	sessCtx mongo.SessionContext,
	req *GatewayProfileWalletAddressRequestIDO,
) (*domain.User, error) {
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
		e["wallet_address"] = "Wallet address is required"
	}
	if ou.WalletAddress != nil {
		// If user has wallet address, make sure it's not the nil address.
		if ou.WalletAddress.Hex() != "0x0000000000000000000000000000000000000000" {
			e["wallet_address"] = fmt.Sprintf("Wallet address already set: %v", ou.WalletAddress.Hex())
		}
		//TODO: LastCoinsDepositAt time.Time
	}
	if t.Account.Balance == 0 {
		e["message"] = "Faucet has no coins in wallet - please try again later"

		s.logger.Error("Wallet has empty balance",
			slog.String("address", t.Account.Address.Hex()))
	}

	if len(e) != 0 {
		s.logger.Warn("Failed validation login",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3: Transfer coins into new wallet address.
	//

	walletAddress := common.HexToAddress(strings.ToLower(req.WalletAddress))
	idoReq := &FaucetCoinTransferRequestIDO{
		ChainID:               s.config.Blockchain.ChainID,
		FromAccountAddress:    s.config.App.WalletAddress,
		AccountWalletPassword: s.config.App.WalletPassword,
		To:                    &walletAddress,
		Value:                 s.config.App.RegistrationCoinsReward,
		Data:                  make([]byte, 0),
	}
	if err := s.faucetCoinTransferService.Execute(sessCtx, idoReq); err != nil {
		s.logger.Error("Failed transfering coins to address",
			slog.Any("error", err))
		return nil, err
	}

	s.logger.Info("ComicCoin Faucet submitted coins to wallet address successfully",
		slog.String("user_id", ou.ID.Hex()),
		slog.String("address", req.WalletAddress),
	)

	//
	// STEP 4: Set wallet address.
	//

	ou.WalletAddress = &walletAddress
	ou.LastCoinsDepositAt = time.Now()
	if err := s.userUpdateUseCase.Execute(sessCtx, ou); err != nil {
		s.logger.Error("user update by id error", slog.Any("error", err))
		return nil, err
	}

	s.logger.Info("ComicCoin Faucet collected wallet address of user",
		slog.String("user_id", ou.ID.Hex()),
		slog.String("address", req.WalletAddress),
	)

	return ou, nil
}
