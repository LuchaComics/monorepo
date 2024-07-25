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
	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) validateRegisterCustomerRequest(ctx context.Context, dirtyData *gateway_s.RegisterCustomerRequestIDO) error {
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
	if dirtyData.AgreeTOS == false {
		e["agree_tos"] = "you must agree to the terms before proceeding"
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
	if dirtyData.StoreID.IsZero() {
		e["store_id"] = "missing value"
	} else {
		if exists, err := impl.StoreStorer.CheckIfExistsByID(ctx, dirtyData.StoreID); exists == false || err != nil {
			if err != nil {
				e["store_id"] = fmt.Sprintf("error occured: %s", err.Error())
			} else if exists == false {
				e["store_id"] = fmt.Sprintf("does not exist for value: %s", dirtyData.StoreID.Hex())
			}
		}
	}
	if dirtyData.HowDidYouHearAboutUs > 7 || dirtyData.HowDidYouHearAboutUs < 1 {
		e["how_did_you_hear_about_us"] = "missing value"
	} else {
		if dirtyData.HowDidYouHearAboutUs == 1 && dirtyData.HowDidYouHearAboutUsOther == "" {
			e["how_did_you_hear_about_us_other"] = "missing value"
		}
	}
	if dirtyData.HowLongCollectingComicBooksForGrading == 0 {
		e["how_long_collecting_comic_books_for_grading"] = "missing value"
	}
	if dirtyData.HasPreviouslySubmittedComicBookForGrading == 0 {
		e["has_previously_submitted_comic_book_for_grading"] = "missing value"
	}
	if dirtyData.HasOwnedGradedComicBooks == 0 {
		e["has_owned_graded_comic_books"] = "missing value"
	}
	if dirtyData.HasRegularComicBookShop == 0 {
		e["has_regular_comic_book_shop"] = "missing value"
	}
	if dirtyData.HasPreviouslyPurchasedFromAuctionSite == 0 {
		e["has_previously_purchased_from_auction_site"] = "missing value"
	}
	if dirtyData.HasPreviouslyPurchasedFromFacebookMarketplace == 0 {
		e["has_previously_purchased_from_facebook_marketplace"] = "missing value"
	}
	if dirtyData.HasRegularlyAttendedComicConsOrCollectibleShows == 0 {
		e["has_regularly_attended_comic_cons_or_collectible_shows"] = "missing value"
	}

	if len(e) != 0 {
		return httperror.NewForBadRequest(&e)
	}
	return nil
}

func (impl *GatewayControllerImpl) RegisterCustomer(ctx context.Context, req *gateway_s.RegisterCustomerRequestIDO) error {
	// Defensive Code: For security purposes we need to remove all whitespaces from the email and lower the characters.
	req.Email = strings.ToLower(req.Email)
	req.Password = strings.ReplaceAll(req.Password, " ", "")

	// Perform our validation and return validation error on any issues detected.
	if err := impl.validateRegisterCustomerRequest(ctx, req); err != nil {
		impl.Logger.Warn("customer registration validation error",
			slog.Any("error", err))
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

		// Lookup the user and store in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByEmail(sessCtx, req.Email)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if u != nil {
			impl.Logger.Warn("user already exists validation error")
			return nil, httperror.NewForBadRequestWithSingleField("email", "email is not unique")
		}
		s, err := impl.StoreStorer.GetByID(sessCtx, req.StoreID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if s == nil {
			impl.Logger.Warn("store does not exists error")
			return nil, httperror.NewForBadRequestWithSingleField("store_id", fmt.Sprintf("does not exist for value: %s", req.StoreID.Hex()))
		}

		// Create our user.
		u, err = impl.createCustomerUserForRequest(sessCtx, req)
		if err != nil {
			return nil, err
		}

		// Create our store.
		u.StoreID = s.ID
		u.StoreName = s.Name
		u.StoreLevel = s.Level
		u.StoreTimezone = s.Timezone
		// u.StoreTimezone = req.Timezone
		u.ModifiedAt = time.Now()
		if err := impl.UserStorer.UpdateByID(sessCtx, u); err != nil {
			impl.Logger.Error("database update error", slog.Any("error", err))
			return nil, err
		}

		// Send our verification email.
		if err := impl.TemplatedEmailer.SendCustomerVerificationEmail(u.Email, u.EmailVerificationCode, u.FirstName); err != nil {
			impl.Logger.Error("failed sending verification email with error", slog.Any("err", err))
			return nil, err
		}

		// return nil, httperror.NewForBadRequestWithSingleField("message", "halted by programmer") // For debugging purposes only!

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

func (impl *GatewayControllerImpl) createCustomerUserForRequest(sessCtx mongo.SessionContext, req *gateway_s.RegisterCustomerRequestIDO) (*user_s.User, error) {
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
		ID:                                    userID,
		FirstName:                             req.FirstName,
		LastName:                              req.LastName,
		Name:                                  fmt.Sprintf("%s %s", req.FirstName, req.LastName),
		LexicalName:                           fmt.Sprintf("%s, %s", req.LastName, req.FirstName),
		Email:                                 req.Email,
		PasswordHash:                          passwordHash,
		PasswordHashAlgorithm:                 impl.Password.AlgorithmName(),
		Role:                                  user_s.UserRoleCustomer,
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
		EmailVerificationCode:                 impl.UUID.NewUUID(),
		EmailVerificationExpiry:               time.Now().Add(72 * time.Hour),
		Status:                                user_s.UserStatusActive,
		HasShippingAddress:                    req.HasShippingAddress,
		ShippingName:                          req.ShippingName,
		ShippingPhone:                         req.ShippingPhone,
		ShippingCountry:                       req.ShippingCountry,
		ShippingRegion:                        req.ShippingRegion,
		ShippingCity:                          req.ShippingCity,
		ShippingPostalCode:                    req.ShippingPostalCode,
		ShippingAddressLine1:                  req.ShippingAddressLine1,
		ShippingAddressLine2:                  req.ShippingAddressLine2,
		PaymentProcessorName:                  impl.PaymentProcessor.GetName(),
		PaymentProcessorCustomerID:            *paymentProcessorCustomerID,
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
	err = impl.UserStorer.Create(sessCtx, u)
	if err != nil {
		impl.Logger.Error("database create error", slog.Any("error", err))
		return nil, err
	}
	impl.Logger.Info("Customer user created.",
		slog.Any("_id", u.ID),
		slog.String("full_name", u.Name),
		slog.String("email", u.Email),
		slog.String("password_hash_algorithm", u.PasswordHashAlgorithm),
		slog.String("password_hash", u.PasswordHash))

	return u, nil
}
