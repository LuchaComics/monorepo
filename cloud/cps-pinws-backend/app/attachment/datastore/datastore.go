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

const (
	StatusActive            = 1
	StatusError             = 2
	StatusArchived          = 3
	OwnershipTypeUser       = 1
	OwnershipTypeSubmission = 2
	OwnershipTypeStore      = 3
	ContentTypeFile         = 1
	ContentTypeImage        = 2
)

type Attachment struct {
	TenantID           primitive.ObjectID `bson:"tenant_id,omitempty" json:"tenant_id,omitempty"`
	TenantName         string             `bson:"tenant_name" json:"tenant_name"`
	TenantTimezone     string             `bson:"tenant_timezone" json:"tenant_timezone"`
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	Name               string             `bson:"name" json:"name"`
	Description        string             `bson:"description" json:"description"`
	Filename           string             `bson:"filename" json:"filename"`
	ObjectKey          string             `bson:"object_key" json:"object_key"`
	ObjectURL          string             `bson:"object_url" json:"object_url"`
	OwnershipID        primitive.ObjectID `bson:"ownership_id" json:"ownership_id"`
	OwnershipType      int8               `bson:"ownership_type" json:"ownership_type"`
	Status             int8               `bson:"status" json:"status"`
	ContentType        int8               `bson:"content_type" json:"content_type"`
	CID                string             `bson:"cid" json:"cid"` // The unique IPFS CID used to identify this file.
}

type AttachmentAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// AttachmentStorer Interface for attachment.
type AttachmentStorer interface {
	Create(ctx context.Context, m *Attachment) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Attachment, error)
	UpdateByID(ctx context.Context, m *Attachment) error
	ListByFilter(ctx context.Context, m *AttachmentPaginationListFilter) (*AttachmentPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *AttachmentPaginationListFilter) ([]*AttachmentAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	// //TODO: Add more...
}

type AttachmentStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) AttachmentStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("attachments")

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

	s := &AttachmentStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
