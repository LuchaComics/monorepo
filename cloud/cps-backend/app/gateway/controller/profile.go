package controller

import (
	"context"
	"fmt"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	user_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/user/datastore"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
	"github.com/LuchaComics/monorepo/cloud/cps-backend/utils/httperror"
)

func (impl *GatewayControllerImpl) Profile(ctx context.Context) (*user_s.User, error) {
	// Extract from our session the following data.
	userID := ctx.Value(constants.SessionUserID).(primitive.ObjectID)

	// Lookup the user in our database, else return a `400 Bad Request` error.
	u, err := impl.UserStorer.GetByID(ctx, userID)
	if err != nil {
		impl.Logger.Error("database error", slog.Any("err", err))
		return nil, err
	}
	if u == nil {
		impl.Logger.Warn("user does not exist validation error")
		return nil, httperror.NewForBadRequestWithSingleField("id", "does not exist")
	}
	return u, nil
}

func (impl *GatewayControllerImpl) ProfileUpdate(ctx context.Context, nu *user_s.User) error {
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
		// Extract from our session the following data.
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

		// Lookup the user in our database, else return a `400 Bad Request` error.
		ou, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return nil, err
		}
		if ou == nil {
			impl.Logger.Warn("user does not exist validation error")
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
		ou.AgreePromotionsEmail = nu.AgreePromotionsEmail
		ou.HasShippingAddress = nu.HasShippingAddress
		ou.ShippingName = nu.ShippingName
		ou.ShippingPhone = nu.ShippingPhone
		ou.ShippingCountry = nu.ShippingCountry
		ou.ShippingRegion = nu.ShippingRegion
		ou.ShippingCity = nu.ShippingCity
		ou.ShippingPostalCode = nu.ShippingPostalCode
		ou.ShippingAddressLine1 = nu.ShippingAddressLine1
		ou.ShippingAddressLine2 = nu.ShippingAddressLine2

		if err := impl.UserStorer.UpdateByID(sessCtx, ou); err != nil {
			impl.Logger.Error("user update by id error", slog.Any("error", err))
			return nil, err
		}

		if ou.Role != user_s.UserRoleRoot {
			if "Stripe, Inc." == impl.PaymentProcessor.GetName() {
				err = impl.PaymentProcessor.UpdateCustomer(
					ou.PaymentProcessorCustomerID,
					fmt.Sprintf("%s %s", ou.FirstName, ou.LastName),
					ou.Email,
					"", // description...
					fmt.Sprintf("%s %s Shipping Address", ou.FirstName, ou.LastName),
					ou.Phone,
					ou.ShippingCity, ou.ShippingCountry, ou.ShippingAddressLine1, ou.ShippingAddressLine2, ou.ShippingPostalCode, ou.ShippingRegion, // Shipping
					ou.City, ou.Country, ou.AddressLine1, ou.AddressLine2, ou.PostalCode, ou.Region, // Billing
				)
				if err != nil {
					impl.Logger.Error("updated customer from payment processor error", slog.Any("error", err))
					return nil, err
				}
				impl.Logger.Debug("updated customer in stripe",
					slog.Any("ou.Region", ou.Region),
				)
			}
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
