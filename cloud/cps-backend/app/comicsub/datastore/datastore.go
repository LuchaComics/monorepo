package datastore

import (
	"context"
	"log"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	c "github.com/LuchaComics/monorepo/cloud/cps-backend/config"
)

const (
	StatusPaymentRequired                            = 1
	StatusWaiting                                    = 1
	StatusReceived                                   = 2
	StatusPending                                    = 3
	StatusInProcess                                  = 4
	StatusComplete                                   = 5
	StatusShipped                                    = 6
	StatusCompletedByRetailPartner                   = 7
	StatusError                                      = 10
	StatusArchived                                   = 11
	ServiceTypePreScreening                          = 1
	ServiceTypePedigree                              = 2
	ServiceTypeCPSCapsule                            = 3
	ServiceTypeCPSCapsuleIndieMintGem                = 4
	ServiceTypeCPSCapsuleSignatureCollection         = 5
	ServiceTypeCPSCapsuleYouGrade                    = 6
	ServiceTypeCPSCapsuleYouGradeSignatureCollection = 7
	FindingPoor                                      = 1
	FindingFair                                      = 2
	FindingGood                                      = 3
	FindingVeryGood                                  = 4
	FindingFine                                      = 5
	FindingVeryFine                                  = 6
	FindingNearMint                                  = 7
	YesItShowsSignsOfTamperingOrRestoration          = 1
	NoItDoesNotShowsSignsOfTamperingOrRestoration    = 2
	GradingScaleLetter                               = 1
	GradingScaleNumber                               = 2
	GradingScaleCPSPercentage                        = 3
	CollectibleTypeGeneric                           = 1
	PrimaryLabelDetailsRegularEdition                = 2
	PrimaryLabelDetailsDirectEdition                 = 3
	PrimaryLabelDetailsNewsstandEdition              = 4
	PrimaryLabelDetailsVariantCover                  = 5
	PrimaryLabelDetailsCanadianPriceVariant          = 6
	PrimaryLabelDetailsFacsimile                     = 7
	PrimaryLabelDetailsReprint                       = 8
	PrimaryLabelDetailsOther                         = 1
	PaymentProcessorStripe                           = 1
)

