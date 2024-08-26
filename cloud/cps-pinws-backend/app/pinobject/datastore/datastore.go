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

	OwnershipTypeUser       = 1
	OwnershipTypeSubmission = 2
	OwnershipTypeStore      = 3
	ContentTypeFile         = 1
	ContentTypeImage        = 2
)

// PinObject is a representation of a pin request. It means it is the IPFS content which we are saving in our system and sharing to the IPFS network, also know as "pinning". This structure has the core variables required to work with IPFS as per their documentation https://ipfs.github.io/pinning-services-api-spec/#section/Schemas/Identifiers, in additon we also have our applications specific varaibles.
type PinObject struct {
	TenantID       primitive.ObjectID `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	TenantName     string             `bson:"tenant_name" json:"tenant_name"`
	TenantTimezone string             `bson:"tenant_timezone" json:"tenant_timezone"`

	// ID variable is the unique identifier we use internally in our system.
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	Status             string             `bson:"status" json:"status"`

	OwnershipID   primitive.ObjectID `bson:"ownership_id" json:"ownership_id"`
	OwnershipType int8               `bson:"ownership_type" json:"ownership_type"`

	// RequestID variable is the public viewable unique identifier of this pin.
	RequestID primitive.ObjectID `bson:"requestid" json:"requestid"`

	// CID variable is the unique identifier of our content on the IPFS network. The official definition is: Content Identifier (CID) points at the root of a DAG that is pinned recursively.
	CID string `bson:"cid" json:"cid"`

	// Name variable used to provide human readable description for the content. This is optional.
	Name string `bson:"name" json:"name"`

	// Addresses provided in origins list are relevant only during the initial pinning, and don't need to be persisted by the pinning service
	Origins []string `bson:"origins" json:"origins"`

	// Addresses in the delegates array are peers designated by pinning service that will receive the pin data over bitswap
	Delegates []string `bson:"delegates" json:"delegates"`

	// Any additional vendor-specific information is returned in optional info.
	Meta map[string]string `bson:"meta" json:"meta"`

	// Any additional vendor-specific information is returned in optional info.
	Info map[string]string `bson:"info" json:"info"`

	Filename    string `bson:"filename" json:"filename"`
	ObjectKey   string `bson:"object_key" json:"object_key"`
	ObjectURL   string `bson:"object_url" json:"object_url"`
	ContentType int8   `bson:"content_type" json:"content_type"`
}

type PinObjectAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// PinObjectStorer Interface for pinobject.
type PinObjectStorer interface {
	Create(ctx context.Context, m *PinObject) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*PinObject, error)
	UpdateByID(ctx context.Context, m *PinObject) error
	ListByFilter(ctx context.Context, m *PinObjectPaginationListFilter) (*PinObjectPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *PinObjectPaginationListFilter) ([]*PinObjectAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
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

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"tenant_name", "text"},
			{"name", "text"},
			{"description", "text"},
			{"filename", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
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
