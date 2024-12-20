package service

import (
	"errors"
	"fmt"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type ComicSubmissionCreateService struct {
	logger                                             *slog.Logger
	userGetByIDUseCase                                 *usecase.UserGetByIDUseCase
	comicSubmissionCountTotalCreatedTodayByUserUseCase *usecase.ComicSubmissionCountTotalCreatedTodayByUserUseCase
	attachmentGetUseCase                               *usecase.AttachmentGetUseCase
	attachmentUpdateUseCase                            *usecase.AttachmentUpdateUseCase
	comicSubmissionCreateUseCase                       *usecase.ComicSubmissionCreateUseCase
}

func NewComicSubmissionCreateService(
	logger *slog.Logger,
	uc1 *usecase.UserGetByIDUseCase,
	uc2 *usecase.ComicSubmissionCountTotalCreatedTodayByUserUseCase,
	uc3 *usecase.AttachmentGetUseCase,
	uc4 *usecase.AttachmentUpdateUseCase,
	uc5 *usecase.ComicSubmissionCreateUseCase,
) *ComicSubmissionCreateService {
	return &ComicSubmissionCreateService{logger, uc1, uc2, uc3, uc4, uc5}
}

type ComicSubmissionCreateRequestIDO struct {
	Name       string             `bson:"name" json:"name"`
	FrontCover primitive.ObjectID `bson:"front_cover" json:"front_cover"`
	BackCover  primitive.ObjectID `bson:"back_cover" json:"back_cover"`
}

type ComicSubmissionCreateResponseIDO domain.ComicSubmission

func (s *ComicSubmissionCreateService) Execute(sessCtx mongo.SessionContext, req *ComicSubmissionCreateRequestIDO) (*ComicSubmissionCreateResponseIDO, error) {
	//
	// STEP 1: Validation
	//

	e := make(map[string]string)

	if req == nil {
		err := errors.New("No request data inputted")
		s.logger.Error("validation error", slog.Any("err", err))
		return nil, err
	}

	s.logger.Debug("Validating...",
		slog.Any("name", req.Name),
		slog.Any("front_cover", req.FrontCover),
		slog.Any("back_cover", req.BackCover))

	if req.Name == "" {
		e["name"] = "Name is required"
	}
	if req.FrontCover.IsZero() {
		e["front_cover"] = "Front cover is required"
	}
	if req.BackCover.IsZero() {
		e["back_cover"] = "Back cover is required"
	}
	if len(e) != 0 {
		s.logger.Warn("validation failure",
			slog.Any("e", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 2: Fetch related records.
	//

	// Extract from our session the following data.
	userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByIDUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("user_id", "does not exist")
	}
	frontCover, err := s.attachmentGetUseCase.Execute(sessCtx, req.FrontCover)
	if err != nil {
		s.logger.Error("Failed retrieving front cover attachment",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("err", err))
		return nil, err
	}
	if frontCover == nil {
		err := fmt.Errorf("Front cover does not exist for: %v", req.FrontCover)
		s.logger.Error("Failed getting front cover",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("front_cover", req.FrontCover),
			slog.Any("error", err))
		return nil, err
	}
	backCover, err := s.attachmentGetUseCase.Execute(sessCtx, req.BackCover)
	if err != nil {
		s.logger.Error("Failed retrieving back cover attachment",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("err", err))
		return nil, err
	}
	if backCover == nil {
		err := fmt.Errorf("Back cover does not exist for: %v", req.BackCover)
		s.logger.Error("Failed getting back cover",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("back_cover", req.BackCover),
			slog.Any("error", err))
		return nil, err
	}

	//
	// STEP 3: If user is not verified, then enforce daily limit restriction.
	//

	if !u.WasVerified {
		todaysCount, err := s.comicSubmissionCountTotalCreatedTodayByUserUseCase.Execute(sessCtx, u.ID, u.Timezone)
		if err != nil {
			s.logger.Error("Failed creating comic submission",
				slog.Any("user_id", userID),
				slog.Any("name", req.Name),
				slog.Any("err", err))
			return nil, err
		}
		if todaysCount > 6 {
			e["message"] = "Daily limit of 4 reached today"
			s.logger.Warn("validation failure",
				slog.Any("e", e))
			return nil, httperror.NewForBadRequest(&e)
		}
	}

	//
	// STEP 4: Create our unique identifier for our comic submission.
	//

	comicSubmissionID := primitive.NewObjectID()

	//
	// STEP 5: Attach the front and back covers to our comic submission.
	//

	frontCover.BelongsToUniqueIdentifier = comicSubmissionID
	frontCover.BelongsToType = domain.AttachmentBelongsToTypeSubmission
	frontCover.ModifiedAt = time.Now()
	frontCover.ModifiedByUserName = u.Name
	frontCover.ModifiedByUserID = u.ID
	frontCover.ModifiedFromIPAddress = ipAddress
	if err := s.attachmentUpdateUseCase.Execute(sessCtx, frontCover); err != nil {
		s.logger.Error("Failed updating attachment",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("err", err))
		return nil, err
	}

	backCover.BelongsToUniqueIdentifier = comicSubmissionID
	backCover.BelongsToType = domain.AttachmentBelongsToTypeSubmission
	backCover.ModifiedAt = time.Now()
	backCover.ModifiedByUserName = u.Name
	backCover.ModifiedByUserID = u.ID
	backCover.ModifiedFromIPAddress = ipAddress
	if err := s.attachmentUpdateUseCase.Execute(sessCtx, backCover); err != nil {
		s.logger.Error("Failed updating attachment",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("err", err))
		return nil, err
	}

	//
	// STEP 6: Create our comic submission.
	//

	comicSubmission := &domain.ComicSubmission{
		ID:                    comicSubmissionID,
		Name:                  req.Name,
		FrontCover:            frontCover,
		BackCover:             backCover,
		Status:                domain.ComicSubmissionStatusInReview,
		UserID:                u.ID,
		CreatedAt:             time.Now(),
		CreatedByUserName:     u.Name,
		CreatedByUserID:       u.ID,
		CreatedFromIPAddress:  ipAddress,
		ModifiedAt:            time.Now(),
		ModifiedByUserName:    u.Name,
		ModifiedByUserID:      u.ID,
		ModifiedFromIPAddress: ipAddress,
		CoinsReward:           5, //TODO: Replace magic number with configuration variable.
		TenantID:              u.TenantID,
	}
	if err := s.comicSubmissionCreateUseCase.Execute(sessCtx, comicSubmission); err != nil {
		s.logger.Error("Failed creating comic submission",
			slog.Any("user_id", userID),
			slog.Any("name", req.Name),
			slog.Any("err", err))
		return nil, err
	}

	return (*ComicSubmissionCreateResponseIDO)(comicSubmission), nil
}
