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
	StatusActive   = 1
	StatusArchived = 2
)

type UserPurchase struct {
	StoreID                    primitive.ObjectID `bson:"store_id" json:"store_id"`
	StoreName                  string             `bson:"store_name" json:"store_name"`
	StoreTimezone              string             `bson:"store_timezone" json:"store_timezone"`
	ID                         primitive.ObjectID `bson:"_id" json:"id"`
	UserName                   string             `bson:"user_name" json:"user_name"`
	UserLexicalName            string             `bson:"user_lexical_name" json:"user_lexical_name"`
	UserID                     primitive.ObjectID `bson:"user_id" json:"user_id"`
	Status                     int8               `bson:"status" json:"status"`
	CreatedAt                  time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	ModifiedAt                 time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	OfferID                    primitive.ObjectID `bson:"offer_id" json:"offer_id"`                   // Copied from `Offer`.
	OfferName                  string             `bson:"offer_name" json:"offer_name"`               // Copied from `Offer`.
	OfferDescription           string             `bson:"offer_description" json:"offer_description"` // Copied from `Offer`.
	OfferType                  int8               `bson:"offer_type" json:"offer_type"`
	OfferPrice                 float64            `bson:"offer_price" json:"offer_price"`                         // Copied from `Offer`.
	OfferPriceCurrency         string             `bson:"offer_price_currency" json:"offer_price_currency"`       // Copied from `Offer`.
	OfferPayFrequency          int8               `bson:"offer_pay_frequency" json:"offer_pay_frequency"`         // Copied from `Offer`.ç
	OfferBusinessFunction      int8               `bson:"offer_business_function" json:"offer_business_function"` // Copied from `Offer`.
	OfferServiceType           int8               `bson:"offer_service_type" json:"offer_service_type"`
	ComicSubmissionID          primitive.ObjectID `bson:"comic_submission_id" json:"comic_submission_id"`
	ComicSubmissionSeriesTitle string             `bson:"comic_submission_series_title" json:"comic_submission_series_title"`
	ComicSubmissionIssueVol    string             `bson:"comic_submission_issue_vol" json:"comic_submission_issue_vol"`
	ComicSubmissionIssueNo     string             `bson:"comic_submission_issue_no" json:"comic_submission_issue_no"`
	// The payment processor we used to create this receipt.
	PaymentProcessor int8 `bson:"payment_processor" json:"payment_processor"`
	// PaymentProcessorReceiptID is the unique id set by the payment processor for this particular receipt.
	PaymentProcessorReceiptID string `bson:"payment_processor_receipt_id" json:"payment_processor_receipt_id"`
	// PaymentProcessorReceiptURL is the external URL to the payment processors receipt hosted service.
	PaymentProcessorReceiptURL string `bson:"payment_processor_receipt_url" json:"payment_processor_receipt_url"`
	PaymentProcessorPurchaseID string `bson:"payment_processor_purchase_id" json:"payment_processor_purchase_id"`
	// PaymentProcessorPurchaseStatus stores the status set by the payment processor.
	PaymentProcessorPurchaseStatus string `bson:"payment_processor_purchase_status" json:"payment_processor_purchase_status"`
	// PaymentProcessorPurchasedAt represents the date/time this comic book submission was purchased on.
	PaymentProcessorPurchasedAt   time.Time `bson:"payment_processor_purchased_at" json:"payment_processor_purchased_at"`
	PaymentProcessorPurchaseError string    `bson:"payment_processor_purchase_error" json:"payment_processor_purchase_error"`
	// AmountSubtotal is the pre-tax amount.
	AmountSubtotal float64 `bson:"amount_subtotal" json:"amount_subtotal"`
	// AmountTax is the sum of all the tax amounts.
	AmountTax float64 `bson:"amount_tax" json:"amount_tax"`
	// AmountTotal of total of all items after discounts and taxes are applied.
	AmountTotal float64 `bson:"amount_total" json:"amount_total"`
}

type UserPurchaseListFilter struct {
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

type UserPurchaseListResult struct {
	Results     []*UserPurchase    `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

// UserPurchaseStorer Interface for store.
type UserPurchaseStorer interface {
	Create(ctx context.Context, m *UserPurchase) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*UserPurchase, error)
	GetByName(ctx context.Context, name string) (*UserPurchase, error)
	GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*UserPurchase, error)
	GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*UserPurchase, error)
	UpdateByID(ctx context.Context, m *UserPurchase) error
	ListByFilter(ctx context.Context, m *UserPurchasePaginationListFilter) (*UserPurchasePaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *UserPurchasePaginationListFilter) ([]*UserPurchaseAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByNameInOrgBranch(ctx context.Context, name string, orgID primitive.ObjectID, branchID primitive.ObjectID) (bool, error)
	// //TODO: Add more...
}

type UserPurchaseAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

type UserPurchaseStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) UserPurchaseStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("user_purchases")

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

	s := &UserPurchaseStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
