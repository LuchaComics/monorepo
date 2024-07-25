package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
)

const (
	StatusActive                        = 1
	StatusClaimed                       = 2
	StatusArchived                      = 3
	BusinessFunctionGrantFreeSubmission = 1
)

type Credit struct {
	StoreID                    primitive.ObjectID `bson:"store_id" json:"store_id"`
	StoreName                  string             `bson:"store_name" json:"store_name"`
	StoreTimezone              string             `bson:"store_timezone" json:"store_timezone"`
	ID                         primitive.ObjectID `bson:"_id" json:"id"`
	UserName                   string             `bson:"user_name" json:"user_name"`
	UserLexicalName            string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID                     primitive.ObjectID `bson:"user_id" json:"user_id"`
	BusinessFunction           int8               `bson:"business_function" json:"business_function"`
	OfferID                    primitive.ObjectID `bson:"offer_id" json:"offer_id"`
	OfferName                  string             `bson:"offer_name" json:"offer_name"`
	OfferServiceType           int8               `bson:"offer_service_type" json:"offer_service_type"`
	Status                     int8               `bson:"status" json:"status"`
	CreatedAt                  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt                 time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ClaimedByComicSubmissionID primitive.ObjectID `bson:"claimed_by_comic_submission_id" json:"claimed_by_comic_submission_id"`
}

type CreditListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	StoreID          primitive.ObjectID
	UserID           primitive.ObjectID
	OfferID          primitive.ObjectID
	OfferServiceType int8
	Status           int8
	BusinessFunction int8
	SearchText       string
}

type CreditListResult struct {
	Results     []*Credit          `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// CreditStorer Interface for store.
type CreditStorer interface {
	Create(ctx context.Context, m *Credit) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Credit, error)
	GetByName(ctx context.Context, name string) (*Credit, error)
	GetByPaymentProcessorCreditID(ctx context.Context, paymentProcessorCreditID string) (*Credit, error)
	GetNextAvailable(ctx context.Context, userID primitive.ObjectID, serviceType int8) (*Credit, error)
	UpdateByID(ctx context.Context, m *Credit) error
	ListByFilter(ctx context.Context, m *CreditPaginationListFilter) (*CreditPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *CreditPaginationListFilter) ([]*CreditAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type CreditAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type CreditStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) CreditStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("credits")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"user_name", "text"},
			{"user_lexical_name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &CreditStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
