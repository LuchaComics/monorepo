package service

import (
	"context"
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
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayUserRegisterService struct {
	config                           *config.Configuration
	logger                           *slog.Logger
	passwordProvider                 password.Provider
	cache                            mongodbcache.Cacher
	jwtProvider                      jwt.Provider
	tenantGetByIDUseCase             *usecase.TenantGetByIDUseCase
	userGetByEmailUseCase            *usecase.UserGetByEmailUseCase
	userCreateUseCase                *usecase.UserCreateUseCase
	userUpdateUseCase                *usecase.UserUpdateUseCase
	sendUserVerificationEmailUseCase *usecase.SendUserVerificationEmailUseCase
}

func NewGatewayUserRegisterService(
	cfg *config.Configuration,
	logger *slog.Logger,
	pp password.Provider,
	cach mongodbcache.Cacher,
	jwtp jwt.Provider,
	uc1 *usecase.TenantGetByIDUseCase,
	uc2 *usecase.UserGetByEmailUseCase,
	uc3 *usecase.UserCreateUseCase,
	uc4 *usecase.UserUpdateUseCase,
	uc5 *usecase.SendUserVerificationEmailUseCase,
) *GatewayUserRegisterService {
	return &GatewayUserRegisterService{cfg, logger, pp, cach, jwtp, uc1, uc2, uc3, uc4, uc5}
}

type RegisterCustomerRequestIDO struct {
	FirstName           string `json:"first_name"`
	LastName            string `json:"last_name"`
	Email               string `json:"email"`
	Password            string `json:"password"`
	PasswordConfirm     string `json:"password_confirm"`
	Phone               string `json:"phone,omitempty"`
	Country             string `json:"country,omitempty"`
	CountryOther        string `json:"country_other,omitempty"`
	Timezone            string `bson:"timezone" json:"timezone"`
	AgreeTermsOfService bool   `json:"agree_terms_of_service,omitempty"`
	AgreePromotions     bool   `json:"agree_promotions,omitempty"`
}

type RegisterCustomerResponseIDO struct {
	User                   *domain.User `json:"user"`
	AccessToken            string       `json:"access_token"`
	AccessTokenExpiryTime  time.Time    `json:"access_token_expiry_time"`
	RefreshToken           string       `json:"refresh_token"`
	RefreshTokenExpiryTime time.Time    `json:"refresh_token_expiry_time"`
}

func (s *GatewayUserRegisterService) Execute(
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
		e["first_name"] = "First name is required"
	}
	if req.LastName == "" {
		e["last_name"] = "Last name is required"
	}
	if req.Email == "" {
		e["email"] = "Email is required"
	}
	if len(req.Email) > 255 {
		e["email"] = "Email is too long"
	}
	if req.Password == "" {
		e["password"] = "Password is required"
	}
	if req.PasswordConfirm == "" {
		e["password_confirm"] = "Password confirm is required"
	}
	if req.PasswordConfirm != req.Password {
		e["password"] = "Password does not match"
		e["password_confirm"] = "Password does not match"
	}
	if req.Country == "" {
		e["country"] = "Country is required"
	} else {
		if req.Country == "Other" && req.CountryOther == "" {
			e["country_other"] = "Specify country is required"
		}
	}
	if req.Timezone == "" {
		e["timezone"] = "Password confirm is required"
	}
	if req.AgreeTermsOfService == false {
		e["agree_terms_of_service"] = "Agreeing to terms of service is required and you must agree to the terms before proceeding"
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
		return nil, httperror.NewForBadRequestWithSingleField("email", "Email address already exists")
	}

	// Create our user.
	u, err = s.createCustomerUserForRequest(sessCtx, req)
	if err != nil {
		s.logger.Error("failed creating customer user error", slog.Any("err", err))
		return nil, err
	}

	// Send our verification email.
	if err := s.sendUserVerificationEmailUseCase.Execute(context.Background(), u); err != nil {
		// s.logger.Error("failed sending verification email with error", slog.Any("err", err))
		// return nil, err

		// Skip any error handling...
	}

	return s.registerWithUser(sessCtx, u)
}

func (s *GatewayUserRegisterService) registerWithUser(sessCtx mongo.SessionContext, u *domain.User) (*RegisterCustomerResponseIDO, error) {
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

func (s *GatewayUserRegisterService) createCustomerUserForRequest(sessCtx mongo.SessionContext, req *RegisterCustomerRequestIDO) (*domain.User, error) {
	// Lookup the store and check to see if it's active or not, if not active then return the specific requests.
	t, err := s.tenantGetByIDUseCase.Execute(sessCtx, s.config.App.TenantID)
	if err != nil {
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if t == nil {
		err := fmt.Errorf("Tenant does not exist for ID: %v", s.config.App.TenantID)
		s.logger.Error("database error", slog.Any("err", err))
		return nil, err
	}

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

	ipAddress, _ := sessCtx.Value(constants.SessionIPAddress).(string)

	userID := primitive.NewObjectID()
	u := &domain.User{
		TenantID:                   t.ID,
		TenantName:                 t.Name,
		ID:                         userID,
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		Name:                       fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:                fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                      req.Email,
		PasswordHash:               passwordHash,
		PasswordHashAlgorithm:      s.passwordProvider.AlgorithmName(),
		Role:                       domain.UserRoleCustomer,
		Phone:                      "",
		Country:                    req.Country,
		Timezone:                   req.Timezone,
		Region:                     "",
		City:                       "",
		PostalCode:                 "",
		AddressLine1:               "",
		AddressLine2:               "",
		AgreeTermsOfService:        req.AgreeTermsOfService,
		AgreePromotions:            req.AgreePromotions,
		CreatedByUserID:            userID,
		CreatedAt:                  time.Now(),
		CreatedByName:              fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		CreatedFromIPAddress:       ipAddress,
		ModifiedByUserID:           userID,
		ModifiedAt:                 time.Now(),
		ModifiedByName:             fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		ModifiedFromIPAddress:      ipAddress,
		WasEmailVerified:           false,
		EmailVerificationCode:      primitive.NewObjectID().Hex(),
		EmailVerificationExpiry:    time.Now().Add(72 * time.Hour),
		Status:                     domain.UserStatusActive,
		HasShippingAddress:         false,
		ShippingName:               "",
		ShippingPhone:              "",
		ShippingCountry:            "",
		ShippingRegion:             "",
		ShippingCity:               "",
		ShippingPostalCode:         "",
		ShippingAddressLine1:       "",
		ShippingAddressLine2:       "",
		PaymentProcessorName:       "",
		PaymentProcessorCustomerID: "",
		// HowDidYouHearAboutUs:                  "",
		// HowDidYouHearAboutUsOther:             "",
		// HowLongCollectingComicBooksForGrading: "",
		// HasPreviouslySubmittedComicBookForGrading:       "",
		// HasOwnedGradedComicBooks:                        "",
		// HasRegularComicBookShop:                         "",
		// HasPreviouslyPurchasedFromAuctionSite:           "",
		// HasPreviouslyPurchasedFromFacebookMarketplace:   "",
		// HasRegularlyAttendedComicConsOrCollectibleShows: "",
	}
	if req.CountryOther != "" {
		u.Country = req.CountryOther
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
