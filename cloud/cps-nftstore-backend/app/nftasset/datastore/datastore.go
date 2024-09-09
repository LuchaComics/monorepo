package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config"
)

// "queued" "pinning" "pinned" "failed"

const (
	StatusPending = "pending"
	StatusQueued  = "queued"
	StatusPinning = "pinning"
	StatusPinned  = "pinned"
	StatusFailed  = "failed"

	ContentTypeFile  = 1
	ContentTypeImage = 2
)

// NFTAsset is a representation of a pin request. It means it is the IPFS content which we are saving in our system and sharing to the IPFS network, also know as "pinning". This structure has the core variables required to work with IPFS as per their documentation https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers, in additon we also have our applications specific varaibles.
type NFTAsset struct {
	// Status represents operational state of this pinned object in IPFS.
	Status string `bson:"status" json:"status"`

	// CID variable is the unique identifier of our content on the IPFS network. The official definition is: Content Identifier (CID) points at the root of a DAG that is pinned recursively.
	CID string `bson:"cid" json:"cid"`

	// Name variable used to provide human readable description for the content. This is optional.
	Name string `bson:"name" json:"name"`

	// The date/time this content was pinned in IPFS network.
	CreatedAt time.Time `bson:"created,omitempty" json:"created,omitempty"`

	Filename    string `bson:"filename" json:"filename,omitempty"`
	ContentType string `bson:"content_type" json:"content_type,omitempty"`

	TenantID       primitive.ObjectID `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	TenantName     string             `bson:"tenant_name" json:"tenant_name,omitempty"`
	TenantTimezone string             `bson:"tenant_timezone" json:"tenant_timezone,omitempty"`

	// ID variable is the unique identifier we use internally in our system.
	ID                    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CreatedFromIPAddress  string             `bson:"created_from_ip_address" json:"created_from_ip_address,omitempty"`
	ModifiedAt            time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedFromIPAddress string             `bson:"modified_from_ip_address" json:"modified_from_ip_address,omitempty"`

	// NFTMetadataID variable used to track ownership of this pin object. Used primarily for organization purposes.
	NFTMetadataID primitive.ObjectID `bson:"nftmetadata_id" json:"nftmetadata_id,omitempty"`
}

type NFTAssetAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// NFTAssetStorer Interface for nftasset.
type NFTAssetStorer interface {
	Create(ctx context.Context, m *NFTAsset) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*NFTAsset, error)
	GetByCID(ctx context.Context, cid string) (*NFTAsset, error)
	GetAllCIDs(ctx context.Context) ([]string, error)
	UpdateByID(ctx context.Context, m *NFTAsset) error
	ListByFilter(ctx context.Context, m *NFTAssetPaginationListFilter) (*NFTAssetPaginationListResult, error)
	ListByNFTMetadataID(ctx context.Context, nftmetadataID primitive.ObjectID) (*NFTAssetPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *NFTAssetPaginationListFilter) ([]*NFTAssetAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	DeleteByCID(ctx context.Context, cid primitive.ObjectID) error
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type NFTAssetStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) NFTAssetStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("nft_assets")

	// For debugging purposes only or if you are going to recreate new indexes.
	if _, err := uc.Indexes().DropAll(context.TODO()); err != nil {
		loggerp.Warn("failed deleting all indexes",
			slog.Any("err", err))

		// Do not crash app, just continue.
	}

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "cid", Value: 1}}}, // Used for ipfs gateway `get` requests
		{Keys: bson.D{{Key: "nftmetadata_id", Value: 1}, {Key: "created", Value: SortOrderDescending}}}, // Note: Used in default list.

		// 4. Compound Text Index for Text Search
		{Keys: bson.D{
			// Frontend Search  - filter by customer
			{Key: "cid", Value: "text"},
			{Key: "name", Value: "text"},
			{Key: "meta", Value: "text"},
			{Key: "info", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &NFTAssetStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
