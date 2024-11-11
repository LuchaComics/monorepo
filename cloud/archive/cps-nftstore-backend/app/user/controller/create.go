package controller

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-nftstore-backend/utils/httperror"
)

type UserCreateRequestIDO struct {
	TenantID                                         primitive.ObjectID `bson:"tenant_id" json:"tenant_id,omitempty"`
	FirstName                                       string             `json:"first_name"`
	LastName                                        string             `json:"last_name"`
	Email                                           string             `json:"email"`
	Phone                                           string             `json:"phone,omitempty"`
	Country                                         string             `json:"country,omitempty"`
	Region                                          string             `json:"region,omitempty"`
	City                                            string             `json:"city,omitempty"`
	PostalCode                                      string             `json:"postal_code,omitempty"`
	AddressLine1                                    string             `json:"address_line1,omitempty"`
	AddressLine2                                    string             `json:"address_line2,omitempty"`
	HowDidYouHearAboutUs                            int8               `json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther                       string             `json:"how_did_you_hear_about_us_other,omitempty"`
	HowLongTenantOperating                           int8               `bson:"how_long_tenant_operating" json:"how_long_tenant_operating,omitempty"`
	GradingComicsExperience                         string             `bson:"grading_comics_experience" json:"grading_comics_experience,omitempty"`
	RetailPartnershipReason                         string             `bson:"retail_partnership_reason" json:"retail_partnership_reason,omitempty"` // "Please describe how you could become a good retail partner for CPS_NFTSTORE"
	CPS_NFTSTOREPartnershipReason                            string             `bson:"cps-nftstore_partnership_reason" json:"cps-nftstore_partnership_reason,omitempty"`       // "Please describe how CPS_NFTSTORE could help you grow your business"
	AgreeTOS                                        bool               `json:"agree_tos,omitempty"`
	AgreePromotionsEmail                            bool               `json:"agree_promotions_email,omitempty"`
	Status                                          int8               `bson:"status" json:"status"`
	Role                                            int8               `bson:"role" json:"role"`
	HasShippingAddress                              bool               `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName                                    string             `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone                                   string             `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry                                 string             `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion                                  string             `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity                                    string             `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode                              string             `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1                            string             `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2                            string             `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	HowLongCollectingComicBooksForGrading           int8               `bson:"how_long_collecting_comic_books_for_grading" json:"how_long_collecting_comic_books_for_grading"`
	HasPreviouslySubmittedComicBookForGrading       int8               `bson:"has_previously_submitted_comic_book_for_grading" json:"has_previously_submitted_comic_book_for_grading"`
	HasOwnedGradedComicBooks                        int8               `bson:"has_owned_graded_comic_books" json:"has_owned_graded_comic_books"`
	HasRegularComicBookShop                         int8               `bson:"has_regular_comic_book_shop" json:"has_regular_comic_book_shop"`
	HasPreviouslyPurchasedFromAuctionSite           int8               `bson:"has_previously_purchased_from_auction_site" json:"has_previously_purchased_from_auction_site"`
	HasPreviouslyPurchasedFromFacebookMarketplace   int8               `bson:"has_previously_purchased_from_facebook_marketplace" json:"has_previously_purchased_from_facebook_marketplace"`
	HasRegularlyAttendedComicConsOrCollectibleShows int8               `bson:"has_regularly_attended_comic_cons_or_collectible_shows" json:"has_regularly_attended_comic_cons_or_collectible_shows"`
}

