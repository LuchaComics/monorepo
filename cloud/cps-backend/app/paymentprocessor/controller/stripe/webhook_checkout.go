package stripe

import (
	"encoding/json"
	"errors"
	"log/slog"
	"strings"

	"github.com/stripe/stripe-go/v75"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	el_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
)

// webhookForCheckoutSessionCompleted function will handle Stripe's `checkout.session.completed` webhook event in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForCheckoutSessionCompleted(sessCtx mongo.SessionContext, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForCheckoutSessionCompleted: starting...", slog.String("webhook", string(event.Type)))

	////
	//// Marshal
	////

	var session stripe.CheckoutSession

	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}

	// DEVELOPERS NOTE: We do this to prevent duplicates.
	c.Kmutex.Lockf("%v", session.PaymentIntent.ID)
	defer func() {
		c.Kmutex.Unlockf("%v", session.PaymentIntent.ID)
	}()

	// Extract the comic submission id.
	// For example: https://cpsapp.ca/submissions/comics/add/confirmation?submission_id=652a094e09db53b4991d8df5&session_id={CHECKOUT_SESSION_ID}
	// Goes into: [https:  cpsapp.ca submissions comics add 652a094e09db53b4991d8df5 confirmation?session_id={CHECKOUT_SESSION_ID}]
	//
	arr := strings.Split(session.SuccessURL, "/")
	csIDStr := arr[6]
	c.Logger.Error("extracted comic book submission",
		slog.Any("id", csIDStr),
		slog.String("webhook", string(event.Type)))
	csID, err := primitive.ObjectIDFromHex(csIDStr)
	if err != nil {
		c.Logger.Error("converting object id from hex error",
			slog.Any("Error", err),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if csID.IsZero() {
		c.Logger.Error("comicbook submission id failed to be extracted from metadata", slog.String("webhook", string(event.Type)))
		return errors.New("comicbook submission id failed to be extracted from metadata")
	}

	////
	//// Get comic submission
	////

	cs, err := c.ComicSubmissionStorer.GetByID(sessCtx, csID)
	if err != nil {
		c.Logger.Error("get cs by id error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if cs == nil {
		c.Logger.Error("cs does not exist error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return errors.New("offer does not exist")
	}

	cs.AmountSubtotal = fromStripeFormat(session.AmountSubtotal)
	cs.AmountTax = fromStripeFormat(session.TotalDetails.AmountTax)
	cs.AmountTotal = fromStripeFormat(session.AmountTotal)
	if err := c.ComicSubmissionStorer.UpdateByID(sessCtx, cs); err != nil {
		c.Logger.Error("updated comic submission error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	c.Logger.Debug("updated comic submission", slog.String("webhook", string(event.Type)))

	////
	//// Update
	////

	up, err := c.UserPurchaseStorer.GetByComicSubmissionID(sessCtx, csID)
	if err != nil {
		c.Logger.Error("get comic submission by id error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	// if up == nil {
	// 	c.Logger.Error("user purchase does not exist error",
	// 		slog.Any("err", err),
	// 		slog.String("webhook", string(event.Type)))
	// 	return errors.New("user purchase does not exist")
	// }
	if up != nil {
		up.AmountSubtotal = fromStripeFormat(session.AmountSubtotal)
		up.AmountTax = fromStripeFormat(session.TotalDetails.AmountTax)
		up.AmountTotal = fromStripeFormat(session.AmountTotal)
		if err := c.UserPurchaseStorer.UpdateByID(sessCtx, up); err != nil {
			c.Logger.Error("updated user purchases error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
			return err
		}
		c.Logger.Debug("updated user purchase", slog.String("webhook", string(event.Type)))
	}

	////
	//// Mark the logevent as processed.
	////

	el.Status = el_d.StatusOK
	if err := c.EventLogStorer.UpdateByID(sessCtx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}

	c.Logger.Debug("webhookForCheckoutSessionCompleted: finished", slog.String("webhook", string(event.Type)))
	return nil
}
