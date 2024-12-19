package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	ComicSubmissionStatusInReview                       = 1
	ComicSubmissionStatusRejected                       = 2
	ComicSubmissionStatusAccepted                       = 3
	ComicSubmissionStatusError                          = 4
	ComicSubmissionStatusArchived                       = 5
	ComicSubmissionStatusFlagged                        = 6
	ComicSubmissionFlagIssueOther                       = 1
	ComicSubmissionFlagIssueDuplicateSubmission         = 2
	ComicSubmissionFlagIssuePoorImageQuality            = 3
	ComicSubmissionFlagIssueCounterfeit                 = 4
	ComicSubmissionFlagIssueInappropriateContent        = 5
	ComicSubmissionFlagActionDoNothing                  = 1
	ComicSubmissionFlagActionLockoutUser                = 2
	ComicSubmissionFlagActionLockoutUserAndBanIPAddress = 3
)

type ComicSubmission struct {
	ID                    primitive.ObjectID `bson:"_id" json:"id"`
	Name                  string             `bson:"name" json:"name"`
	FrontCover            *Attachment        `bson:"front_cover" json:"front_cover"`
	BackCover             *Attachment        `bson:"back_cover" json:"back_cover"`
	Status                int8               `bson:"status" json:"status"`
	UserID                primitive.ObjectID `bson:"user_id" json:"user_id"`
	CreatedAt             time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName     string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID       primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedFromIPAddress  string             `bson:"created_from_ip_address" json:"created_from_ip_address"`
	ModifiedAt            time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName    string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID      primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedFromIPAddress string             `bson:"modified_from_ip_address" json:"modified_from_ip_address"`
	CoinsReward           uint64             `bson:"coins_reward" json:"coins_reward"`
	WasAwarded            bool               `bson:"was_awarded" json:"was_awarded"`
	FlagIssue             int8               `bson:"flag_issue" json:"flag_issue"`
	FlagIssueOther        string             `bson:"flag_issue_other" json:"flag_issue_other"`
	FlagAction            int8               `bson:"flag_action" json:"flag_action"`
	TenantID              primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
}

type ComicSubmissionFilter struct {
	TenantID       primitive.ObjectID `json:"tenant_id"` // Required for data partitioning
	Name           *string            `json:"name,omitempty"`
	Status         int8               `json:"status,omitempty"`
	UserID         primitive.ObjectID `json:"user_id,omitempty"`
	CreatedAtStart *time.Time         `json:"created_at_start,omitempty"`
	CreatedAtEnd   *time.Time         `json:"created_at_end,omitempty"`

	// Cursor-based pagination
	LastID        *primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt *time.Time          `json:"last_created_at,omitempty"`
	Limit         int64               `json:"limit"`
}

type ComicSubmissionFilterResult struct {
	Submissions   []*ComicSubmission `json:"submissions"`
	HasMore       bool               `json:"has_more"`
	LastID        primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt time.Time          `json:"last_created_at,omitempty"`
}

// ComicSubmissionRepository Interface for a file that has content which lives in the cloud.
type ComicSubmissionRepository interface {
	Create(ctx context.Context, m *ComicSubmission) error
	CountByUserID(ctx context.Context, userID primitive.ObjectID) (uint64, error)
	CountByStatusAndByUserID(ctx context.Context, status int8, userID primitive.ObjectID) (uint64, error)
	CountTotalCreatedTodayByUserID(ctx context.Context, userID primitive.ObjectID, timezone string) (uint64, error)
	CountCoinsRewardByFilter(ctx context.Context, filter *ComicSubmissionFilter) (uint64, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*ComicSubmission, error)
	CountByFilter(ctx context.Context, filter *ComicSubmissionFilter) (uint64, error)
	ListByFilter(ctx context.Context, filter *ComicSubmissionFilter) (*ComicSubmissionFilterResult, error)
	UpdateByID(ctx context.Context, m *ComicSubmission) error
	TotalCoinsAwarded(ctx context.Context) (uint64, error)
	TotalCoinsAwardedByUserID(ctx context.Context, userID primitive.ObjectID) (uint64, error)
	// DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// ListByFilter(ctx context.Context, m *ComicSubmissionPaginationListFilter) (*ComicSubmissionPaginationListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *ComicSubmissionPaginationListFilter) ([]*ComicSubmissionAsSelectOption, error)
}

// type ComicSubmissionRepositoryImpl struct {
// 	Logger     *slog.Logger
// 	DbClient   *mongo.Client
// 	Collection *mongo.Collection
// }
//
// func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) ComicSubmissionRepository {
// 	// ctx := context.Background()
// 	uc := client.Database(appCfg.DB.Name).Collection("tenants")
//
// 	// The following few lines of code will create the index for our app for this
// 	// colleciton.
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{
// 			{"name", "text"},
// 		},
// 	}
// 	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
// 	if err != nil {
// 		// It is important that we crash the app on startup to meet the
// 		// requirements of `google/wire` framework.
// 		log.Fatal(err)
// 	}
//
// 	s := &ComicSubmissionRepositoryImpl{
// 		Logger:     loggerp,
// 		DbClient:   client,
// 		Collection: uc,
// 	}
// 	return s
// }
