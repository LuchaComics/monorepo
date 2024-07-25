package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	gateway_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/gateway/datastore"
	store_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/store/datastore"
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func validateRegisterBusinessRequest(dirtyData *gateway_s.RegisterBusinessRequestIDO) error {
	e := make(map[string]string)

	if dirtyData.FirstName == "" {
		e["first_name"] = "missing value"
	}
	if dirtyData.LastName == "" {
		e["last_name"] = "missing value"
	}
	if dirtyData.Email == "" {
		e["email"] = "missing value"
	}
	if len(dirtyData.Email) > 255 {
		e["email"] = "too long"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.PasswordRepeated == "" {
		e["password_repeated"] = "missing value"
	}
	if dirtyData.PasswordRepeated != dirtyData.Password {
		e["password"] = "does not match"
		e["password_repeated"] = "does not match"
	}
	if dirtyData.ComicBookStoreName == "" {
		e["comic_book_store_name"] = "missing value"
	}
	if dirtyData.Phone == "" {
		e["phone"] = "missing value"
	}
	if dirtyData.Country == "" {
		e["country"] = "missing value"
	}
	if dirtyData.Region == "" {
		e["region"] = "missing value"
	}
	if dirtyData.City == "" {
		e["city"] = "missing value"
	}
	if dirtyData.PostalCode == "" {
		e["postal_code"] = "missing value"
	}
	if dirtyData.Password == "" {
		e["password"] = "missing value"
	}
	if dirtyData.AddressLine1 == "" {
		e["address_line1"] = "missing value"
	}
	// if dirtyData.HowDidYouHearAboutUs == 0 {
	// 	e["how_did_you_hear_about_us"] = "missing value"
	// }
	if dirtyData.AgreeTOS == false {
		e["agree_tos"] = "you must agree to the terms before proceeding"
	}
	if dirtyData.HowDidYouHearAboutUs > 7 || dirtyData.HowDidYouHearAboutUs < 1 {
		e["how_did_you_hear_about_us"] = "missing value"
	} else {
		if dirtyData.HowDidYouHearAboutUs == 1 && dirtyData.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "missing value"
		}
	}
	if dirtyData.HowLongStoreOperating > 4 || dirtyData.HowLongStoreOperating < 1 {
		e["how_long_store_operating"] = "missing value"
	}

	// The following logic will enforce shipping address input validation.
	if dirtyData.HasShippingAddress == true {
		if dirtyData.ShippingName == "" {
			e["shipping_name"] = "missing value"
		}
		if dirtyData.ShippingPhone == "" {
			e["shipping_phone"] = "missing value"
		}
		if dirtyData.ShippingCountry == "" {
			e["shipping_country"] = "missing value"
		}
		if dirtyData.ShippingRegion == "" {
			e["shipping_region"] = "missing value"
		}
		if dirtyData.ShippingCity == "" {
			e["shipping_city"] = "missing value"
		}
		if dirtyData.ShippingPostalCode == "" {
			e["shipping_postal_code"] = "missing value"
		}
		if dirtyData.ShippingAddressLine1 == "" {
			e["shipping_address_line1"] = "missing value"
		}
	}
	if dirtyData.RetailPartnershipReason == "" {
		e["retail_partnership_reason"] = "missing value"
	}
	if dirtyData.CPSPartnershipReason == "" {
		e["cps_partnership_reason"] = "missing value"
	}
	if dirtyData.EstimatedSubmissionsPerMonth == 0 {
		e["estimated_submissions_per_month"] = "missing value"
	}
	if dirtyData.HasOtherGradingService == 0 {
		e["has_other_grading_service"] = "missing value"
	} else {
		// if dirtyData.OtherGradingServiceName == "" {
		// 	e["other_grading_service_name"] = "missing value"
		// }
	}
	if dirtyData.RequestWelcomePackage == 0 {
		e["request_welcome_package"] = "missing value"
	}
	if dirtyData.Timezone == "" {
		e["timezone"] = "missing value"
	} else {
		// Confirm the timezone is one that exists.
		location, err := time.LoadLocation(dirtyData.Timezone)
		if err != nil || location == nil {
			e["timezone"] = "unsupported value"
		}
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *GatewayControllerImpl) RegisterBusiness(ctx context.Context, req *gateway_s.RegisterBusinessRequestIDO) error {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	req.Email = strings.ToLower(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err := validateRegisterBusinessRequest(req); err != nil {
		return err
	}

	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {

		// Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByEmail(sessCtx, req.Email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u != nil {
			impl.Logger.Warn("user already exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
		}

		// Create our user.
		u, err = impl.createBusinessUserForRequest(sessCtx, req)
		if err != nil {
			return nil, err
		}

		// Create our store.
		orgID, err := impl.createStoreForUser(sessCtx, req, u)
		if err != nil {
			impl.Logger.Error("database create error", slog.Any("error", err))
			return nil, err
		}
		u.StoreID = orgID // Attach to our user profile.
		u.StoreName = req.ComicBookStoreName
		u.StoreLevel = 1
		u.StoreTimezone = req.Timezone
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("database update error", slog.Any("error", err))
			return nil, err
		}

		// Send our verification email.
		if err := impl.TemplatedEmailer.SendBusinessVerificationEmail(u.Email, u.EmailVerificationCode, u.FirstName); err != nil {
			impl.Logger.Error("failed sending verification email with error", slog.Any("err", err))
			return nil, err
		}

		return nil, nil
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}

func (impl *GatewayControllerImpl) createBusinessUserForRequest(sessCtx mongo.SessionContext, req *gateway_s.RegisterBusinessRequestIDO) (*user_s.User, error) {
	passwordHash, err := impl.Password.GenerateHashFromPassword(req.Password)
	if err != nil {
		impl.Logger.Error("hashing error", slog.Any("error", err))
		return nil, err
	}

	// Create an account with our payment processor.
	var paymentProcessorCustomerID *string
	if "Stripe, Inc." == impl.PaymentProcessor.GetName() {
		paymentProcessorCustomerID, err = impl.PaymentProcessor.CreateCustomer(
			fmt.Sprintf("%s %s", req.FirstName, req.LastName),
			req.Email,
			"", // description...
			fmt.Sprintf("%s %s Shipping Address", req.FirstName, req.LastName),
			req.Phone,
			req.ShippingCity, req.ShippingCountry, req.ShippingAddressLine1, req.ShippingAddressLine2, req.ShippingPostalCode, req.ShippingRegion, // Shipping
			req.City, req.Country, req.AddressLine1, req.AddressLine2, req.PostalCode, req.Region, // Billing
		)
		if err != nil {
			impl.Logger.Error("creating customer from payment processor error", slog.Any("error", err))
			return nil, err
		}
	}

	userID := primitive.NewObjectID()
	u := &user_s.User{
		ID:                         userID,
		FirstName:                  req.FirstName,
		LastName:                   req.LastName,
		Name:                       fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:                fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                      req.Email,
		PasswordHash:               passwordHash,
		PasswordHashAlgorithm:      impl.Password.AlgorithmName(),
		Role:                       user_s.UserRoleRetailer,
		Phone:                      req.Phone,
		Country:                    req.Country,
		Region:                     req.Region,
		City:                       req.City,
		PostalCode:                 req.PostalCode,
		AddressLine1:               req.AddressLine1,
		AddressLine2:               req.AddressLine2,
		HowDidYouHearAboutUs:       req.HowDidYouHearAboutUs,
		HowDidYouHearAboutUsOther:  req.HowDidYouHearAboutUsOther,
		AgreeTOS:                   req.AgreeTOS,
		AgreePromotionsEmail:       req.AgreePromotionsEmail,
		CreatedByUserID:            userID,
		CreatedAt:                  time.Now(),
		CreatedByName:              fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		ModifiedByUserID:           userID,
		ModifiedAt:                 time.Now(),
		ModifiedByName:             fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		WasEmailVerified:           false,
		EmailVerificationCode:      impl.UUID.NewUUID(),
		EmailVerificationExpiry:    time.Now().Add(72 * time.Hour),
		Status:                     user_s.UserStatusActive,
		HasShippingAddress:         req.HasShippingAddress,
		ShippingName:               req.ShippingName,
		ShippingPhone:              req.ShippingPhone,
		ShippingCountry:            req.ShippingCountry,
		ShippingRegion:             req.ShippingRegion,
		ShippingCity:               req.ShippingCity,
		ShippingPostalCode:         req.ShippingPostalCode,
		ShippingAddressLine1:       req.ShippingAddressLine1,
		ShippingAddressLine2:       req.ShippingAddressLine2,
		PaymentProcessorName:       impl.PaymentProcessor.GetName(),
		PaymentProcessorCustomerID: *paymentProcessorCustomerID,
	}
	err = impl.UserStorer.Create(sessCtx, u)
	if err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}
	impl.Logger.Info("Business user created.",
		slog.Any("_id", u.ID),
		slog.String("full_name", u.Name),
		slog.String("email", u.Email),
		slog.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		slog.String("password_hash", u.PasswordHash))

	return u, nil
}

func (impl *GatewayControllerImpl) createStoreForUser(sessCtx mongo.SessionContext, req *gateway_s.RegisterBusinessRequestIDO, u *user_s.User) (primitive.ObjectID, error) {
	o := &store_s.Store{
		ID:                           primitive.NewObjectID(),
		Name:                         req.ComicBookStoreName,
		WebsiteURL:                   req.WebsiteURL,
		EstimatedSubmissionsPerMonth: req.EstimatedSubmissionsPerMonth,
		HasOtherGradingService:       req.HasOtherGradingService,
		OtherGradingServiceName:      req.OtherGradingServiceName,
		RequestWelcomePackage:        req.RequestWelcomePackage,
		GradingComicsExperience:      req.GradingComicsExperience,
		RetailPartnershipReason:      req.RetailPartnershipReason,
		CPSPartnershipReason:         req.CPSPartnershipReason,
		HowLongStoreOperating:        req.HowLongStoreOperating,
		Type:                         store_s.RetailerType,
		CreatedAt:                    time.Now(),
		CreatedByUserID:              u.ID,
		CreatedByUserName:            u.Name,
		ModifiedAt:                   time.Now(),
		ModifiedByUserID:             u.ID,
		ModifiedByUserName:           u.Name,
		Status:                       store_s.StorePendingStatus,
		Level:                        1, // Default
		Timezone:                     req.Timezone,
	}
	err := impl.StoreStorer.Create(sessCtx, o)
	if err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return primitive.NewObjectID(), err
	}
	impl.Logger.Info("Store created.",
		slog.Any("_id", u.ID),
		slog.String("name", u.Name))

	return o.ID, nil
}
