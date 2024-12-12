package domain

import (
	"context"
	"log/slog"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const (
	UserStatusActive   = 1
	UserStatusArchived = 100
	UserRoleRoot       = 1
	UserRoleRetailer   = 2
	UserRoleCustomer   = 3
)

type User struct {
	ID                        primitive.ObjectID `bson:"_id" json:"id"`
	TenantID                  primitive.ObjectID `bson:"tenant_id" json:"tenant_id,omitempty"`
	TenantName                string             `bson:"tenant_name" json:"tenant_name"`
	FirstName                 string             `bson:"first_name" json:"first_name"`
	LastName                  string             `bson:"last_name" json:"last_name"`
	Name                      string             `bson:"name" json:"name"`
	LexicalName               string             `bson:"lexical_name" json:"lexical_name"`
	Email                     string             `bson:"email" json:"email"`
	PasswordHashAlgorithm     string             `bson:"password_hash_algorithm" json:"password_hash_algorithm,omitempty"`
	PasswordHash              string             `bson:"password_hash" json:"password_hash,omitempty"`
	Role                      int8               `bson:"role" json:"role"`
	WasEmailVerified          bool               `bson:"was_email_verified" json:"was_email_verified,omitempty"`
	EmailVerificationCode     string             `bson:"email_verification_code,omitempty" json:"email_verification_code,omitempty"`
	EmailVerificationExpiry   time.Time          `bson:"email_verification_expiry,omitempty" json:"email_verification_expiry,omitempty"`
	Phone                     string             `bson:"phone" json:"phone,omitempty"`
	Country                   string             `bson:"country" json:"country,omitempty"`
	Region                    string             `bson:"region" json:"region,omitempty"`
	City                      string             `bson:"city" json:"city,omitempty"`
	PostalCode                string             `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1              string             `bson:"address_line1" json:"address_line1,omitempty"`
	AddressLine2              string             `bson:"address_line2" json:"address_line2,omitempty"`
	HasShippingAddress        bool               `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName              string             `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone             string             `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry           string             `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion            string             `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity              string             `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode        string             `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1      string             `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2      string             `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	TenantLogoS3Key           string             `bson:"tenant_logo_s3_key" json:"tenant_logo_s3_key,omitempty"`
	TenantLogoTitle           string             `bson:"tenant_logo_title" json:"tenant_logo_title,omitempty"`
	TenantLogoFileURL         string             `bson:"tenant_logo_file_url" json:"tenant_logo_file_url,omitempty"`     // (Optional, added by endpoint)
	TenantLogoFileURLExpiry   time.Time          `bson:"tenant_logo_file_url_expiry" json:"tenant_logo_file_url_expiry"` // (Optional, added by endpoint)
	HowDidYouHearAboutUs      int8               `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string             `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreeTermsOfService       bool               `bson:"agree_terms_of_service" json:"agree_terms_of_service,omitempty"`
	AgreePromotions           bool               `bson:"agree_promotions" json:"agree_promotions,omitempty"`
	CreatedByUserID           primitive.ObjectID `bson:"created_by_user_id" json:"created_by_user_id"`
	CreatedAt                 time.Time          `bson:"created_at" json:"created_at,omitempty"`
	CreatedByName             string             `bson:"created_by_name" json:"created_by_name"`
	ModifiedByUserID          primitive.ObjectID `bson:"modified_by_user_id" json:"modified_by_user_id"`
	ModifiedAt                time.Time          `bson:"modified_at" json:"modified_at,omitempty"`
	ModifiedByName            string             `bson:"modified_by_name" json:"modified_by_name"`
	Status                    int8               `bson:"status" json:"status"`
	Comments                  []*UserComment     `bson:"comments" json:"comments"`
	IsStarred                 bool               `bson:"is_starred" json:"is_starred,omitempty"`
	// The name of the payment processor we are using to handle payments with
	// this particular member.
	PaymentProcessorName string `bson:"payment_processor_name" json:"payment_processor_name"`
	// The unique identifier used by the payment processor which has a somesort of
	// copy of this member's details saved and we can reference that customer on
	// the payment processor using this `customer_id`.
	PaymentProcessorCustomerID string `bson:"payment_processor_customer_id" json:"payment_processor_customer_id"`

	// OTPEnabled controls whether we force 2FA or not during login.
	OTPEnabled bool `bson:"otp_enabled" json:"otp_enabled"`

	// OTPVerified indicates user has successfully validated their opt token afer enabling 2FA thus turning it on.
	OTPVerified bool `bson:"otp_verified" json:"otp_verified"`

	// OTPValidated automatically gets set as `false` on successful login and then sets `true` once successfully validated by 2FA.
	OTPValidated bool `bson:"otp_validated" json:"otp_validated"`

	// OTPSecret the unique one-time password secret to be shared between our
	// backend and 2FA authenticator sort of apps that support `TOPT`.
	OTPSecret string `bson:"otp_secret" json:"-"`

	// OTPAuthURL is the URL used to share.
	OTPAuthURL string `bson:"otp_auth_url" json:"-"`

	// OTPBackupCodeHash is the one-time use backup code which resets the 2FA settings and allow the user to setup 2FA from scratch for the user.
	OTPBackupCodeHash string `bson:"otp_backup_code_hash" json:"-"`

	// OTPBackupCodeHashAlgorithm tracks the hashing algorithm used.
	OTPBackupCodeHashAlgorithm string `bson:"otp_backup_code_hash_algorithm" json:"-"`

	HowLongCollectingComicBooksForGrading           int8 `bson:"how_long_collecting_comic_books_for_grading" json:"how_long_collecting_comic_books_for_grading"`
	HasPreviouslySubmittedComicBookForGrading       int8 `bson:"has_previously_submitted_comic_book_for_grading" json:"has_previously_submitted_comic_book_for_grading"`
	HasOwnedGradedComicBooks                        int8 `bson:"has_owned_graded_comic_books" json:"has_owned_graded_comic_books"`
	HasRegularComicBookShop                         int8 `bson:"has_regular_comic_book_shop" json:"has_regular_comic_book_shop"`
	HasPreviouslyPurchasedFromAuctionSite           int8 `bson:"has_previously_purchased_from_auction_site" json:"has_previously_purchased_from_auction_site"`
	HasPreviouslyPurchasedFromFacebookMarketplace   int8 `bson:"has_previously_purchased_from_facebook_marketplace" json:"has_previously_purchased_from_facebook_marketplace"`
	HasRegularlyAttendedComicConsOrCollectibleShows int8 `bson:"has_regularly_attended_comic_cons_or_collectible_shows" json:"has_regularly_attended_comic_cons_or_collectible_shows"`

	// WalletAddress variable holds the address of the user's wallet
	// which is used by this faucet application to send.
	WalletAddress *common.Address `bson:"wallet_address" json:"wallet_address"`

	// LastCoinsDepositAt variable keeps track of when this faucet sent coins
	// to this user's account.
	LastCoinsDepositAt time.Time `bson:"last_coins_deposit_at" json:"last_coins_deposit_at"`
}

type UserComment struct {
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

type UserListFilter struct {
	// Pagination related.
	Cursor    primitive.ObjectID
	PageSize  int64
	SortField string
	SortOrder int8 // 1=ascending | -1=descending

	// Filter related.
	TenantID        primitive.ObjectID
	Role            int8
	Status          int8
	UUIDs           []string
	ExcludeArchived bool
	SearchText      string
	FirstName       string
	LastName        string
	Email           string
	Phone           string
	CreatedAtGTE    time.Time
	IsStarred       int8 //0=All, 1=True, 2=False
}

type UserListResult struct {
	Results     []*User            `json:"results"`
	NextCursor  primitive.ObjectID `json:"next_cursor"`
	HasNextPage bool               `json:"has_next_page"`
}

type UserAsSelectOption struct {
	Value primitive.ObjectID `bson:"_id" json:"value"` // Extract from the database `_id` field and output through API as `value`.
	Label string             `bson:"name" json:"label"`
}

// UserPaginationListResult represents the paginated list results for
// the associate records.
type UserPaginationListResult struct {
	Results     []*User `json:"results"`
	NextCursor  string  `json:"next_cursor"`
	HasNextPage bool    `json:"has_next_page"`
}

// UserStorer Interface for user.
type UserRepository interface {
	Create(ctx context.Context, m *User) error
	GetByID(ctx context.Context, id primitive.ObjectID) (*User, error)
	GetByEmail(ctx context.Context, email string) (*User, error)
	GetByVerificationCode(ctx context.Context, verificationCode string) (*User, error)
	DeleteByID(ctx context.Context, id primitive.ObjectID) error
	CheckIfExistsByEmail(ctx context.Context, email string) (bool, error)
	UpdateByID(ctx context.Context, m *User) error
	// ListByFilter(ctx context.Context, f *UserPaginationListFilter) (*UserPaginationListResult, error)
	// ListAsSelectOptionByFilter(ctx context.Context, f *UserPaginationListFilter) ([]*UserAsSelectOption, error)
	// ListAllRootStaff(ctx context.Context) (*UserPaginationListResult, error)
	// ListAllRetailerStaffForTenantID(ctx context.Context, storeID primitive.ObjectID) (*UserPaginationListResult, error)
	// //TODO: Add more...
}

type UserStorerImpl struct {
	Logger     *slog.Logger
	DbClient   *mongo.Client
	Collection *mongo.Collection
}

// func NewDatastore(appCfg *c.Configuration, loggerp *slog.Logger, client *mongo.Client) UserStorer {
// 	// ctx := context.Background()
// 	uc := client.Database(appCfg.DB.Name).Collection("users")
//
// 	// The following few lines of code will create the index for our app for this
// 	// colleciton.
// 	indexModel := mongo.IndexModel{
// 		Keys: bson.D{
// 			{"tenant_name", "text"},
// 			{"name", "text"},
// 			{"lexical_name", "text"},
// 			{"email", "text"},
// 			{"phone", "text"},
// 			{"country", "text"},
// 			{"region", "text"},
// 			{"city", "text"},
// 			{"postal_code", "text"},
// 			{"address_line1", "text"},
// 		},
// 	}
// 	_, err := uc.Indexes().CreateOne(context.TODO(), indexModel)
// 	if err != nil {
// 		// It is important that we crash the app on startup to meet the
// 		// requirements of `google/wire` framework.
// 		log.Fatal(err)
// 	}
//
// 	s := &UserStorerImpl{
// 		Logger:     loggerp,
// 		DbClient:   client,
// 		Collection: uc,
// 	}
// 	return s
// }
