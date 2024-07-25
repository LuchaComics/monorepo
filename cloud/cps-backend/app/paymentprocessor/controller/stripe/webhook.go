package stripe

import (
	"context"
	"time"

	"log/slog"

	"github.com/stripe/stripe-go/v75"
	"github.com/stripe/stripe-go/v75/webhook"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	el_d "github.com/LuchaComics/monorepo/cloud/cps-backend/app/eventlog/datastore"
)

func (impl *StripePaymentProcessorControllerImpl) Webhook(ctx context.Context, header string, b []byte) error {
	// TECHDEBT: Find a way to replace parts of this code into adapter.

	// DEVELOPERS NOTE:
	// HOW DO WE MANUALLY RESEND AN EVENT THAT WAS PREVIOUSLY CALLED?
	// STEP 1:
	// Perform the action you want in the app: Ex: purchase.
	//
	// STEP 2:
	// Go to: https://dashboard.stripe.com/test/events
	//
	// STEP 3:
	// Lookup the event ID, for example: `evt_1NfoHAC1dNpgYbqFJfSeeCNK`.
	//
	// STEP 4:
	// In a new terminal window, run the code to resend (note: https://stripe.com/docs/cli/events/resend):
	// $ stripe events resend evt_1NfoHAC1dNpgYbqFJfSeeCNK

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
		event, err := webhook.ConstructEvent(b, header, impl.PaymentProcessor.GetWebhookSecretKey())
		if err != nil {
			impl.Logger.Error("construct event error",
				slog.Any("err", err),
				slog.Any("WebhookSecretKey", impl.PaymentProcessor.GetWebhookSecretKey()))
			return nil, err
		}

		// Log in our system the webhook event and process in our system this event.
		eventlog, err := impl.logWebhookEvent(sessCtx, event)
		if err != nil {
			impl.Logger.Error("failed logging stripe webhook event", slog.Any("event_type", event.Type))
			return nil, err
		}

		impl.Logger.Debug("stripe webhook executing...", slog.Any("event_type", event.Type))

		// DEVELOPERS NOTE: The full list can be found via https://stripe.com/docs/api/events/types
		switch eventlog.SecondaryType {
		case "product.created":
			return nil, impl.webhookForProductCreated(sessCtx, event, eventlog)
		case "product.updated":
			return nil, impl.webhookForProductUpdated(sessCtx, event, eventlog)
		case "plan.created", "plan.updated":
			return nil, impl.webhookForPlanCreatedOrUpdated(sessCtx, event, eventlog)
		case "price.created", "price.updated":
			return nil, impl.webhookForPriceCreatedOrUpdated(sessCtx, event, eventlog)
		case "charge.succeeded":
			return nil, impl.webhookForChargeSucceeded(sessCtx, event, eventlog)
		case "payment_intent.created":
			impl.Logger.Warn("skip processing `payment_intent.created` stripe event ")
			return nil, nil
		case "payment_intent.succeeded":
			return nil, impl.webhookForPaymentIntentSucceeded(sessCtx, event, eventlog)
		case "checkout.session.completed":
			return nil, impl.webhookForCheckoutSessionCompleted(sessCtx, event, eventlog)
		default:
			impl.Logger.Warn("skip processing stripe event", slog.Any("eventType", event.Type))
			return nil, nil
		}
	}

	// Start a transaction
	if _, err := session.WithTransaction(ctx, transactionFunc); err != nil {
		impl.Logger.Error("session failed error",
			slog.Any("error", err))
		return err
	}

	return nil
}

func (impl *StripePaymentProcessorControllerImpl) logWebhookEvent(sessCtx mongo.SessionContext, event stripe.Event) (*el_d.EventLog, error) {
	impl.Logger.Debug("logging stripe webhook event...")

	// DEVELOPERS NOTE:
	// We will take advantage of Golang's new "generics" functionality and store
	// the event data into our database. The mongodb library will understand
	// how to handle the generiimpl.

	eventlog := &el_d.EventLog{
		PrimaryType:   el_d.PrimaryTypeStripeWebhookEvent,
		SecondaryType: string(event.Type),
		CreatedAt:     time.Now(),
		Content:       event.Data.Object, // Store the event payload, not metadata.
		Status:        el_d.StatusPending,
		ExternalID:    event.ID,
		ID:            primitive.NewObjectID(),
	}
	if err := impl.EventLogStorer.Create(sessCtx, eventlog); err != nil {
		impl.Logger.Error("marshalling create error", slog.Any("err", err))
		return nil, err
	}
	impl.Logger.Debug("logged stripe webhook event")
	return eventlog, nil
}