type ComicSubmission struct {
	ID                                 primitive.ObjectID          `bson:"_id" json:"id"`
	StoreID                            primitive.ObjectID          `bson:"store_id,omitempty" json:"store_id,omitempty"`
	StoreName                          string                      `bson:"store_name" json:"store_name"`
	StoreSpecialCollection             int8                        `bson:"store_special_colleciton" json:"store_special_colleciton"`
	StoreTimezone                      string                      `bson:"store_timezone" json:"store_timezone"`
	CPSRNClassification                string                      `bson:"cpsrn_classification" json:"cpsrn_classification"`
	CPSRN                              string                      `bson:"cpsrn" json:"cpsrn"`
	CreatedAt                          time.Time                   `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID                    primitive.ObjectID          `bson:"created_by_user_id,omitempty" json:"created_by_user_id,omitempty"`
	CreatedByUserRole                  int8                        `bson:"created_by_user_role" json:"created_by_user_role"`
	ModifiedAt                         time.Time                   `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID                   primitive.ObjectID          `bson:"modified_by_user_id,omitempty" json:"modified_by_user_id,omitempty"`
	ModifiedByUserRole                 int8                        `bson:"modified_by_user_role" json:"modified_by_user_role"`
	ServiceType                        int8                        `bson:"service_type" json:"service_type"`
	Status                             int8                        `bson:"status" json:"status"`
	SubmissionDate                     time.Time                   `bson:"submission_date" json:"submission_date"`
	Item                               string                      `bson:"item" json:"item"` // Created by system.
	SeriesTitle                        string                      `bson:"series_title" json:"series_title"`
	IssueVol                           string                      `bson:"issue_vol" json:"issue_vol"`
	IssueNo                            string                      `bson:"issue_no" json:"issue_no"`
	IssueCoverYear                     int64                       `bson:"issue_cover_year" json:"issue_cover_year"`
	IssueCoverMonth                    int8                        `bson:"issue_cover_month" json:"issue_cover_month"`
	PublisherName                      int8                        `bson:"publisher_name" json:"publisher_name"`
	PublisherNameOther                 string                      `bson:"publisher_name_other" json:"publisher_name_other"`
	IsKeyIssue                         bool                        `bson:"is_key_issue" json:"is_key_issue"`
	KeyIssue                           int8                        `bson:"key_issue" json:"key_issue"`
	KeyIssueOther                      string                      `bson:"key_issue_other" json:"key_issue_other"`
	KeyIssueDetail                     string                      `bson:"key_issue_detail" json:"key_issue_detail"`
	IsInternationalEdition             bool                        `bson:"is_international_edition" json:"is_international_edition"`
	IsVariantCover                     bool                        `bson:"is_variant_cover" json:"is_variant_cover"`
	VariantCoverDetail                 string                      `bson:"variant_cover_detail" json:"variant_cover_detail"`
	Printing                           int8                        `bson:"printing" json:"printing"`
	PrimaryLabelDetails                int8                        `bson:"primary_label_details" json:"primary_label_details"`
	PrimaryLabelDetailsOther           string                      `bson:"primary_label_details_other" json:"primary_label_details_other"`
	SpecialNotes                       string                      `bson:"special_notes" json:"special_notes"`
	GradingNotes                       string                      `bson:"grading_notes" json:"grading_notes"`
	CreasesFinding                     string                      `bson:"creases_finding" json:"creases_finding"`
	TearsFinding                       string                      `bson:"tears_finding" json:"tears_finding"`
	MissingPartsFinding                string                      `bson:"missing_parts_finding" json:"missing_parts_finding"`
	StainsFinding                      string                      `bson:"stains_finding" json:"stains_finding"`
	DistortionFinding                  string                      `bson:"distortion_finding" json:"distortion_finding"`
	PaperQualityFinding                string                      `bson:"paper_quality_finding" json:"paper_quality_finding"`
	SpineFinding                       string                      `bson:"spine_finding" json:"spine_finding"`
	CoverFinding                       string                      `bson:"cover_finding" json:"cover_finding"`
	ShowsSignsOfTamperingOrRestoration int8                        `bson:"shows_signs_of_tampering_or_restoration" json:"shows_signs_of_tampering_or_restoration"`
	GradingScale                       int8                        `bson:"grading_scale" json:"grading_scale"`
	OverallLetterGrade                 string                      `bson:"overall_letter_grade" json:"overall_letter_grade"`
	OverallNumberGrade                 float64                     `bson:"overall_number_grade" json:"overall_number_grade"`
	CpsPercentageGrade                 float64                     `bson:"cps_percentage_grade" json:"cps_percentage_grade"`
	IsOverallLetterGradeNearMintPlus   bool                        `bson:"is_overall_letter_grade_near_mint_plus" json:"is_overall_letter_grade_near_mint_plus"`
	InspectorID                        primitive.ObjectID          `bson:"inspector_id,omitempty" json:"inspector_id,omitempty"` // This is the customer this submission belongs to.
	InspectorFirstName                 string                      `bson:"inspector_first_name" json:"inspector_first_name"`
	InspectorLastName                  string                      `bson:"inspector_last_name" json:"inspector_last_name"`
	InspectorSignature                 string                      `bson:"inspector_signature" json:"user_signature"`
	CustomerID                         primitive.ObjectID          `bson:"customer_id,omitempty" json:"customer_id,omitempty"` // This is the customer this submission belongs to.
	CustomerFirstName                  string                      `bson:"customer_first_name" json:"customer_first_name"`
	CustomerLastName                   string                      `bson:"customer_last_name" json:"customer_last_name"`
	Comments                           []*ComicSubmissionComment   `bson:"comments" json:"comments,omitempty"`
	CollectibleType                    int8                        `bson:"collectible_type" json:"collectible_type"`
	Signatures                         []*ComicSubmissionSignature `bson:"signatures" json:"signatures,omitempty"`
	// FileAttachments                    []*ComicSubmissionFileAttachment  `bson:"file_attachments" json:"file_attachments,omitempty"`
	// ImageAttachments                   []*ComicSubmissionImageAttachment `bson:"image_attachments" json:"image_attachments,omitempty"`
	FindingsFormObjectKey       string    `bson:"findings_form_object_key" json:"findings_form_object_key"`
	FindingsFormObjectURL       string    `bson:"findings_form_object_url" json:"findings_form_object_url"`
	FindingsFormObjectURLExpiry time.Time `bson:"findings_form_object_url_expiry" json:"findings_form_object_url_expiry"`
	LabelObjectKey              string    `bson:"label_object_key" json:"label_object_key"`
	LabelObjectURL              string    `bson:"label_object_url" json:"label_object_url"`
	LabelObjectURLExpiry        time.Time `bson:"label_object_url_expiry" json:"label_object_url_expiry"`
	// CreditID stores the unique ID from the `Credit` table of the credit used to purchase this comic submission.
	CreditID primitive.ObjectID `bson:"credit_id,omitempty" json:"credit_id,omitempty"`
	// PaymentProcessorName represents the name of the payment processor we used in the purchase.
	PaymentProcessor int8 `bson:"payment_processor" json:"payment_processor"`
	// PaymentProcessorPaymentIntentID represent the unique id returned by the payment processor that this comic book submisison was successfully purchased by the customer. If
	// this value is blank ("") then the user has not made a purchase, or burned a credit instead.
	PaymentProcessorPurchaseID string `bson:"payment_processor_purchase_id" json:"payment_processor_purchase_id"`
	// PaymentProcessorPurchaseStatus stores the status set by the payment processor.
	PaymentProcessorPurchaseStatus string `bson:"payment_processor_purchase_status" json:"payment_processor_purchase_status"`
	// PaymentProcessorPurchasedAt represents the date/time this comic book submission was purchased on.
	PaymentProcessorPurchasedAt   time.Time `bson:"payment_processor_purchased_at" json:"payment_processor_purchased_at"`
	PaymentProcessorPurchaseError string    `bson:"payment_processor_purchase_error" json:"payment_processor_purchase_error"`
	// PaymentProcessorReceiptID is the unique id set by the payment processor for this particular receipt.
	PaymentProcessorReceiptID string `bson:"payment_processor_receipt_id" json:"payment_processor_receipt_id"`
	// PaymentProcessorReceiptURL is the external URL to the payment processors receipt hosted service.
	PaymentProcessorReceiptURL string `bson:"payment_processor_receipt_url" json:"payment_processor_receipt_url"`
	// AmountSubtotal is the pre-tax amount.
	AmountSubtotal float64 `bson:"amount_subtotal" json:"amount_subtotal"`
	// AmountTax is the sum of all the tax amounts.
	AmountTax float64 `bson:"amount_tax" json:"amount_tax"`
	// AmountTotal of total of all items after discounts and taxes are applied.
	AmountTotal float64 `bson:"amount_total" json:"amount_total"`
}

