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
	StatusQueued  = "queued"
	StatusPinning = "pinning"
	StatusPinned  = "pinned"
	StatusFailed  = "failed"
)

// PinObject is a representation of a pin request. It means it is the IPFS content which we are saving in our system and sharing to the IPFS network, also know as "pinning".
// This structure is not exactly as designed in the IPFS documentation https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers, this structure
// is minified and light on purpose for our purpose.
type PinObject struct {
	ID          primitive.ObjectID `bson:"_id" json:"id,omitempty"`    // ID variable is the unique identifier we use internally in our system.
	IPNSPath    string             `bson:"ipns_path" json:"ipns_path"` // Optional variable which is set if this file is mounted to the IPNS.
	CID         string             `bson:"cid" json:"cid"`
	Content     []byte             `bson:"-" json:"content,omitempty"` // FileContent variable holds all the content of this pin. Variable will not be saved to database, only returned in API endpoint.
	Filename    string             `bson:"filename" json:"filename,omitempty"`
	ContentType string             `bson:"content_type" json:"content_type,omitempty"`
	CreatedAt   time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"` // The date/time this content was pinned in IPFS network.
	ModifiedAt  time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
}

type PinObjectAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// PinObjectStorer Interface for pinobject.
type PinObjectStorer interface {
	Create(ctx context.Context, m *PinObject) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*PinObject, error)
	GetByCID(ctx context.Context, cid string) (*PinObject, error)
	GetByIPNSPath(ctx context.Context, ipnsPath string) (*PinObject, error)
	GetAllCIDs(ctx context.Context) ([]string, error)
	UpdateByID(ctx context.Context, m *PinObject) error
	ListByFilter(ctx context.Context, m *PinObjectPaginationListFilter) (*PinObjectPaginationListResult, error)
	ListByTenantID(ctx context.Context, tenantID primitive.ObjectID) (*PinObjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *PinObjectPaginationListFilter) ([]*PinObjectAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	DeleteByCID(ctx context.Context, cid string) error
	// //TODO: Add more...
}

type PinObjectStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) PinObjectStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("pin_objects")

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
		{Keys: bson.D{
			// Frontend Search  - filter by customer
			{Key: "cid", Value: "text"},
			{Key: "Filename", Value: "text"},
		}},
	})
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &PinObjectStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
