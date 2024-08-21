package controller

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	s3_storage "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/adapter/storage/s3"
	attachment_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	domain "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/attachment/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
	"github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/provider/uuid"
)

// AttachmentController Interface for attachment business logic controller.
type AttachmentController interface {
	Create(ctx context.Context, req *AttachmentCreateRequestIDO) (*domain.Attachment, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*domain.Attachment, error)
	UpdateByID(ctx context.Context, ns *AttachmentUpdateRequestIDO) (*domain.Attachment, error)
	ListByFilter(ctx context.Context, f *domain.AttachmentPaginationListFilter) (*domain.AttachmentPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *domain.AttachmentPaginationListFilter) ([]*domain.AttachmentAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	PermanentlyDeleteByID(ctx context.Context, id primitive.ObjectID) error
}

type AttachmentControllerImpl struct {
	Config           *config.Conf
	Logger           *slog.Logger
	UUID             uuid.Provider
	S3               s3_storage.S3Storager
	DbClient         *mongo.Client
	AttachmentStorer attachment_s.AttachmentStorer
	UserStorer       user_s.UserStorer
}

func NewController(
	appCfg *config.Conf,
	loggerp *slog.Logger,
	uuidp uuid.Provider,
	s3 s3_storage.S3Storager,
	client *mongo.Client,
	org_storer attachment_s.AttachmentStorer,
	usr_storer user_s.UserStorer,
) AttachmentController {
	s := &AttachmentControllerImpl{
		Config:           appCfg,
		Logger:           loggerp,
		UUID:             uuidp,
		S3:               s3,
		DbClient:         client,
		AttachmentStorer: org_storer,
		UserStorer:       usr_storer,
	}
	s.Logger.Debug("attachment controller initialization started...")
	s.Logger.Debug("attachment controller initialized")
	return s
}