// type ComicSubmissionImageAttachment struct {
// 	StoreID   primitive.ObjectID `bson:"store_id,omitempty" json:"store_id,omitempty"`
// 	StoreName string             `bson:"store_name" json:"store_name"`
// 	SubmissionID     primitive.ObjectID `bson:"submission_id" json:"submission_id"`
// 	ID               primitive.ObjectID `bson:"_id" json:"id"`
// 	Filename         string             `bson:"filename" json:"filename"`
// 	ObjectKey        string             `bson:"object_key" json:"object_key"`
// 	ObjectURL        string             `bson:"object_url" json:"object_url"`
// 	ObjectURLExpiry  time.Time          `bson:"object_url_expiry" json:"object_url_expiry"`
// 	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
// 	CreatedByUserID  primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
// 	CreatedByName    string             `bson:"created_by_name" json:"created_by_name"`
// 	ModifiedAt       time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
// 	ModifiedByUserID primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
// 	ModifiedByName   string             `bson:"modified_by_name" json:"modified_by_name"`
// }
//
// type ComicSubmissionFileAttachment struct {
// 	StoreID     primitive.ObjectID `bson:"store_id,omitempty" json:"store_id,omitempty"`
// 	StoreName   string             `bson:"store_name" json:"store_name"`
// 	SubmissionID       primitive.ObjectID `bson:"submission_id" json:"submission_id"`
// 	ID                 primitive.ObjectID `bson:"_id" json:"id"`
// 	CreatedAt          time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
// 	CreatedByUserName  string             `bson:"created_by_user_name" json:"created_by_user_name"`
// 	CreatedByUserID    primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
// 	ModifiedAt         time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
// 	ModifiedByUserName string             `bson:"modified_by_user_name" json:"modified_by_user_name"`
// 	ModifiedByUserID   primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
// 	Name               string             `bson:"name" json:"name"`
// 	Description        string             `bson:"description" json:"description"`
// 	Filename           string             `bson:"filename" json:"filename"`
// 	ObjectKey          string             `bson:"object_key" json:"object_key"`
// 	ObjectURL          string             `bson:"object_url" json:"object_url"`
// 	ObjectURLExpiry    time.Time          `bson:"object_url_expiry" json:"object_url_expiry"`
// 	Status             int8               `bson:"status" json:"status"`
// }

