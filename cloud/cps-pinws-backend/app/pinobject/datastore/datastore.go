package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/config"
)

// "queued" "pinning" "pinned" "failed"

const (
	StatusQueued  = "queued"
	StatusPinning = "pinning"
	StatusPinned  = "pinned"
	StatusFailed  = "failed"

	ContentTypeFile  = 1
	ContentTypeImage = 2
)

// PinObject is a representation of a pin request. It means it is the IPFS content which we are saving in our system and sharing to the IPFS network, also know as "pinning". This structure has the core variables required to work with IPFS as per their documentation https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers, in additon we also have our applications specific varaibles.
type PinObject struct {
	// RequestID variable is the public viewable unique identifier of this pin.
	RequestID primitive.ObjectID `bson:"requestid" json:"requestid"`

	// Status represents operational state of this pinned object in IPFS.
	Status string `bson:"status" json:"status"`

	// CID variable is the unique identifier of our content on the IPFS network. The official definition is: Content Identifier (CID) points at the root of a DAG that is pinned recursively.
	CID string `bson:"cid" json:"cid"`

	// Name variable used to provide human readable description for the content. This is optional.
	Name string `bson:"name" json:"name"`

	// The date/time this content was pinned in IPFS network. Developers note: Normally we write it as `CreatedAt`, but IPFS specs require us to write it this way.
	Created time.Time `bson:"created,omitempty" json:"created,omitempty"`

	// Addresses provided in origins list are relevant only during the initial pinning, and don't need to be persisted by the pinning service
	Origins []string `bson:"origins" json:"origins"`

	// Any additional vendor-specific information is returned in optional info.
	Meta map[string]string `bson:"meta" json:"meta"`

	// Addresses in the delegates array are peers designated by pinning service that will receive the pin data over bitswap
	Delegates []string `bson:"delegates" json:"delegates"`

	// Any additional vendor-specific information is returned in optional info.
	Info map[string]string `bson:"info" json:"info"`

	// CODE BELOW IS EXTENSION TO THE IPFS SPEC.
	//------------------------------------------//

	// FileContent variable holds all the content of this pin. Variable will not be saved to database, only returned in API endpoint.
	Content []byte `bson:"-" json:"content,omitempty"`

	Filename    string `bson:"filename" json:"filename,omitempty"`
	ObjectKey   string `bson:"object_key" json:"object_key,omitempty"`
	ObjectURL   string `bson:"object_url" json:"object_url,omitempty"`
	ContentType int8   `bson:"content_type" json:"content_type,omitempty"`

	TenantID       primitive.ObjectID `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	TenantName     string             `bson:"tenant_name" json:"tenant_name,omitempty"`
	TenantTimezone string             `bson:"tenant_timezone" json:"tenant_timezone,omitempty"`

	// ID variable is the unique identifier we use internally in our system.
	ID                    primitive.ObjectID `bson:"_id" json:"id,omitempty"`
	CreatedFromIPAddress  string             `bson:"created_from_ip_address" json:"created_from_ip_address,omitempty"`
	ModifiedAt            time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedFromIPAddress string             `bson:"modified_from_ip_address" json:"modified_from_ip_address,omitempty"`

	// ProjectID variable used to track ownership of this pin object. Used primarily for organization purposes.
	ProjectID primitive.ObjectID `bson:"project_id" json:"project_id,omitempty"`
}

type PinObjectAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// PinObjectStorer Interface for pinobject.
type PinObjectStorer interface {
	Create(ctx context.Context, m *PinObject) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*PinObject, error)
	GetByRequestID(ctx context.Context, requestID primitive.ObjectID) (*PinObject, error)
	UpdateByID(ctx context.Context, m *PinObject) error
	UpdateByRequestID(ctx context.Context, m *PinObject) error
	ListByFilter(ctx context.Context, m *PinObjectPaginationListFilter) (*PinObjectPaginationListResult, error)
	ListByProjectID(ctx context.Context, projectID primitive.ObjectID) (*PinObjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *PinObjectPaginationListFilter) ([]*PinObjectAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	DeleteByCID(ctx context.Context, cid primitive.ObjectID) error
	DeleteByRequestID(ctx context.Context, requestID primitive.ObjectID) error
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
		loggerp.Error("failed deleting all indexes",
			slog.Any("err", err))

		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	// Note:
	// * 1 for ascending
	// * -1 for descending
	// * "text" for text indexes

	// The following few lines of code will create the index for our app for this
	// colleciton.
	_, err := uc.Indexes().CreateMany(context.TODO(), []mongo.IndexModel{
		{Keys: bson.D{{Key: "requestid", Value: 1}}},                                                // Used for get requests
		{Keys: bson.D{{Key: "project_id", Value: 1}, {Key: "created", Value: SortOrderDescending}}}, // Note: Used in default list.

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

	s := &PinObjectStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
