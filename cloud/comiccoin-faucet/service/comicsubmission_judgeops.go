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
	cloudStorageDeleteUseCase     *usecase.CloudStorageDeleteUseCase
	userGetByIDUseCase            *usecase.UserGetByIDUseCase
	comicSubmissionGetByIDUseCase *usecase.ComicSubmissionGetByIDUseCase
	comicSubmissionUpdateUseCase  *usecase.ComicSubmissionUpdateUseCase
}

func NewComicSubmissionJudgeOperationService(
	cfg *config.Configuration,
	logger *slog.Logger,
	s1 *FaucetCoinTransferService,
	uc1 *usecase.CloudStorageDeleteUseCase,
	uc2 *usecase.UserGetByIDUseCase,
	uc3 *usecase.ComicSubmissionGetByIDUseCase,
	uc4 *usecase.ComicSubmissionUpdateUseCase,
) *ComicSubmissionJudgeOperationService {
	return &ComicSubmissionJudgeOperationService{cfg, logger, s1, uc1, uc2, uc3, uc4}
}

type ComicSubmissionJudgeVerdictRequestIDO struct {
	ComicSubmissionID  primitive.ObjectID `bson:"comic_submission_id" json:"comic_submission_id"`
	Status             int8               `bson:"status" json:"status"`
	FlagIssue          int8               `bson:"flag_issue" json:"flag_issue"`
	FlagIssueOther     string             `bson:"flag_issue_other" json:"flag_issue_other"`
	FlagAction         int8               `bson:"flag_action" json:"flag_action"`
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
	} else {
		if req.Status == domain.ComicSubmissionStatusFlagged {
			if req.FlagIssue == 0 {
				e["flag_issue"] = "Flag issue is required"
			} else if req.FlagIssue == domain.ComicSubmissionFlagIssueOther && req.FlagIssueOther == "" {
				e["flag_issue_other"] = "Flag issue (other) is required"
			}
			if req.FlagAction == 0 {
				e["flag_action"] = "Flag action is required"
			}
		}
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
		s.logger.Debug("Faucet is granting user some ComicCoins...",
			slog.Any("wallet_address", customerUser.WalletAddress),
			slog.Any("comiccoins_rewarded", comicSubmission.CoinsReward))

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
		s.logger.Debug("Faucet is granted user some ComicCoins",
			slog.Any("comiccoins_rewarded", comicSubmission.CoinsReward))

		// Update the comic submission to indicate we successfully sent
		// the reward.
		comicSubmission.WasAwarded = true
	}

	//
	// STEP 4: Status is flagged
	//

	if req.Status == domain.ComicSubmissionStatusFlagged {
		s.logger.Debug("Executing flag issue",
			slog.Any("flag_issue", req.FlagIssue),
		)

		// TODO: Save the hash value to block future images by

		if err := s.cloudStorageDeleteUseCase.Execute(sessCtx, []string{comicSubmission.FrontCover.ObjectKey}); err != nil {
			s.logger.Error("Failed deleting front cover attachment from cloud storage",
				slog.Any("err", err))
			return nil, err
		}
		if err := s.cloudStorageDeleteUseCase.Execute(sessCtx, []string{comicSubmission.BackCover.ObjectKey}); err != nil {
			s.logger.Error("Failed deleting front cover attachment from cloud storage",
				slog.Any("err", err))
			return nil, err
		}
	}

	//
	// STEP 5: Update the state in the database.
	//

	comicSubmission.Status = req.Status
	comicSubmission.FlagIssue = req.FlagIssue
	comicSubmission.FlagIssueOther = req.FlagIssueOther
	comicSubmission.FlagAction = req.FlagAction
	comicSubmission.ModifiedAt = time.Now()
	comicSubmission.ModifiedByUserName = adminUser.Name
	comicSubmission.ModifiedByUserID = req.AdminUserID
	comicSubmission.ModifiedFromIPAddress = req.AdminUserIPAddress
	if err := s.comicSubmissionUpdateUseCase.Execute(sessCtx, comicSubmission); err != nil {
		s.logger.Error("Failed update",
			slog.Any("err", err))
		return nil, err
	}

	//
	// STEP 6: Execute flag action (if submission was flagged),
	//

	if req.Status == domain.ComicSubmissionStatusFlagged && req.FlagAction > domain.ComicSubmissionFlagActionDoNothing {
		s.logger.Debug("Executing flag action",
			slog.Any("flag_action", req.FlagAction),
		)
	}

	// s.logger.Debug("fetched",
	// 	slog.Any("id", id),
	// 	slog.Any("detail", detail))

	return comicSubmission, nil
}