type ComicSubmissionComment struct {
	ID               primitive.ObjectID `bson:"_id" json:"id"`
	StoreID          primitive.ObjectID `bson:"store_id" json:"store_id"`
	CreatedAt        time.Time          `bson:"created_at,omitempty" json:"created_at,omitempty"`
	CreatedByUserID  primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedByName    string             `bson:"created_by_name" json:"created_by_name"`
	ModifiedAt       time.Time          `bson:"modified_at,omitempty" json:"modified_at,omitempty"`
	ModifiedByUserID primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedByName   string             `bson:"modified_by_name" json:"modified_by_name"`
	Content          string             `bson:"content" json:"content"`
}

type ComicSubmissionSignature struct {
	Role string `bson:"role" json:"role"`
	Name string `bson:"name" json:"name"`
}

type ComicSubmissionAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// ComicSubmissionStorer Interface for submission.
type ComicSubmissionStorer interface {
	Create(ctx context.Context, m *ComicSubmission) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*ComicSubmission, error)
	GetByCPSRN(ctx context.Context, cpsrn string) (*ComicSubmission, error)
	GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*ComicSubmission, error)
	UpdateByID(ctx context.Context, m *ComicSubmission) error
	ListByFilter(ctx context.Context, f *ComicSubmissionPaginationListFilter) (*ComicSubmissionPaginationListResult, error)
	ListAsSelectOptionByFilter(ctx context.Context, f *ComicSubmissionPaginationListFilter) ([]*ComicSubmissionAsSelectOption, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CountAll(ctx context.Context) (int64, error)
	CountByFilter(ctx context.Context, f *ComicSubmissionPaginationListFilter) (int64, error)
	// //TODO: Add more...
}

type ComicSubmissionStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

func NewDatastore(appCfg *c.Conf, loggerp *slog.Logger, client *mongo.Client) ComicSubmissionStorer {
	// ctx := context.Background()
	uc := client.Database(appCfg.DB.Name).Collection("comic_submissions")

	// The following few lines of code will create the index for our app for this
	// colleciton.
	indexModel := mongo.IndexModel{
		Keys: bson.D{
			{"store_name", "text"},
			{"cpsrn", "text"},
			{"item", "text"},
			{"publisher_name_other", "text"},
			{"special_notes", "text"},
			{"grading_notes", "text"},
			{"primary_label_details_other", "text"},
		},
	}
	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
	if err != nil {
		// It is important that we crash the app on startup to meet the
		// requirements of `google/wire` framework.
		log.Fatal(err)
	}

	s := &ComicSubmissionStorerImpl{
		Logger:     loggerp,
		DbClient:   client,
		Collection: uc,
	}
	return s
}
