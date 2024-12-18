package service

import (
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionJudgeOperationService struct {
	config                        *config.Configuration
	logger                        *slog.Logger
	faucetCoinTransferService     *FaucetCoinTransferService
	userGetByIDUseCase            *usecase.UserGetByIDUseCase
	comicSubmissionGetByIDUseCase *usecase.ComicSubmissionGetByIDUseCase
	comicSubmissionUpdateUseCase  *usecase.ComicSubmissionUpdateUseCase
}

func NewComicSubmissionJudgeOperationService(
	cfg *config.Configuration,
	logger *slog.Logger,
	s1 *FaucetCoinTransferService,
	uc1 *usecase.UserGetByIDUseCase,
	uc2 *usecase.ComicSubmissionGetByIDUseCase,
	uc3 *usecase.ComicSubmissionUpdateUseCase,
) *ComicSubmissionJudgeOperationService {
	return &ComicSubmissionJudgeOperationService{cfg, logger, s1, uc1, uc2, uc3}
}

type ComicSubmissionJudgeVerdictRequestIDO struct {
	ComicSubmissionID  primitive.ObjectID `bson:"comic_submission_id" json:"comic_submission_id"`
	Status             int8               `bson:"status" json:"status"`
	AdminUserID        primitive.ObjectID `bson:"admin_user_id" json:"admin_user_id"`
	AdminUserIPAddress string             `bson:"admin_user_ip_address" json:"admin_user_ip_address"`
}

func (s *ComicSubmissionJudgeOperationService) Execute(
	sessCtx mongo.SessionContext,
	req *ComicSubmissionJudgeVerdictRequestIDO,
) (*domain.ComicSubmission, error) {
	s.logger.Warn("Begin to validate",
		slog.Any("ComicSubmissionID", req.ComicSubmissionID),
		slog.Any("Status", req.Status),
		slog.Any("AdminUserID", req.AdminUserID),
		slog.Any("AdminUserIPAddress", req.AdminUserIPAddress),
	)

	//
	// STEP 1: Validation.
	//

	e := make(map[string]string)
	if req.ComicSubmissionID.IsZero() {
		e["comic_submission_id"] = "Comic submission identifier is required"
	}
	if req.Status == 0 {
		e["status"] = "Status is required"
	}
	if req.AdminUserID.IsZero() {
		e["admin_user_id"] = "Admin user identifier is required"
	}
	if req.AdminUserIPAddress == "" {
		e["admin_user_ip_address"] = "Admin user IP address is required"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validating",
			slog.Any("req", req),
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Get related records.
	//

	adminUser, err := s.userGetByIDUseCase.Execute(sessCtx, req.AdminUserID)
	if err != nil {
		s.logger.Error("failed getting admin",
			slog.Any("err", err))
		return nil, err
	}
	if adminUser == nil {
		err := fmt.Errorf("Administrative user does not exist with ID: %v", req.AdminUserID)
		s.logger.Error("failed getting admin",
			slog.Any("err", err))
		return nil, err
	}

	comicSubmission, err := s.comicSubmissionGetByIDUseCase.Execute(sessCtx, req.ComicSubmissionID)
	if err != nil {
		s.logger.Error("failed getting comic submission",
			slog.Any("err", err))
		return nil, err
	}
	if comicSubmission == nil {
		err := fmt.Errorf("Comic submission does not exist with ID: %v", req.ComicSubmissionID)
		s.logger.Error("failed getting comic submission",
			slog.Any("err", err))
		return nil, err
	}

	customerUser, err := s.userGetByIDUseCase.Execute(sessCtx, comicSubmission.UserID)
	if err != nil {
		s.logger.Error("failed getting customer",
			slog.Any("err", err))
		return nil, err
	}
	if customerUser == nil {
		err := fmt.Errorf("Customer user does not exist with ID: %v", comicSubmission.UserID)
		s.logger.Error("failed getting customer",
			slog.Any("err", err))
		return nil, err
	}

	//
	// STEP 3: Reward the user if approved without previous reward.
	//

	if req.Status == domain.ComicSubmissionStatusAccepted && !comicSubmission.WasAwarded {
		req := &FaucetCoinTransferRequestIDO{
			ChainID:               s.config.Blockchain.ChainID,
			FromAccountAddress:    s.config.App.WalletAddress,
			AccountWalletPassword: s.config.App.WalletPassword,
			To:                    customerUser.WalletAddress,
			Value:                 comicSubmission.CoinsReward,
			Data:                  ([]byte)(s.config.App.FrontendDomain),
		}
		if err := s.faucetCoinTransferService.Execute(sessCtx, req); err != nil {
			s.logger.Error("Failed faucet coin transfer",
				slog.Any("err", err))
			return nil, err
		}
		s.logger.Debug("Granting user ComicCoins",
			slog.Any("comiccoins_rewarded", comicSubmission.CoinsReward),
		)

		// Update the comic submission to indicate we successfully sent
		// the reward.
		comicSubmission.WasAwarded = true
	}

	//
	// STEP 4: Update the state in the database.
	//

	comicSubmission.Status = req.Status
	comicSubmission.ModifiedAt = time.Now()
	comicSubmission.ModifiedByUserName = adminUser.Name
	comicSubmission.ModifiedByUserID = req.AdminUserID
	comicSubmission.ModifiedFromIPAddress = req.AdminUserIPAddress
	if err := s.comicSubmissionUpdateUseCase.Execute(sessCtx, comicSubmission); err != nil {
		s.logger.Error("Failed update",
			slog.Any("err", err))
		return nil, err
	}

	// s.logger.Debug("fetched",
	// 	slog.Any("id", id),
	// 	slog.Any("detail", detail))

	return comicSubmission, nil
}
