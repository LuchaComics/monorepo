package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
)

const (
	StatusActive   = 1
	StatusArchived = 2
)

type Collection struct {
	TenantID       primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	TenantName     string             `bson:"tenant_name" json:"tenant_name"`
	TenantTimezone string             `bson:"tenant_timezone" json:"tenant_timezone"`
	ID             primitive.ObjectID `bson:"_id" json:"id"`
	Status         int8               `bson:"status" json:"status"`
	CreatedAt      time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt     time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`

	// Name variable to track user name description of this collection.
	Name string `bson:"name" json:"name"`

	// SecretHashAlgorithm is the name of the hashing algorithm used for the secret value. Do not return in the API endpoint.
	SecretHashAlgorithm string `bson:"secret_hash_algorithm" json:"-"`

	// SecretHash is a random value which is returned to the client (so the client only knows the plaintext value) while
	// the server stores the hash of the value. This is used for our API token authentication for this specific collection.
	// Do not return in the API endpoint.
	SecretHash string `bson:"secret_hash" json:"-"`

	// ApiKey variable used only during the initial collection creation step
	// afterwords this value is not saved for security purposes.
	ApiKey string `bson:"-" json:"api_key,omitempty"`
}

type CollectionListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	TenantID        primitive.ObjectID
	UserID          primitive.ObjectID
	ExcludeArchived bool
	SearchText      string
}

type CollectionListResult struct {
	Results     []*Collection         `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// CollectionStorer Interface for tenant.
type CollectionStorer interface {
	Create(ctx context.Context, m *Collection) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Collection, error)
	GetByName(ctx context.Context, name string) (*Collection, error)
	GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*Collection, error)
	GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*Collection, error)
	UpdateByID(ctx context.Context, m *Collection) error
	ListByFilter(ctx context.Context, m *CollectionPaginationListFilter) (*CollectionPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *CollectionPaginationListFilter) ([]*CollectionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type CollectionAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type CollectionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) CollectionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("collections")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &CollectionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
