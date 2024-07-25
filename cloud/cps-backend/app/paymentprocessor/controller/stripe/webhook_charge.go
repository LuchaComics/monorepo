package stripe

import (
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/stripe/stripe-go/v75"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	el_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
	r_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/receipt/datastore"
	up_s "github.com/LuchaComics/monorepo/cloud/cps-backend/app/userpurchase/datastore"
)

// webhookForChargeSucceeded function will handle Stripe's `charge.succeeded` webhook event in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForChargeSucceeded(sessCtx mongo.SessionContext, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForChargeSucceeded: starting...", slog.String("webhook", string(event.Type)))

	////
	//// Marshal & extract metadata.
	////

	var chrg stripe.Charge

	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &chrg); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}

	// DEVELOPERS NOTE: We do this to prevent duplicates.
	c.Kmutex.Lockf("%v", chrg.PaymentIntent.ID)
	defer func() {
		c.Kmutex.Unlockf("%v", chrg.PaymentIntent.ID)
	}()

	csID, err := primitive.ObjectIDFromHex(chrg.Metadata["ComicSubmissionID"])
	if err != nil {
		c.Logger.Error("converting object id from hex error",
			slog.Any("Error", err),
			slog.Any("Metadata Key", "ComicSubmissionID"),
			slog.Any("Metadata Value", chrg.Metadata["ComicSubmissionID"]))
		return err
	}
	if csID.IsZero() {
		c.Logger.Error("comicbook submission id failed to be extracted from metadata",
			slog.Any("Metadata Key", "ComicSubmissionID"),
			slog.Any("Metadata Value", chrg.Metadata["ComicSubmissionID"]))
		return errors.New("comicbook submission id failed to be extracted from metadata")
	}

	uID, err := primitive.ObjectIDFromHex(chrg.Metadata["UserID"])
	if err != nil {
		c.Logger.Error("converting object id from hex error",
			slog.Any("Error", err),
			slog.Any("Metadata Key", "UserID"),
			slog.Any("Metadata Value", chrg.Metadata["UserID"]))
		return err
	}
	if uID.IsZero() {
		c.Logger.Error("user id failed to be extracted from metadata",
			slog.Any("Metadata Key", "UserID"),
			slog.Any("Metadata Value", chrg.Metadata["UserID"]))
		return errors.New("user id failed to be extracted from metadata")
	}

	oID, err := primitive.ObjectIDFromHex(chrg.Metadata["OfferID"])
	if err != nil {
		c.Logger.Error("converting object id from hex error",
			slog.Any("Error", err),
			slog.Any("Metadata Key", "OfferID"),
			slog.Any("Metadata Value", chrg.Metadata["OfferID"]))
		return err
	}
	if oID.IsZero() {
		c.Logger.Error("offer id failed to be extracted from metadata",
			slog.Any("Metadata Key", "OfferID"),
			slog.Any("Metadata Value", chrg.Metadata["OfferID"]))
		return errors.New("offer id failed to be extracted from metadata")
	}

	////
	//// Get related records.
	////

	u, err := c.UserStorer.GetByID(sessCtx, uID)
	if err != nil {
		c.Logger.Error("get user by pp customer id error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	if u == nil {
		c.Logger.Error("customer does not exist error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return errors.New("customer does not exist")
	}

	o, err := c.OfferStorer.GetByID(sessCtx, oID)
	if err != nil {
		c.Logger.Error("get offer by id error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	if o == nil {
		c.Logger.Error("offer does not exist error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return errors.New("offer does not exist")
	}

	cs, err := c.ComicSubmissionStorer.GetByID(sessCtx, csID)
	if err != nil {
		c.Logger.Error("get comic submission by id error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	if cs == nil {
		c.Logger.Error("comic submission does not exist error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return errors.New("comic submission does not exist")
	}

	up, err := c.UserPurchaseStorer.GetByComicSubmissionID(sessCtx, csID)
	if err != nil {
		c.Logger.Error("get comic submission by id error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}

	////
	//// Create or update user purchase.
	////

	// Use the user's provided time zone or default to UTC.
	location, _ := time.LoadLocation("UTC")

	if up == nil {
		////
		//// Create user purchase receipt
		////

		up = &up_s.UserPurchase{
			StoreID:                     u.StoreID,
			StoreName:                   u.StoreName,
			StoreTimezone:               u.StoreTimezone,
			ID:                          primitive.NewObjectID(),
			UserID:                      u.ID,
			UserName:                    u.Name,
			UserLexicalName:             u.LexicalName,
			Status:                      r_s.StatusActive,
			CreatedAt:                   time.Now().In(location),
			ModifiedAt:                  time.Now().In(location),
			OfferID:                     o.ID,
			OfferName:                   o.Name,
			OfferDescription:            o.Description,
			OfferType:                   o.Type,
			OfferPrice:                  o.Price,
			OfferPriceCurrency:          o.PriceCurrency,
			OfferPayFrequency:           o.PayFrequency,
			OfferBusinessFunction:       o.BusinessFunction,
			OfferServiceType:            o.ServiceType,
			ComicSubmissionID:           cs.ID,
			ComicSubmissionSeriesTitle:  cs.SeriesTitle,
			ComicSubmissionIssueVol:     cs.IssueVol,
			ComicSubmissionIssueNo:      cs.IssueNo,
			PaymentProcessor:            r_s.PaymentProcessorStripe,
			PaymentProcessorReceiptID:   chrg.ID,
			PaymentProcessorReceiptURL:  chrg.ReceiptURL,
			PaymentProcessorPurchaseID:  chrg.PaymentIntent.ID,
			PaymentProcessorPurchasedAt: time.Now().In(location),
			AmountTotal:                 fromStripeFormat(chrg.Amount),
		}
		if err := c.UserPurchaseStorer.Create(sessCtx, up); err != nil {
			c.Logger.Error("create user purchase error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
			return err
		}
		c.Logger.Debug("created user purchase for comic book submission purchase", slog.String("webhook", string(event.Type)))
	} else {
		////
		//// Update user purchase receipt
		////

		up.StoreID = u.StoreID
		up.StoreName = u.StoreName
		up.StoreTimezone = u.StoreTimezone
		up.UserID = u.ID
		up.UserName = u.Name
		up.UserLexicalName = u.LexicalName
		up.Status = r_s.StatusActive
		up.ModifiedAt = time.Now().In(location)
		up.OfferID = o.ID
		up.OfferName = o.Name
		up.OfferDescription = o.Description
		up.OfferType = o.Type
		up.OfferPrice = o.Price
		up.OfferPriceCurrency = o.PriceCurrency
		up.OfferPayFrequency = o.PayFrequency
		up.OfferBusinessFunction = o.BusinessFunction
		up.OfferServiceType = o.ServiceType
		up.ComicSubmissionID = cs.ID
		up.ComicSubmissionSeriesTitle = cs.SeriesTitle
		up.ComicSubmissionIssueVol = cs.IssueVol
		up.ComicSubmissionIssueNo = cs.IssueNo
		up.PaymentProcessor = r_s.PaymentProcessorStripe
		up.PaymentProcessorReceiptID = chrg.ID
		up.PaymentProcessorReceiptURL = chrg.ReceiptURL
		up.PaymentProcessorPurchaseID = chrg.PaymentIntent.ID
		up.PaymentProcessorPurchasedAt = time.Now().In(location)
		up.AmountTotal = fromStripeFormat(chrg.Amount)
		if err := c.UserPurchaseStorer.UpdateByID(sessCtx, up); err != nil {
			c.Logger.Error("update user purchase error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
			return err
		}
		c.Logger.Debug("update user purchase for comic book submission purchase", slog.String("webhook", string(event.Type)))
	}

	////
	//// Update comics submission.
	////

	cs.PaymentProcessor = up.PaymentProcessor
	cs.PaymentProcessorReceiptID = chrg.ID
	cs.PaymentProcessorReceiptURL = chrg.ReceiptURL
	cs.PaymentProcessorPurchaseID = chrg.PaymentIntent.ID
	cs.AmountTotal = fromStripeFormat(chrg.Amount)
	// cs.PaymentProcessorPurchaseStatus = up.PaymentProcessorPurchaseStatus
	cs.PaymentProcessorPurchasedAt = time.Now().In(location)
	cs.PaymentProcessorPurchaseError = "" // Reset error.
	if err := c.ComicSubmissionStorer.UpdateByID(sessCtx, cs); err != nil {
		c.Logger.Error("update user purchase error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	c.Logger.Debug("update comic book submission", slog.String("webhook", string(event.Type)))

	////
	//// Mark the logevent as processed.
	////

	el.Status = el_d.StatusOK
	if err := c.EventLogStorer.UpdateByID(sessCtx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}

	c.Logger.Debug("webhookForChargeSucceeded: finished", slog.String("webhook", string(event.Type)))
	return nil
}
