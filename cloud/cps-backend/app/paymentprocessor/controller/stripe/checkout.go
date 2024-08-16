package stripe

import (
	"context"
	"errors"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/LuchaComics/monorepo/cloud/cps-backend/config/constants"
)

func (impl *StripePaymentProcessorControllerImpl) CreateStripeCheckoutSessionURLForComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (string, error) {
	////
	//// Start the transaction.
	////

	session, err := impl.DbClient.StartSession()
	if err != nil {
		impl.Logger.Error("start session error",
			slog.Any("error", err))
		return "", err
	}
	defer session.EndSession(ctx)

	// Define a transaction function with a series of operations
	transactionFunc := func(sessCtx mongo.SessionContext) (interface{}, error) {
		// Extract from our session the following data.
		// This user is the logged in retailer admin as they are the only ones
		// whom can purchase.
		userID := sessCtx.Value(constants.SessionUserID).(primitive.ObjectID)

		// STEP 1: Lookup the submission in our database, else return a `400 Bad Request` error.
		cs, err := impl.ComicSubmissionStorer.GetByID(sessCtx, comicSubmissionID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return "", err
		}
		if cs == nil {
			impl.Logger.Warn("comic submission does not exist validation error")
			return "", errors.New("comic submission id does not exist")
		}

		// STEP 2: Lookup the offer in our database, else return a `400 Bad Request` error.
		o, err := impl.OfferStorer.GetByServiceType(sessCtx, cs.ServiceType)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return "", err
		}
		if o == nil {
			impl.Logger.Warn("offer does not exist validation error")
			return "", errors.New("offer does not exist")
		}

		// STEP 3: Lookup the user in our database, else return a `400 Bad Request` error.
		u, err := impl.UserStorer.GetByID(sessCtx, userID)
		if err != nil {
			impl.Logger.Error("database error", slog.Any("err", err))
			return "", err
		}
		if u == nil {
			impl.Logger.Warn("user does not exist validation error")
			return "", errors.New("user does not exist")
		}

		// Defensive code: Prevent executing this function if different processor.
		if o.PaymentProcessorName != "Stripe, Inc." {
			impl.Logger.Warn("not stripe payment processor assigned to offer.",
				slog.String("offer_payment_processor_name", o.PaymentProcessorName),
				slog.String("user_payment_processor_name", u.PaymentProcessorName))
			return "", errors.New("offer is using payment processor which is not supported")
		}
		if u.FirstName == "Root" {
			impl.Logger.Warn("not stripe payment processor assigned to user.",
				slog.String("first_name", u.FirstName))
			return "", errors.New("system administrator are not allowed to purchase")
		}
		if u.PaymentProcessorName != "Stripe, Inc." {
			impl.Logger.Warn("not stripe payment processor assigned to user.",
				slog.String("user_id", userID.Hex()),
				slog.String("offer_payment_processor_name", o.PaymentProcessorName),
				slog.String("user_payment_processor_name", u.PaymentProcessorName))
			return "", errors.New("user is using payment processor which is not supported")
		}

		// Defensive code: Prevent executing if no `customer id` exist from stripe.
		if u.PaymentProcessorCustomerID == "" {
			impl.Logger.Warn("not stripe payment processor customer id assigned to user.")
			return "", errors.New("user has no customer id set by payment processor")
		}

		// Defensive code: Prevent executing if `product id` have not been created.
		if o.StripePriceID == "" {
			impl.Logger.Warn("this product is not ready")
			return "", errors.New("this product is not ready")
		}

		hasShippingAddress := u.ShippingCity != "" || u.ShippingCountry != "" || u.ShippingAddressLine1 != ""

		impl.Logger.Debug("creating stripe checkout session",
			slog.String("priceID", o.StripePriceID),
			slog.Any("offerID", o.ID),
			slog.Any("hasShippingAddress", hasShippingAddress))

		// DEVELOPERS NOTE:
		// THIS IS HOW WE ATTACH OUR METADATA TO OUR STRIPE SUBMISSION. THIS ALLOWS
		// USE TO DEREFERENCE THE METADATA LATER IN WEBHOOKS. THIS IS IMPORTANT TO
		// HOW OUR APP WORKS.
		metadata := make(map[string]string)
		metadata["ComicSubmissionID"] = comicSubmissionID.Hex()
		metadata["UserID"] = u.ID.Hex()
		metadata["OfferID"] = o.ID.Hex()
		metadata["Type"] = "Comic Book Submission"

		// DEVELOPERS NOTE:
		// THIS IS HOW WE SUBMIT OUR APPS CONFIGURAITON FOR THE PRODUCT AND
		// STRIPE WILL GENERATE A CHECKOUT SESSION URL TO USE IN OUR APP.
		var redirectURL string
		redirectURL, err = impl.PaymentProcessor.CreateOneTimeCheckoutSessionURL( // TODO: FIX TO SUPPORT URL WITH COMIC SUBMISSION IS DONE.
			impl.Emailer.GetFrontendDomainName(),
			"/submissions/comics/add/confirmation?submission_id="+comicSubmissionID.Hex(),              // Accepted URL
			"/submissions/comics/add/checkout?submission_id="+comicSubmissionID.Hex()+"&canceled=true", // Cancelled URL
			u.PaymentProcessorCustomerID,
			o.StripePriceID,
			metadata,
			hasShippingAddress,
		)
		if err != nil {
			return "", err
		}
		impl.Logger.Debug("stripe checkout session ready", slog.String("redirectURL", redirectURL))
		return redirectURL, nil
	}

	// Start a transaction
	res, err := session.WithTransaction(ctx, transactionFunc)
	if err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return "", err
	}

	return res.(string), nil
}

type CompleteStripeCheckoutSessionResponse struct {
	Name          string  `json:"name"`
	Description   string  `bson:"description" json:"description"`
	Price         float64 `bson:"price" json:"price"`
	PriceCurrency string  `bson:"price_currency" json:"price_currency"`
	PayFrequency  int8    `bson:"pay_frequency" json:"pay_frequency"`
	SessionID     string  `json:"session_id"`
	PaymentStatus string  `json:"payment_status"`
	Status        string  `json:"status"`
}