func (impl *UserControllerImpl) userFromCreateRequest(requestData *UserCreateRequestIDO) (*user_s.User, error) {
	return &user_s.User{
		TenantID:                               requestData.TenantID,
		FirstName:                             requestData.FirstName,
		LastName:                              requestData.LastName,
		Email:                                 requestData.Email,
		Phone:                                 requestData.Phone,
		Country:                               requestData.Country,
		Region:                                requestData.Region,
		City:                                  requestData.City,
		PostalCode:                            requestData.PostalCode,
		AddressLine1:                          requestData.AddressLine1,
		AddressLine2:                          requestData.AddressLine2,
		HowDidYouHearAboutUs:                  requestData.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther:             requestData.HowDidYouHearAboutUsOther,
		AgreeTOS:                              requestData.AgreeTOS,
		AgreePromotionsEmail:                  requestData.AgreePromotionsEmail,
		Status:                                requestData.Status,
		Role:                                  requestData.Role,
		HasShippingAddress:                    requestData.HasShippingAddress,
		ShippingName:                          requestData.ShippingName,
		ShippingPhone:                         requestData.ShippingPhone,
		ShippingCountry:                       requestData.ShippingCountry,
		ShippingRegion:                        requestData.ShippingRegion,
		ShippingCity:                          requestData.ShippingCity,
		ShippingPostalCode:                    requestData.ShippingPostalCode,
		ShippingAddressLine1:                  requestData.ShippingAddressLine1,
		ShippingAddressLine2:                  requestData.ShippingAddressLine2,
		HowLongCollectingComicBooksForGrading: requestData.HowLongCollectingComicBooksForGrading,
		HasPreviouslySubmittedComicBookForGrading:       requestData.HasPreviouslySubmittedComicBookForGrading,
		HasOwnedGradedComicBooks:                        requestData.HasOwnedGradedComicBooks,
		HasRegularComicBookShop:                         requestData.HasRegularComicBookShop,
		HasPreviouslyPurchasedFromAuctionSite:           requestData.HasPreviouslyPurchasedFromAuctionSite,
		HasPreviouslyPurchasedFromFacebookMarketplace:   requestData.HasPreviouslyPurchasedFromFacebookMarketplace,
		HasRegularlyAttendedComicConsOrCollectibleShows: requestData.HasRegularlyAttendedComicConsOrCollectibleShows,
	}, nil
}

func (impl *UserControllerImpl) Create(ctx context.Context, requestData *UserCreateRequestIDO) (*user_s.User, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return nil, err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		m, err := impl.userFromCreateRequest(requestData)
		if err != nil {
			return nil, err
		}

		// Extract from our session the following data.
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)

		// Apply filtering based on ownership and role.
		if userRole != user_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByEmail(sessCtx, m.Email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u != nil {
			impl.Logger.Warn("user already exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
		}

		// Lookup the store in our database, else return a `400 Bad Request` error.
		o, err := impl.TenantStorer.GetByID(sessCtx, m.TenantID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if o == nil {
			impl.Logger.Warn("store does not exist exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("tenant_id", "store does not exist")
		}

		// Modify the user based on role.

		// Add defaults.
		m.Email = strings.ToLower(m.Email)
		m.TenantID = o.ID
		m.TenantName = o.Name
		m.ID = primitive.NewObjectID()
		m.CreatedAt = time.Now()
		m.CreatedByUserID = userID
		m.CreatedByName = userName
		m.ModifiedAt = time.Now()
		m.ModifiedByUserID = userID
		m.ModifiedByName = userName
		m.Name = fmt.Sprintf("%s %s", m.FirstName, m.LastName)
		m.LexicalName = fmt.Sprintf("%s, %s", m.LastName, m.FirstName)
		m.WasEmailVerified = true

		// Generate a temporary password.
		temporaryPassword := primitive.NewObjectID().Hex()

		// Hash our password with the temporary password and attach to account.
		temporaryPasswordHash, err := impl.Password.GenerateHashFromPassword(temporaryPassword)
		if err != nil {
			impl.Logger.Error("hashing error", slog.Any("error", err))
			return nil, err
		}
		m.PasswordHashAlgorithm = impl.Password.AlgorithmName()
		m.PasswordHash = temporaryPasswordHash

		// Save to our database.
		if err := impl.UserStorer.Create(sessCtx, m); err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}

		// Send email to user of the new password.
		if err := impl.TemplatedEmailer.SendNewUserTemporaryPasswordEmail(m.Email, m.FirstName, temporaryPassword); err != nil {
			impl.Logger.Error("failed sending verification email with error", slog.Any("err", err))
			return nil, err
		}

		return m, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	return res.(*user_s.User), nil
}
