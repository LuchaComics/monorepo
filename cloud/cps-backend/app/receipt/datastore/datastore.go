package datastore

import (
	"context"
	"log"
	"log/slog"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
)

const (
	StatusActive           = 1
	StatusArchived         = 2
	PaymentProcessorStripe = 1
)

type Receipt struct {
	StoreID         primitive.ObjectID `bson:"store_id" json:"store_id"`
	StoreName       string             `bson:"store_name" json:"store_name"`
	ID              primitive.ObjectID `bson:"_id" json:"id"`
	UserName        string             `bson:"user_name" json:"user_name"`
	UserLexicalName string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID          primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status          int8               `bson:"status" json:"status"`
	CreatedAt       time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt      time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`

	// The name of the payment processor we used to create this receipt.
	PaymentProcessor int8 `bson:"payment_processor" json:"payment_processor"`

	// PaymentProcessorReceiptID is the unique id set by the payment processor for this particular receipt.
	PaymentProcessorReceiptID string `bson:"payment_processor_receipt_id" json:"payment_processor_receipt_id"`

	// PaymentProcessorReceiptURL is the external URL to the payment processors receipt hosted service.
	PaymentProcessorReceiptURL string `bson:"payment_processor_receipt_url" json:"payment_processor_receipt_url"`

	PaymentProcessorPurchaseID  string    `bson:"payment_processor_purchase_id" json:"payment_processor_purchase_id"`
	PaymentProcessorPurchasedAt time.Time `bson:"payment_processor_purchased_at" json:"payment_processor_purchased_at"`
}

type StripeReceipt struct {
	// The unique identification created by Stripe to present this particular receipt.
	ID string `bson:"id" json:"id"`
	// Time at which the object was created. Measured in seconds since the Unix epoch.
	Created int64 `json:"created"`
	// Whether payment was successfully collected for this receipt. An receipt can be paid (most commonly) with a charge or with credit from the customer's account balance.
	Paid bool `json:"paid"`
	// The URL for the hosted receipt page, which allows customers to view and pay an receipt. If the receipt has not been finalized yet, this will be null.
	HostedReceiptURL string `bson:"hosted_receipt_url" json:"hosted_receipt_url"`
	// The link to download the PDF for the receipt. If the receipt has not been finalized yet, this will be null.
	ReceiptPDF string `json:"receipt_pdf"`
}

type ReceiptListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	StoreID         primitive.ObjectID
	UserID          primitive.ObjectID
	ExcludeArchived bool
	SearchText      string
}

type ReceiptListResult struct {
	Results     []*Receipt         `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// ReceiptStorer Interface for store.
type ReceiptStorer interface {
	Create(ctx context.Context, m *Receipt) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*Receipt, error)
	GetByName(ctx context.Context, name string) (*Receipt, error)
	GetByPaymentProcessorReceiptID(ctx context.Context, paymentProcessorReceiptID string) (*Receipt, error)
	GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*Receipt, error)
	UpdateByID(ctx context.Context, m *Receipt) error
	ListByFilter(ctx context.Context, m *ReceiptPaginationListFilter) (*ReceiptPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *ReceiptPaginationListFilter) ([]*ReceiptAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type ReceiptAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type ReceiptStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) ReceiptStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("receipts")

	// The following few lines of code will create the index for our app for
	// this colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"user_name", "text"},
			{"user_lexical_name", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &ReceiptStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
