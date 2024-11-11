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
	RequestWelcomePackageNo            = 2
	SpecialCollection040001            = 1
)

type Tenant struct {
	ID                           primitive.ObjectID `bson:"_id" json:"id"`
	ModifiedAt                   time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserName           string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
	ModifiedByUserID             primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	Type                         int8               `bson:"type" json:"type"`
	Status                       int8               `bson:"status" json:"status"`
	Name                         string             `bson:"name" json:"name"` // Created by system.
	CreatedAt                    time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserName            string             `bson:"created_by_user_name" json:"created_by_user_name"`
	CreatedByUserID              primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	Comments                     []*TenantComment   `bson:"comments" json:"comments"`
	WebsiteURL                   string             `bson:"website_url" json:"website_url"`
	EstimatedSubmissionsPerMonth int8               `bson:"estimated_submissions_per_month" json:"estimated_submissions_per_month"`
	HasOtherGradingService       int8               `bson:"has_other_grading_service" json:"has_other_grading_service"`
	OtherGradingServiceName      string             `bson:"other_grading_service_name" json:"other_grading_service_name"`
	RequestWelcomePackage        int8               `bson:"request_welcome_package" json:"request_welcome_package"`
	HowLongTenantOperating       int8               `bson:"how_long_tenant_operating" json:"how_long_tenant_operating,omitempty"`
	GradingComicsExperience      string             `bson:"grading_comics_experience" json:"grading_comics_experience,omitempty"`
	RetailPartnershipReason      string             `bson:"retail_partnership_reason" json:"retail_partnership_reason,omitempty"`       // "Please describe how you could become a good retail partner for CPS_NFTSTORE"
	CPS_NFTSTOREPartnershipReason   string             `bson:"cps-nftstore_partnership_reason" json:"cps-nftstore_partnership_reason,omitempty"` // "Please describe how CPS_NFTSTORE could help you grow your business"
	PendingReviewEmailSent       bool               `bson:"pending_review_email_sent" json:"pending_review_email_sent,omitempty"`
	Level                        int8               `bson:"level" json:"level,omitempty"`
	// SpecialCollection controls what special numbering system to apply on
	// generating a CSPRN.
	SpecialCollection int8   `bson:"special_collection" json:"special_collection"`
	Timezone          string `bson:"timezone" json:"timezone"` // Created by system.
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

// TenantStorer Interface for tenant.
type TenantStorer interface {
	Create(ctx context.Context, m *Tenant) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Tenant, error)
	UpdateByID(ctx context.Context, m *Tenant) error
	ListByFilter(ctx context.Context, m *TenantPaginationListFilter) (*TenantPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *TenantPaginationListFilter) ([]*TenantAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByID(ctx context.Context, id primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type TenantStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) TenantStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("tenants")

	// The following few lines of code will create the index for our app for this
	// colleciton.
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

	s := &TenantStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
