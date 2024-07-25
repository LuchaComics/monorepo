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

// webhookForPlanCreatedOrUpdated function will handle Stripe's `plan.created` webhook event to create an offer in our system.
func (c *StripePaymentProcessorControllerImpl) webhookForPlanCreatedOrUpdated(sessCtx mongo.SessionContext, event stripe.Event, el *el_d.EventLog) error {
	c.Logger.Debug("webhookForPlanCreatedOrUpdated: starting...",
		slog.String("webhook", string(event.Type)))

	////
	//// Get the plan and product.
	////

	var plan stripe.Plan

	// Successfully cast to []byte
	if err := json.Unmarshal(event.Data.Raw, &plan); err != nil {
		c.Logger.Error("unmarshalling error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}
	el.Status = el_d.StatusOK

	off, err := c.OfferStorer.GetByStripeProductID(sessCtx, plan.Product.ID)
	if err != nil {
		c.Logger.Error("get offer error",
			slog.Any("plan.Product", plan.Product),
			slog.String("webhook", string(event.Type)))
		return err
	}
	if off == nil {
		c.Logger.Error("offer does not exist error",
			slog.String("price.Product.ID", plan.Product.ID),
			slog.String("webhook", string(event.Type)))
		return fmt.Errorf("offer does not exist for `product_id` value of %v", plan.Product.ID)
	}

	////
	//// Update the offer.
	////

	off.Price = fromStripeFormat(plan.Amount)
	off.PriceCurrency = strings.ToUpper(string(plan.Currency))
	switch plan.Interval {
	case stripe.PlanIntervalDay:
		off.PayFrequency = off_d.PayFrequencyDay
		break
	case stripe.PlanIntervalMonth:
		off.PayFrequency = off_d.PayFrequencyMonthly
		break
	case stripe.PlanIntervalWeek:
		off.PayFrequency = off_d.PayFrequencyWeek
		break
	case stripe.PlanIntervalYear:
		off.PayFrequency = off_d.PayFrequencyAnnual
		break
	default:
		off.PayFrequency = off_d.PayFrequencyOneTime
		break
	}

	if plan.Interval != "" {
		off.IsSubscription = true
	}

	if err := c.OfferStorer.UpdateByID(sessCtx, off); err != nil {
		c.Logger.Error("create offer error",
			slog.Any("err", err),
			slog.String("webhook", string(event.Type)))
		return err
	}
	c.Logger.Debug("webhookForPlanCreatedOrUpdated: successful processed",
		slog.String("webhook", string(event.Type)))

	////
	//// Mark the logevent as processed.
	////

	if err := c.EventLogStorer.UpdateByID(sessCtx, el); err != nil {
		c.Logger.Error("create offer error", slog.Any("err", err))
		return err
	}

	c.Logger.Debug("webhookForPlanCreatedOrUpdated: finished",
		slog.String("webhook", string(event.Type)))
	return nil
}
