package service

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/jwt"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/password"
	sstring "github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/security/securestring"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/storage/database/mongodbcache"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayRegisterCustomerService struct {
	config                *config.Configuration
	logger                *slog.Logger
	passwordProvider      password.Provider
	cache                 mongodbcache.Cacher
	jwtProvider           jwt.Provider
	tenantGetByIDUseCase  *usecase.TenantGetByIDUseCase
	userGetByEmailUseCase *usecase.UserGetByEmailUseCase
	userCreateUseCase     *usecase.UserCreateUseCase
	userUpdateUseCase     *usecase.UserUpdateUseCase
}

func NewGatewayRegisterCustomerService(
	cfg *config.Configuration,
	logger *slog.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 *usecase.TenantGetByIDUseCase,
	uc2 *usecase.UserGetByEmailUseCase,
	uc3 *usecase.UserCreateUseCase,
	uc4 *usecase.UserUpdateUseCase,
) *GatewayRegisterCustomerService {
	return &GatewayRegisterCustomerService{cfg, logger, pp, cach, jwtp, uc1, uc2, uc3, uc4}
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
	StoreID                                         primitive.ObjectID `json:"store_id"`
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
	User                   *domain.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

func (s *GatewayRegisterCustomerService) Execute(
	sessCtx mongo.SessionContext,
	req *RegisterCustomerRequestIDO,
) (*RegisterCustomerResponseIDO, error) {
	//
	// STEP 1: Sanization of input.
	//

	// Defensive Code: For security purposes we need to perform some sanitization on the inputs.
	req.Email = strings.ToLower(req.Email)
	req.Email = strings.ReplaceAll(req.Email, " ", "")
	req.Email = strings.ReplaceAll(req.Email, "\t", "")
	req.Email = strings.TrimSpace(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")
	req.Password = strings.ReplaceAll(req.Password, "\t", "")
	req.Password = strings.TrimSpace(req.Password)
	// password, err := sstring.NewSecureString(unsecurePassword)
	// if err != nil {
	// 	s.logger.Error("secure string error", slog.Any("err", err))
	// 	return nil, err
	// }

	//
	// STEP 2: Validation of input.
	//

	e := make(map[string]string)
	if req.FirstName == "" {
		e["first_name"] = "missing value"
	}
	if req.LastName == "" {
		e["last_name"] = "missing value"
	}
	if req.Email == "" {
		e["email"] = "missing value"
	}
	if len(req.Email) > 255 {
		e["email"] = "too long"
	}
	if req.Password == "" {
		e["password"] = "missing value"
	}
	if req.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if req.PasswordRepeated != req.Password {
		e["password"] = "does not match"
		e["password_repeated"] = "does not match"
	}
	if req.Phone == "" {
		e["phone"] = "missing value"
	}
	if req.Country == "" {
		e["country"] = "missing value"
	}
	if req.Region == "" {
		e["region"] = "missing value"
	}
	if req.City == "" {
		e["city"] = "missing value"
	}
	if req.PostalCode == "" {
		e["postal_code"] = "missing value"
	}
	if req.Password == "" {
		e["password"] = "missing value"
	}
	if req.AddressLine1 == "" {
		e["address_line1"] = "missing value"
	}
	if req.AgreeTOS == false {
		e["agree_tos"] = "you must agree to the terms before proceeding"
	}
	// The following logic will enforce shipping address input validation.
	if req.HasShippingAddress == true {
		if req.ShippingName == "" {
			e["shipping_name"] = "missing value"
		}
		if req.ShippingPhone == "" {
			e["shipping_phone"] = "missing value"
		}
		if req.ShippingCountry == "" {
			e["shipping_country"] = "missing value"
		}
		if req.ShippingRegion == "" {
			e["shipping_region"] = "missing value"
		}
		if req.ShippingCity == "" {
			e["shipping_city"] = "missing value"
		}
		if req.ShippingPostalCode == "" {
			e["shipping_postal_code"] = "missing value"
		}
		if req.ShippingAddressLine1 == "" {
			e["shipping_address_line1"] = "missing value"
		}
	}
	// if req.StoreID.IsZero() {
	// 	e["store_id"] = "missing value"
	// } else {
	// 	if exists, err := s.StoreStorer.CheckIfExistsByID(ctx, req.StoreID); exists == false || err != nil {
	// 		if err != nil {
	// 			e["store_id"] = fmt.Sprintf("error occured: %s", err.Error())
	// 		} else if exists == false {
	// 			e["store_id"] = fmt.Sprintf("does not exist for value: %s", req.StoreID.Hex())
	// 		}
	// 	}
	// }
	if req.HowDidYouHearAboutUs > 7 || req.HowDidYouHearAboutUs < 1 {
		e["how_did_you_hear_about_us"] = "missing value"
	} else {
		if req.HowDidYouHearAboutUs == 1 && req.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "missing value"
		}
	}
	if req.HowLongCollectingComicBooksForGrading == 0 {
		e["how_long_collecting_comic_books_for_grading"] = "missing value"
	}
	if req.HasPreviouslySubmittedComicBookForGrading == 0 {
		e["has_previously_submitted_comic_book_for_grading"] = "missing value"
	}
	if req.HasOwnedGradedComicBooks == 0 {
		e["has_owned_graded_comic_books"] = "missing value"
	}
	if req.HasRegularComicBookShop == 0 {
		e["has_regular_comic_book_shop"] = "missing value"
	}
	if req.HasPreviouslyPurchasedFromAuctionSite == 0 {
		e["has_previously_purchased_from_auction_site"] = "missing value"
	}
	if req.HasPreviouslyPurchasedFromFacebookMarketplace == 0 {
		e["has_previously_purchased_from_facebook_marketplace"] = "missing value"
	}
	if req.HasRegularlyAttendedComicConsOrCollectibleShows == 0 {
		e["has_regularly_attended_comic_cons_or_collectible_shows"] = "missing value"
	}
	if len(e) != 0 {
		s.logger.Warn("Failed validation register",
			slog.Any("error", e))
		return nil, httperror.NewForBadRequest(&e)
	}

	//
	// STEP 3:
	//

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := s.userGetByEmailUseCase.Execute(sessCtx, req.Email)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u != nil {
		s.logger.Warn("user already exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("email", "already exists")
	}

	// Lookup the store and check to see if it's active or not, if not active then return the specific requests.
	t, err := s.tenantGetByIDUseCase.Execute(sessCtx, s.config.App.TenantID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if t == nil {
		err := fmt.Errorf("Tenant does not exist for ID: %v", u.TenantID.Hex())
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}

	// Create our user.
	u, err = s.createCustomerUserForRequest(sessCtx, req)
	if err != nil {
		s.logger.Error("failed creating customer user error", slog.Any("err", err))
		return nil, err
	}

	// // Send our verification email.
	// if err := impl.TemplatedEmailer.SendCustomerVerificationEmail(u.Email, u.EmailVerificationCode, u.FirstName); err != nil {
	// 	impl.Logger.Error("failed sending verification email with error", slog.Any("err", err))
	// 	return nil, err
	// }

	return s.registerWithUser(sessCtx, u)
}

func (s *GatewayRegisterCustomerService) registerWithUser(sessCtx mongo.SessionContext, u *domain.User) (*RegisterCustomerResponseIDO, error) {
	uBin, err := json.Marshal(u)
	if err != nil {
		s.logger.Error("marshalling error", slog.Any("err", err))
		return nil, err
	}

	// Set expiry duration.
	atExpiry := 24 * time.Hour
	rtExpiry := 14 * 24 * time.Hour

	// Start our session using an access and refresh token.
	sessionUUID := primitive.NewObjectID().Hex()

	err = s.cache.SetWithExpiry(sessCtx, sessionUUID, uBin, rtExpiry)
	if err != nil {
		s.logger.Error("cache set with expiry error", slog.Any("err", err))
		return nil, err
	}

	// Generate our JWT token.
	accessToken, accessTokenExpiry, refreshToken, refreshTokenExpiry, err := s.jwtProvider.GenerateJWTTokenPair(sessionUUID, atExpiry, rtExpiry)
	if err != nil {
		s.logger.Error("jwt generate pairs error", slog.Any("err", err))
		return nil, err
	}

	// Return our auth keys.
	return &RegisterCustomerResponseIDO{
		User:                   u,
		AccessToken:            accessToken,
		AccessTokenExpiryTime:  accessTokenExpiry,
		RefreshToken:           refreshToken,
		RefreshTokenExpiryTime: refreshTokenExpiry,
	}, nil
}

func (s *GatewayRegisterCustomerService) createCustomerUserForRequest(sessCtx mongo.SessionContext, req *RegisterCustomerRequestIDO) (*domain.User, error) {
	password, err := sstring.NewSecureString(req.Password)
	if err != nil {
		s.logger.Error("password securing error", slog.Any("err", err))
		return nil, err
	}

	passwordHash, err := s.passwordProvider.GenerateHashFromPassword(password)
	if err != nil {
		s.logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	userID := primitive.NewObjectID()
	u := &domain.User{
		ID:                                    userID,
		FirstName:                             req.FirstName,
		LastName:                              req.LastName,
		Name:                                  fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:                           fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                                 req.Email,
		PasswordHash:                          passwordHash,
		PasswordHashAlgorithm:                 s.passwordProvider.AlgorithmName(),
		Role:                                  domain.UserRoleCustomer,
		Phone:                                 req.Phone,
		Country:                               req.Country,
		Region:                                req.Region,
		City:                                  req.City,
		PostalCode:                            req.PostalCode,
		AddressLine1:                          req.AddressLine1,
		AddressLine2:                          req.AddressLine2,
		AgreeTOS:                              req.AgreeTOS,
		AgreePromotionsEmail:                  req.AgreePromotionsEmail,
		CreatedByUserID:                       userID,
		CreatedAt:                             time.Now(),
		CreatedByName:                         fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		ModifiedByUserID:                      userID,
		ModifiedAt:                            time.Now(),
		ModifiedByName:                        fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		WasEmailVerified:                      false,
		EmailVerificationCode:                 primitive.NewObjectID().Hex(),
		EmailVerificationExpiry:               time.Now().Add(72 * time.Hour),
		Status:                                domain.UserStatusActive,
		HasShippingAddress:                    req.HasShippingAddress,
		ShippingName:                          req.ShippingName,
		ShippingPhone:                         req.ShippingPhone,
		ShippingCountry:                       req.ShippingCountry,
		ShippingRegion:                        req.ShippingRegion,
		ShippingCity:                          req.ShippingCity,
		ShippingPostalCode:                    req.ShippingPostalCode,
		ShippingAddressLine1:                  req.ShippingAddressLine1,
		ShippingAddressLine2:                  req.ShippingAddressLine2,
		PaymentProcessorName:                  "",
		PaymentProcessorCustomerID:            "",
		HowDidYouHearAboutUs:                  req.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther:             req.HowDidYouHearAboutUsOther,
		HowLongCollectingComicBooksForGrading: req.HowLongCollectingComicBooksForGrading,
		HasPreviouslySubmittedComicBookForGrading:       req.HasPreviouslySubmittedComicBookForGrading,
		HasOwnedGradedComicBooks:                        req.HasOwnedGradedComicBooks,
		HasRegularComicBookShop:                         req.HasRegularComicBookShop,
		HasPreviouslyPurchasedFromAuctionSite:           req.HasPreviouslyPurchasedFromAuctionSite,
		HasPreviouslyPurchasedFromFacebookMarketplace:   req.HasPreviouslyPurchasedFromFacebookMarketplace,
		HasRegularlyAttendedComicConsOrCollectibleShows: req.HasRegularlyAttendedComicConsOrCollectibleShows,
	}
	err = s.userCreateUseCase.Execute(sessCtx, u)
	if err != nil {
		s.logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}
	s.logger.Info("Customer user created.",
		slog.Any("_id", u.ID),
		slog.String("full_name", u.Name),
		slog.String("email", u.Email),
		slog.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		slog.String("password_hash", u.PasswordHash))

	return u, nil
}
