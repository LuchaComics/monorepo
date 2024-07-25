package stripe

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/stripe/stripe-go/v75"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	el_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
	off_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
)

// webhookForProductCreated function will handle Stripe's `product.created` webhook event to create an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForProductCreated(sessCtx mongo.SessionContext, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForProductCreated: starting...", slog.String("webhook", string(event.Type)))

	// The following fields must be filled out:
	// - name
	// - description
	// - images (only one)
	// - metadata: { "StoreID": "648763d3f6fbead15f5bd4d2" },
	// - Recurring
	// - Monthly billing period
	// - price

	////
	//// Get the product.
	////

	var product stripe.Product
	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &product); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	el.Status = el_d.StatusOK

	////
	//// Create the offer.
	////

	off := &off_d.Offer{
		ID:                   primitive.NewObjectID(),
		Name:                 product.Name,
		Description:          product.Description,
		Status:               off_d.StatusPending,
		PaymentProcessorName: "Stripe, Inc.",
		StripeProductID:      product.ID,
		CreatedAt:            time.Now(),
		ModifiedAt:           time.Now(),
		BusinessFunction:     off_d.BusinessFunctionUnspecified,
	}

	if len(product.Images) > 0 {
		off.StripeImageURL = product.Images[0]
	}

	if err := c.OfferStorer.Create(sessCtx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	c.Logger.Debug("webhookForProductCreated: successful processed", slog.String("webhook", string(event.Type)))

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(sessCtx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForProductCreated: finished", slog.String("webhook", string(event.Type)))
	return nil
}

// webhookForProductUpdated function will handle Stripe's `product.updated` webhook event to update an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForProductUpdated(sessCtx mongo.SessionContext, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForProductUpdated: starting...", slog.String("webhook", string(event.Type)))

	////
	//// Get the product.
	////

	var product stripe.Product
	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &product); err != nil {
		c.Logger.Error("unmarshalling error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	el.Status = el_d.StatusOK

	off, err := c.OfferStorer.GetByStripeProductID(sessCtx, product.ID)
	if err != nil {
		c.Logger.Error("get offer error", slog.Any("plan.Product", product.ID), slog.String("webhook", string(event.Type)))
		return err
	}
	if off == nil {
		c.Logger.Error("offer does not exist error", slog.String("price.Product.ID", product.ID), slog.String("webhook", string(event.Type)))
		return fmt.Errorf("offer does not exist for `product_id` value of %v", product.ID)
	}

	////
	//// Lookup price and update pricing
	////

	// DEVELOPERS NOTE: This code is somewhat not-needed since we have another
	// webhook to handle price changes. However this code has proven really
	// useful in developer mode when making manual modifications. Keep this code
	// in here for this reason.

	if product.DefaultPrice != nil {
		price, err := c.PaymentProcessor.GetPrice(product.DefaultPrice.ID)
		if err != nil {
			c.Logger.Error("get price error", slog.Any("plan.Product", product.ID), slog.String("webhook", string(event.Type)))
			return err
		}
		if price != nil {
			off.StripePriceID = price.ID
			off.Price = fromStripeFormat(price.UnitAmount)
			off.PriceCurrency = strings.ToUpper(string(price.Currency))
			off.PayFrequency = off_d.PayFrequencyOneTime
			c.Logger.Debug("updated pricing of product", slog.String("webhook", string(event.Type)))
		}
	}

	////
	//// Update the offer.
	////

	off.Name = product.Name
	off.Description = product.Description
	if len(product.Images) > 0 {
		off.StripeImageURL = product.Images[0]
	}
	off.ModifiedAt = time.Now()
	if err := c.OfferStorer.UpdateByID(sessCtx, off); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}
	c.Logger.Debug("updated product")

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(sessCtx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err), slog.String("webhook", string(event.Type)))
		return err
	}

	c.Logger.Debug("webhookForProductUpdated: finished", slog.String("webhook", string(event.Type)))
	return nil
}
