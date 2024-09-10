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

type NFTCollection struct {
	TenantID              primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	TenantName            string             `bson:"tenant_name" json:"tenant_name"`
	TenantTimezone        string             `bson:"tenant_timezone" json:"tenant_timezone"`
	ID                    primitive.ObjectID `bson:"_id" json:"id"`
	Status                int8               `bson:"status" json:"status"`
	CreatedAt             time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedFromIPAddress  string             `bson:"created_from_ip_address" json:"created_from_ip_address,omitempty"`
	ModifiedAt            time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedFromIPAddress string             `bson:"modified_from_ip_address" json:"modified_from_ip_address,omitempty"`
	Name                  string             `bson:"name" json:"name"`                               // Internal name for staff use; not displayed to customers.
	IPNSKeyName           string             `bson:"ipns_key_name" json:"ipns_key_name"`             // Unique key name specific to this collection.
	IPNSName              string             `bson:"ipns_name" json:"ipns_name"`                     // IPNS address of this collection.
	IPFSDirectoryName     string             `bson:"ipfs_directory_name" json:"ipfs_directory_name"` // Directory name for storing NFTs.
	IPFSDirectoryCID      string             `bson:"ipfs_directory_cid" json:"ipfs_directory_cid"`   // CID of the directory for storing NFTs.
	TokensCount           uint64             `bson:"tokens_count" json:"tokens_count"`               // Number of tokens in this collection.
	MetadataFileCIDs      map[uint64]string  `bson:"metadata_file_cids" json:"metadata_file_cids"`   // Used for mapping TokenID's to their metadata file CID.
}

type NFTCollectionListFilter struct {
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

type NFTCollectionListResult struct {
	Results     []*NFTCollection   `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// NFTCollectionStorer Interface for tenant.
type NFTCollectionStorer interface {
	Create(ctx context.Context, m *NFTCollection) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*NFTCollection, error)
	GetByName(ctx context.Context, name string) (*NFTCollection, error)
	GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*NFTCollection, error)
	GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*NFTCollection, error)
	UpdateByID(ctx context.Context, m *NFTCollection) error
	ListByFilter(ctx context.Context, m *NFTCollectionPaginationListFilter) (*NFTCollectionPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *NFTCollectionPaginationListFilter) ([]*NFTCollectionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type NFTCollectionAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type NFTCollectionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) NFTCollectionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("nft_collections")

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

	s := &NFTCollectionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
