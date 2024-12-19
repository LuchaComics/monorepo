package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	AttachmentStatusActive            = 1
	AttachmentStatusError             = 2
	AttachmentStatusArchived          = 3
	AttachmentContentTypeFile         = 1
	AttachmentContentTypeImage        = 2
	AttachmentBelongsToTypeUnassigned = 1
	AttachmentBelongsToTypeSubmission = 2
)

type Attachment struct {
	ID                        primitive.ObjectID `bson:"_id" json:"id"`
	CreatedAt                 time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName         string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID           primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedFromIPAddress      string             `bson:"created_from_ip_address" json:"created_from_ip_address"`
	ModifiedAt                time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName        string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID          primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedFromIPAddress     string             `bson:"modified_from_ip_address" json:"modified_from_ip_address"`
	Name                      string             `bson:"name" json:"name"`
	Description               string             `bson:"description" json:"description"`
	Filename                  string             `bson:"filename" json:"filename"`
	Sha256Hash                string             `bson:"sha_256_hash" json:"sha_256_hash"`
	ObjectKey                 string             `bson:"object_key" json:"object_key"`
	ObjectURL                 string             `bson:"-" json:"object_url"`
	Status                    int8               `bson:"status" json:"status"`
	ContentType               string             `bson:"content_type" json:"content_type"`
	UserID                    primitive.ObjectID `bson:"user_id" json:"user_id"`
	TenantID                  primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	BelongsToUniqueIdentifier primitive.ObjectID `bson:"belongs_to_unique_identifier" json:"belongs_to_unique_identifier"`
	BelongsToType             int8               `bson:"belongs_to_type" json:"belongs_to_type"`
}

type AttachmentFilter struct {
	TenantID       primitive.ObjectID `json:"tenant_id"` // Required for data partitioning
	Name           *string            `json:"name,omitempty"`
	Status         int8               `json:"status,omitempty"`
	UserID         primitive.ObjectID `json:"user_id,omitempty"`
	CreatedAtStart *time.Time         `json:"created_at_start,omitempty"`
	CreatedAtEnd   *time.Time         `json:"created_at_end,omitempty"`
	BelongsToType  int8               `bson:"belongs_type" json:"belongs_type"`

	// Cursor-based pagination
	LastID        *primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt *time.Time          `json:"last_created_at,omitempty"`
	Limit         int64               `json:"limit"`
}

type AttachmentFilterResult struct {
	Attachments   []*Attachment      `json:"attachments"`
	HasMore       bool               `json:"has_more"`
	LastID        primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt time.Time          `json:"last_created_at,omitempty"`
}

// AttachmentRepository Interface for a file that has content which lives in the cloud.
type AttachmentRepository interface {
	Create(ctx context.Context, m *Attachment) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Attachment, error)
	UpdateByID(ctx context.Context, m *Attachment) error
	CountByFilter(ctx context.Context, filter *AttachmentFilter) (uint64, error)
	ListByFilter(ctx context.Context, filter *AttachmentFilter) (*AttachmentFilterResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}
