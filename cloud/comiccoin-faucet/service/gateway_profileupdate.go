package service

import (
	"fmt"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/common/httperror"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/config/constants"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/domain"
	"github.com/LuchaComics/monorepo/cloud/comiccoin-faucet/usecase"
)

type GatewayProfileUpdateService struct {
	logger             *slog.Logger
	userGetByIDUseCase *usecase.UserGetByIDUseCase
	userUpdateUseCase  *usecase.UserUpdateUseCase
}

func NewGatewayProfileUpdateService(
	logger *slog.Logger,
	uc1 *usecase.UserGetByIDUseCase,
	uc2 *usecase.UserUpdateUseCase,
) *GatewayProfileUpdateService {
	return &GatewayProfileUpdateService{logger, uc1, uc2}
}

type GatewayProfileUpdateRequestIDO struct {
	FirstName                 string `bson:"first_name" json:"first_name"`
	LastName                  string `bson:"last_name" json:"last_name"`
	Email                     string `bson:"email" json:"email"`
	Phone                     string `bson:"phone" json:"phone,omitempty"`
	Country                   string `bson:"country" json:"country,omitempty"`
	Region                    string `bson:"region" json:"region,omitempty"`
	City                      string `bson:"city" json:"city,omitempty"`
	PostalCode                string `bson:"postal_code" json:"postal_code,omitempty"`
	AddressLine1              string `bson:"address_line1" json:"address_line1,omitempty"`
	AddressLine2              string `bson:"address_line2" json:"address_line2,omitempty"`
	HasShippingAddress        bool   `bson:"has_shipping_address" json:"has_shipping_address,omitempty"`
	ShippingName              string `bson:"shipping_name" json:"shipping_name,omitempty"`
	ShippingPhone             string `bson:"shipping_phone" json:"shipping_phone,omitempty"`
	ShippingCountry           string `bson:"shipping_country" json:"shipping_country,omitempty"`
	ShippingRegion            string `bson:"shipping_region" json:"shipping_region,omitempty"`
	ShippingCity              string `bson:"shipping_city" json:"shipping_city,omitempty"`
	ShippingPostalCode        string `bson:"shipping_postal_code" json:"shipping_postal_code,omitempty"`
	ShippingAddressLine1      string `bson:"shipping_address_line1" json:"shipping_address_line1,omitempty"`
	ShippingAddressLine2      string `bson:"shipping_address_line2" json:"shipping_address_line2,omitempty"`
	HowDidYouHearAboutUs      int8   `bson:"how_did_you_hear_about_us" json:"how_did_you_hear_about_us,omitempty"`
	HowDidYouHearAboutUsOther string `bson:"how_did_you_hear_about_us_other" json:"how_did_you_hear_about_us_other,omitempty"`
	AgreePromotions      bool   `bson:"agree_promotions_email" json:"agree_promotions_email,omitempty"`
	agreeTermsOfService                  bool   `bson:"agree_terms_of_service" json:"agree_terms_of_service,omitempty"`
}

func (s *GatewayProfileUpdateService) Execute(sessCtx mongo.SessionContext, nu *GatewayProfileUpdateRequestIDO) (*domain.User, error) {
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
	ou.AgreePromotions = nu.AgreePromotions
	ou.HasShippingAddress = nu.HasShippingAddress
	ou.ShippingName = nu.ShippingName
	ou.ShippingPhone = nu.ShippingPhone
	ou.ShippingCountry = nu.ShippingCountry
	ou.ShippingRegion = nu.ShippingRegion
	ou.ShippingCity = nu.ShippingCity
	ou.ShippingPostalCode = nu.ShippingPostalCode
	ou.ShippingAddressLine1 = nu.ShippingAddressLine1
	ou.ShippingAddressLine2 = nu.ShippingAddressLine2

	if err := s.userUpdateUseCase.Execute(sessCtx, ou); err != nil {
		s.logger.Error("user update by id error", slog.Any("error", err))
		return nil, err
	}

	return ou, nil
}
