package domain

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	// TenantPendingStatus indicates this tenant needs to be reviewed by CPS_NFTSTORE and approved / rejected.
	TenantPendingStatus                = 1
	TenantActiveStatus                 = 2
	TenantRejectedStatus               = 3
	TenantErrorStatus                  = 4
	TenantArchivedStatus               = 5
	RootType                           = 1
	RetailerType                       = 2
	EstimatedSubmissionsPerMonth1To4   = 1
	EstimatedSubmissionsPerMonth5To10  = 2
	EstimatedSubmissionsPerMonth11To20 = 3
	EstimatedSubmissionsPerMonth20To49 = 4
	EstimatedSubmissionsPerMonth50Plus = 5
	HasOtherGradingServiceYes          = 1
	HasOtherGradingServiceNo           = 2
	RequestWelcomePackageYes           = 1
	RequestWaelcomePackageNo           = 2
	SpecialCollection040001            = 1
)

type Tenant struct {
	ID         primitive.ObjectID `bson:"_id" json:"id"`
	Name       string             `bson:"name" json:"name"` // Created by system.
	ChainID    uint16             `bson:"chain_id" json:"chain_id"`
	Status     int8               `bson:"status" json:"status"`
	ModifiedAt time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	CreatedAt  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	Account    *Account           `bson:"account" json:"account"`
}

type TenantComment struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	TenantID         primitive.ObjectID `bson:"tenant_id" json:"tenant_id"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID  primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByName    string             `bson:"created_by_name" json:"created_by_name"`
	ModifiedAt       time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByName   string             `bson:"modified_by_name" json:"modified_by_name"`
	Content          string             `bson:"content" json:"content"`
}

type TenantListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	TenantID         primitive.ObjectID
	CreatedByUserID  primitive.ObjectID
	ModifiedByUserID primitive.ObjectID
	Status           int8
	ExcludeArchived  bool
	SearchText       string
	CreatedAtGTE     time.Time
}

type TenantListResult struct {
	Results     []*Tenant          `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type TenantAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// TenantRepository Interface for tenant.
type TenantRepository interface {
	Create(ctx context.Context, m *Tenant) error
	GetByName(ctx context.Context, name string) (*Tenant, error)
	GetByID(ctx context.Context, id primitive.ObjectID) (*Tenant, error)
	UpdateByID(ctx context.Context, m *Tenant) error
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// ListByFilter(ctx context.Context, m *TenantPaginationListFilter) (*TenantPaginationListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *TenantPaginationListFilter) ([]*TenantAsSelectOption, error)
}

// type TenantRepositoryImpl struct {
// 	Logger     *slog.Logger
// 	DbClient   *mongo.Client
// 	Collection *mongo.Collection
// }
//
// func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) TenantRepository {
// 	// ctx := context.Background()
// 	uc := client.Database(appCfg.DB.Name).Collection("tenants")
//
// 	// The following few lines of code will create the index for our app for this
// 	// colleciton.
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{
// 			{"name", "text"},
// 		},
// 	}
// 	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
// 	if err != nil {
// 		// It is important that we crash the app on startup to meet the
// 		// requirements of `google/wire` framework.
// 		log.Fatal(err)
// 	}
//
// 	s := &TenantRepositoryImpl{
// 		Logger:     loggerp,
// 		DbClient:   client,
// 		Collection: uc,
// 	}
// 	return s
// }
