package stripe

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	"github.com/stripe/stripe-go/v75"
	"go.mongodb.org/mongo-driver/mongo"

	el_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
	off_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
)

// webhookForPriceCreatedOrUpdated function will handle Stripe's `price.created` webhook event to create an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForPriceCreatedOrUpdated(sessCtx mongo.SessionContext, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForPriceCreatedOrUpdated: starting...", slog.String("webhook", string(event.Type)))

	////
	//// Get the price and product.
	////

	var price stripe.Price

	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &price); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	el.Status = el_d.StatusOK

	off, err := c.OfferStorer.GetByStripeProductID(sessCtx, price.Product.ID)
	if err != nil {
		c.Logger.Error("get offer error", slog.Any("price.Product", price.Product), slog.String("webhook", string(event.Type)))
		return err
	}
	if off == nil {
		c.Logger.Error("offer does not exist error", slog.String("price.Product.ID", price.Product.ID), slog.String("webhook", string(event.Type)))
		return fmt.Errorf("offer does not exist for `product_id` value of %v", price.Product.ID)
	}

	////
	//// Update the offer.
	////

	off.StripePriceID = price.ID
	off.Price = fromStripeFormat(price.UnitAmount)
	off.PriceCurrency = strings.ToUpper(string(price.Currency))
	off.PayFrequency = off_d.PayFrequencyOneTime

	if err := c.OfferStorer.UpdateByID(sessCtx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	c.Logger.Debug("Created/updated product price", slog.String("webhook", string(event.Type)))

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(sessCtx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForPriceCreatedOrUpdated: finished", slog.String("webhook", string(event.Type)))
	return nil
}
