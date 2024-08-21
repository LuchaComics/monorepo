package datastore

import (
	"time"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-pinws-backend/app/user/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const (
	UserStatusActive   = 1
	UserInactiveStatus = 2
)

type RegisterBusinessRequestIDO struct {
	FirstName                    string `json:"first_name"`
	LastName                     string `json:"last_name"`
	Email                        string `json:"email"`
	Password                     string `json:"password"`
	PasswordRepeated             string `json:"password_repeated"`
	ComicBookTenantName           string `json:"comic_book_tenant_name,omitempty"`
	Phone                        string `json:"phone,omitempty"`
	Country                      string `json:"country,omitempty"`
	Region                       string `json:"region,omitempty"`
	City                         string `json:"city,omitempty"`
	PostalCode                   string `json:"postal_code,omitempty"`
	AddressLine1                 string `json:"address_line1,omitempty"`
	AddressLine2                 string `json:"address_line2,omitempty"`
	StoreLogo                    string `json:"tenant_logo,omitempty"`
	HowDidYouHearAboutUs         int8   `json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther    string `json:"how_did_you_hear_about_us_other,omitempty"`
	HowLongTenantOperating        int8   `json:"how_long_tenant_operating,omitempty"`
	GradingComicsExperience      string `json:"grading_comics_experience,omitempty"`
	RetailPartnershipReason      string `bson:"retail_partnership_reason" json:"retail_partnership_reason,omitempty"` // "Please describe how you could become a good retail partner for CPS_PINWS"
	CPS_PINWSPartnershipReason         string `bson:"cps-pinws_partnership_reason" json:"cps-pinws_partnership_reason,omitempty"`       // "Please describe how CPS_PINWS could help you grow your business"
	AgreeTOS                     bool   `json:"agree_tos,omitempty"`
	AgreePromotionsEmail         bool   `json:"agree_promotions_email,omitempty"`
	HasShippingAddress           bool   `json:"has_shipping_address,omitempty"`
	ShippingName                 string `json:"shipping_name,omitempty"`
	ShippingPhone                string `json:"shipping_phone,omitempty"`
	ShippingCountry              string `json:"shipping_country,omitempty"`
	ShippingRegion               string `json:"shipping_region,omitempty"`
	ShippingCity                 string `json:"shipping_city,omitempty"`
	ShippingPostalCode           string `json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1         string `json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2         string `json:"shipping_address_line2,omitempty"`
	WebsiteURL                   string `bson:"website_url" json:"website_url"`
	EstimatedSubmissionsPerMonth int8   `bson:"estimated_submissions_per_month" json:"estimated_submissions_per_month"`
	HasOtherGradingService       int8   `bson:"has_other_grading_service" json:"has_other_grading_service"`
	OtherGradingServiceName      string `bson:"other_grading_service_name" json:"other_grading_service_name"`
	RequestWelcomePackage        int8   `bson:"request_welcome_package" json:"request_welcome_package"`
	Timezone                     string `bson:"timezone" json:"timezone"`
}

type RegisterBusinessResponseIDO struct {
	User                   *user_s.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

type RegisterCustomerRequestIDO struct {
	FirstName                                       string             `json:"first_name"`
	LastName                                        string             `json:"last_name"`
	Email                                           string             `json:"email"`
	Password                                        string             `json:"password"`
	PasswordRepeated                                string             `json:"password_repeated"`
	Phone                                           string             `json:"phone,omitempty"`
	Country                                         string             `json:"country,omitempty"`
	Region                                          string             `json:"region,omitempty"`
	City                                            string             `json:"city,omitempty"`
	PostalCode                                      string             `json:"postal_code,omitempty"`
	AddressLine1                                    string             `json:"address_line1,omitempty"`
	AddressLine2                                    string             `json:"address_line2,omitempty"`
	AgreeTOS                                        bool               `json:"agree_tos,omitempty"`
	AgreePromotionsEmail                            bool               `json:"agree_promotions_email,omitempty"`
	HasShippingAddress                              bool               `json:"has_shipping_address,omitempty"`
	ShippingName                                    string             `json:"shipping_name,omitempty"`
	ShippingPhone                                   string             `json:"shipping_phone,omitempty"`
	ShippingCountry                                 string             `json:"shipping_country,omitempty"`
	ShippingRegion                                  string             `json:"shipping_region,omitempty"`
	ShippingCity                                    string             `json:"shipping_city,omitempty"`
	ShippingPostalCode                              string             `json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1                            string             `json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2                            string             `json:"shipping_address_line2,omitempty"`
	TenantID                                         primitive.ObjectID `json:"tenant_id"`
	HowDidYouHearAboutUs                            int8               `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther                       string             `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	HowLongCollectingComicBooksForGrading           int8               `bson:"how_long_collecting_comic_books_for_grading" json:"how_long_collecting_comic_books_for_grading"`
	HasPreviouslySubmittedComicBookForGrading       int8               `bson:"has_previously_submitted_comic_book_for_grading" json:"has_previously_submitted_comic_book_for_grading"`
	HasOwnedGradedComicBooks                        int8               `bson:"has_owned_graded_comic_books" json:"has_owned_graded_comic_books"`
	HasRegularComicBookShop                         int8               `bson:"has_regular_comic_book_shop" json:"has_regular_comic_book_shop"`
	HasPreviouslyPurchasedFromAuctionSite           int8               `bson:"has_previously_purchased_from_auction_site" json:"has_previously_purchased_from_auction_site"`
	HasPreviouslyPurchasedFromFacebookMarketplace   int8               `bson:"has_previously_purchased_from_facebook_marketplace" json:"has_previously_purchased_from_facebook_marketplace"`
	HasRegularlyAttendedComicConsOrCollectibleShows int8               `bson:"has_regularly_attended_comic_cons_or_collectible_shows" json:"has_regularly_attended_comic_cons_or_collectible_shows"`
}

type RegisterCustomerResponseIDO struct {
	User                   *user_s.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

type LoginResponseIDO struct {
	User                   *user_s.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

type VerifyResponseIDO struct {
	Message  string `json:"message"`
	UserRole int8   `bson:"user_role" json:"user_role"`
}
