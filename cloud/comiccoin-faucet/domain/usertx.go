package domain

import (
	"context"
	"math/big"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserTransactionStatusSubmitted = 1
	UserTransactionStatusRejected  = 2
	UserTransactionStatusAccepted  = 3
)

// UserTransaction structure represents the blockchain transaction that
// belongs to our user in our application.
type UserTransaction struct {
	Transaction
	ID                 primitive.ObjectID `bson:"_id" json:"id"`
	Status             int8               `bson:"status" json:"status"`
	UserID             primitive.ObjectID `bson:"user_id" json:"user_id"` // The user ID that this transaction belongs to.
	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	TenantID           primitive.ObjectID `bson:"tenant_id" json:"tenant_id"` // The faucet this belongs to.
}

type UserTransactionFilter struct {
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

type UserTransactionFilterResult struct {
	UserTransactions []*UserTransaction `json:"user_transactions"`
	HasMore          bool               `json:"has_more"`
	LastID           primitive.ObjectID `json:"last_id,omitempty"`
	LastCreatedAt    time.Time          `json:"last_created_at,omitempty"`
}

// UserTransactionRepository Interface for a UserTransaction model in the database.
type UserTransactionRepository interface {
	Create(ctx context.Context, m *UserTransaction) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*UserTransaction, error)
	GetByNonce(ctx context.Context, nonce *big.Int) (*UserTransaction, error)
	UpdateByID(ctx context.Context, m *UserTransaction) error
	CountByFilter(ctx context.Context, filter *UserTransactionFilter) (uint64, error)
	ListByFilter(ctx context.Context, filter *UserTransactionFilter) (*UserTransactionFilterResult, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
}
