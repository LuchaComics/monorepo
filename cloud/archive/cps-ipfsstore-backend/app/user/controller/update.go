package controller

import (
	"context"
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-ipfsstore-backend/utils/httperror"
)

type UserUpdateRequestIDO struct {
	ID                                              primitive.ObjectID `bson:"_id" json:"id"`
	TenantID                                        primitive.ObjectID `bson:"tenant_id" json:"tenant_id,omitempty"`
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
	HowDidYouHearAboutUs                            int8               `json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther                       string             `json:"how_did_you_hear_about_us_other,omitempty"`
	HowLongTenantOperating                          int8               `bson:"how_long_tenant_operating" json:"how_long_tenant_operating,omitempty"`
	GradingComicsExperience                         string             `bson:"grading_comics_experience" json:"grading_comics_experience,omitempty"`
	RetailPartnershipReason                         string             `bson:"retail_partnership_reason" json:"retail_partnership_reason,omitempty"`
	CPS_IPFSSTOREPartnershipReason                      string             `bson:"cps-ipfsstore_partnership_reason" json:"cps-ipfsstore_partnership_reason,omitempty"` // "Please describe how CPS_IPFSSTORE could help you grow your business"
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

func (impl *UserControllerImpl) userFromUpdateRequest(requestData *UserUpdateRequestIDO) (*user_s.User, error) {
	passwordHash, err := impl.Password.GenerateHashFromPassword(requestData.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	return &user_s.User{
		ID:                                    requestData.ID,
		TenantID:                              requestData.TenantID,
		FirstName:                             requestData.FirstName,
		LastName:                              requestData.LastName,
		Email:                                 requestData.Email,
		PasswordHash:                          passwordHash,
		PasswordHashAlgorithm:                 impl.Password.AlgorithmName(),
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

func (impl *UserControllerImpl) UpdateByID(ctx context.Context, requestData *UserUpdateRequestIDO) (*user_s.User, error) {
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
		nu, err := impl.userFromUpdateRequest(requestData)
		if err != nil {
			return nil, err
		}

		// Extract from our session the following data.
		userID, _ := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)
		userName, _ := sessCtx.Value(constants.SessionUserName).(string)

		// Extract from our session the following data.
		userRole := sessCtx.Value(constants.SessionUserRole).(int8)

		// Apply filtering based on ownership and role.
		if userRole != user_s.UserRoleRoot {
			return nil, httperror.NewForForbiddenWithSingleField("message", "you do not have permission")
		}

		// Lookup the user in our database, else return a `400 Bad Request` error.
		ou, err := impl.UserStorer.GetByID(sessCtx, nu.ID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if ou == nil {
			impl.Logger.Warn("user does not exist validation error")
			return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
		}

		// Lookup the store in our database, else return a `400 Bad Request` error.
		o, err := impl.TenantStorer.GetByID(sessCtx, nu.TenantID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if o == nil {
			impl.Logger.Warn("store does not exist exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("tenant_id", "store does not exist")
		}

		ou.TenantID = o.ID
		ou.TenantName = o.Name
		ou.FirstName = nu.FirstName
		ou.LastName = nu.LastName
		ou.Name = fmt.Sprintf("%s %s", nu.FirstName, nu.LastName)
		ou.LexicalName = fmt.Sprintf("%s, %s", nu.LastName, nu.FirstName)
		ou.Email = nu.Email
		ou.Phone = nu.Phone
		ou.Country = nu.Country
		ou.Region = nu.Region
		ou.City = nu.City
		ou.PostalCode = nu.PostalCode
		ou.AddressLine1 = nu.AddressLine1
		ou.AddressLine2 = nu.AddressLine2
		ou.HowDidYouHearAboutUs = nu.HowDidYouHearAboutUs
		ou.HowDidYouHearAboutUsOther = nu.HowDidYouHearAboutUsOther
		ou.ModifiedByUserID = userID
		ou.ModifiedByName = userName
		ou.HasShippingAddress = nu.HasShippingAddress
		ou.ShippingName = nu.ShippingName
		ou.ShippingPhone = nu.ShippingPhone
		ou.ShippingCountry = nu.ShippingCountry
		ou.ShippingRegion = nu.ShippingRegion
		ou.ShippingCity = nu.ShippingCity
		ou.ShippingPostalCode = nu.ShippingPostalCode
		ou.ShippingAddressLine1 = nu.ShippingAddressLine1
		ou.ShippingAddressLine2 = nu.ShippingAddressLine2
		ou.HowLongCollectingComicBooksForGrading = nu.HowLongCollectingComicBooksForGrading
		ou.HasPreviouslySubmittedComicBookForGrading = nu.HasPreviouslySubmittedComicBookForGrading
		ou.HasOwnedGradedComicBooks = nu.HasOwnedGradedComicBooks
		ou.HasRegularComicBookShop = nu.HasRegularComicBookShop
		ou.HasPreviouslyPurchasedFromAuctionSite = nu.HasPreviouslyPurchasedFromAuctionSite
		ou.HasPreviouslyPurchasedFromFacebookMarketplace = nu.HasPreviouslyPurchasedFromFacebookMarketplace
		ou.HasRegularlyAttendedComicConsOrCollectibleShows = nu.HasRegularlyAttendedComicConsOrCollectibleShows

		if err := impl.UserStorer.UpdateByID(sessCtx, ou); err != nil {
			impl.Logger.Error("user update by id error", slog.Any("error", err))
			return nil, err
		}

		////
		//// End transaction with success.
		////

		return ou, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return nil, err
	}

	user := res.(*user_s.User)

	////
	//// Update related (in background)
	////

	go func(usr *user_s.User) {
		if err := impl.updateRelatedStoreInBackground(usr); err != nil {
			impl.Logger.Error("update related stores failed", slog.Any("error", err))
		}
	}(user)

	// End.

	return user, nil
}
