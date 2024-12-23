package service

import (
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayProfileApplyForVerificationService struct {
	logger             *slog.Logger
	userGetByIDUseCase *usecase.UserGetByIDUseCase
	userUpdateUseCase  *usecase.UserUpdateUseCase
}

func NewGatewayProfileApplyForVerificationService(
	logger *slog.Logger,
	uc1 *usecase.UserGetByIDUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *GatewayProfileApplyForVerificationService {
	return &GatewayProfileApplyForVerificationService{logger, uc1, uc2}
}

type GatewayProfileApplyForVerificationRequestIDO struct {
	Phone                                           string `bson:"phone" json:"phone,omitempty"`
	Country                                         string `bson:"country" json:"country,omitempty"`
	Region                                          string `bson:"region" json:"region,omitempty"`
	City                                            string `bson:"city" json:"city,omitempty"`
	PostalCode                                      string `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1                                    string `bson:"address_line1" json:"address_line1,omitempty"`
	AddressLine2                                    string `bson:"address_line2" json:"address_line2,omitempty"`
	HasShippingAddress                              bool   `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName                                    string `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone                                   string `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry                                 string `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion                                  string `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity                                    string `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode                              string `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1                            string `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2                            string `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	HowDidYouHearAboutUs                            int8   `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther                       string `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	HowLongCollectingComicBooksForGrading           int8   `bson:"how_long_collecting_comic_books_for_grading" json:"how_long_collecting_comic_books_for_grading"`
	HasPreviouslySubmittedComicBookForGrading       int8   `bson:"has_previously_submitted_comic_book_for_grading" json:"has_previously_submitted_comic_book_for_grading"`
	HasOwnedGradedComicBooks                        int8   `bson:"has_owned_graded_comic_books" json:"has_owned_graded_comic_books"`
	HasRegularComicBookShop                         int8   `bson:"has_regular_comic_book_shop" json:"has_regular_comic_book_shop"`
	HasPreviouslyPurchasedFromAuctionSite           int8   `bson:"has_previously_purchased_from_auction_site" json:"has_previously_purchased_from_auction_site"`
	HasPreviouslyPurchasedFromFacebookMarketplace   int8   `bson:"has_previously_purchased_from_facebook_marketplace" json:"has_previously_purchased_from_facebook_marketplace"`
	HasRegularlyAttendedComicConsOrCollectibleShows int8   `bson:"has_regularly_attended_comic_cons_or_collectible_shows" json:"has_regularly_attended_comic_cons_or_collectible_shows"`
}

func (s *GatewayProfileApplyForVerificationService) Execute(sessCtx mongo.SessionContext, req *GatewayProfileApplyForVerificationRequestIDO) (*domain.User, error) {
	//
	// STEP 1: Validation
	//

	e := make(map[string]string)

	if req == nil {
		err := errors.New("No request data inputted")
		s.logger.Error("validation error", slog.Any("err", err))
		return nil, err
	}

	if req.Phone == "" {
		e["phone"] = "Phone is required"
	}
	if req.Country == "" {
		e["country"] = "Country is required"
	}
	if req.Region == "" {
		e["region"] = "Region is required"
	}
	if req.City == "" {
		e["city"] = "City is required"
	}
	if req.PostalCode == "" {
		e["postal_code"] = "Postal code is required"
	}
	if req.AddressLine1 == "" {
		e["address_line1"] = "Address line 1 is required"
	}
	if req.HasShippingAddress {
		if req.ShippingName == "" {
			e["shipping_name"] = "Shipping name is required"
		}
		if req.ShippingPhone == "" {
			e["shipping_phone"] = "Shipping phone is required"
		}
		if req.ShippingCountry == "" {
			e["shipping_country"] = "Shipping country is required"
		}
		if req.ShippingRegion == "" {
			e["shipping_region"] = "Shipping region is required"
		}
		if req.ShippingCity == "" {
			e["shipping_city"] = "Shipping city is required"
		}
		if req.ShippingPostalCode == "" {
			e["shipping_postal_code"] = "Shipping postal code is required"
		}
		if req.ShippingAddressLine1 == "" {
			e["shipping_address_line1"] = "Shipping address line 1 is required"
		}
	}
	if req.HowDidYouHearAboutUs == 0 {
		e["how_did_you_hear_about_us"] = "How did you hear about us is required"
	} else {
		if req.HowDidYouHearAboutUs == 1 && req.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "How did you hear about us (other) is required"
		}
	}
	if req.HowLongCollectingComicBooksForGrading == 0 {
		e["how_long_collecting_comic_books_for_grading"] = "How long collecting comic books for grading is required"
	}
	if req.HasPreviouslySubmittedComicBookForGrading == 0 {
		e["how_long_collecting_comic_books_for_grading"] = "Has previously submitted comic book for grading is required"
	}
	if req.HasOwnedGradedComicBooks == 0 {
		e["has_owned_graded_comic_books"] = "Has owned graded comic books is required"
	}
	if req.HasRegularComicBookShop == 0 {
		e["has_regular_comic_book_shop"] = "Has regular comic book shop is required"
	}
	if req.HasPreviouslyPurchasedFromAuctionSite == 0 {
		e["has_previously_purchased_from_auction_site"] = "Has previously purchased from auction site required"
	}
	if req.HasPreviouslyPurchasedFromFacebookMarketplace == 0 {
		e["has_previously_purchased_from_facebook_marketplace"] = "Has previously purchased from facebook marketplace required"
	}
	if req.HasRegularlyAttendedComicConsOrCollectibleShows == 0 {
		e["has_regularly_attended_comic_cons_or_collectible_shows"] = "Has regularly attended comic cons or collectible shows is required"
	}

	if len(e) != 0 {
		s.logger.Warn("validation failure",
			slog.Any("e", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	// Extract from our session the following data.
	userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	ou, err := s.userGetByIDUseCase.Execute(sessCtx, userID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if ou == nil {
		s.logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}

	ou.DidApplyForVerification = true
	ou.Phone = req.Phone
	ou.Country = req.Country
	ou.Region = req.Region
	ou.City = req.City
	ou.PostalCode = req.PostalCode
	ou.AddressLine1 = req.AddressLine1
	ou.AddressLine2 = req.AddressLine2
	ou.HowDidYouHearAboutUs = req.HowDidYouHearAboutUs
	ou.HowDidYouHearAboutUsOther = req.HowDidYouHearAboutUsOther
	ou.HasShippingAddress = req.HasShippingAddress
	ou.ShippingName = req.ShippingName
	ou.ShippingPhone = req.ShippingPhone
	ou.ShippingCountry = req.ShippingCountry
	ou.ShippingRegion = req.ShippingRegion
	ou.ShippingCity = req.ShippingCity
	ou.ShippingPostalCode = req.ShippingPostalCode
	ou.ShippingAddressLine1 = req.ShippingAddressLine1
	ou.ShippingAddressLine2 = req.ShippingAddressLine2
	ou.HowLongCollectingComicBooksForGrading = req.HowLongCollectingComicBooksForGrading
	ou.HasPreviouslySubmittedComicBookForGrading = req.HasPreviouslySubmittedComicBookForGrading
	ou.HasOwnedGradedComicBooks = req.HasOwnedGradedComicBooks
	ou.HasRegularComicBookShop = req.HasRegularComicBookShop
	ou.HasPreviouslyPurchasedFromAuctionSite = req.HasPreviouslyPurchasedFromAuctionSite
	ou.HasPreviouslyPurchasedFromFacebookMarketplace = req.HasPreviouslyPurchasedFromFacebookMarketplace
	ou.HasRegularlyAttendedComicConsOrCollectibleShows = req.HasRegularlyAttendedComicConsOrCollectibleShows

	if err := s.userUpdateUseCase.Execute(sessCtx, ou); err != nil {
		s.logger.Error("user update by id error", slog.Any("error", err))
		return nil, err
	}

	return ou, nil
}
