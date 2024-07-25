package controller

import (
	"context"
	"time"

	c_ds "github.com/LuchaComics/monorepo/cloud/cps-backend/app/comicsub/datastore"
	o_ds "github.com/LuchaComics/monorepo/cloud/cps-backend/app/offer/datastore"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl *OfferControllerImpl) createDefaults(ctx context.Context) error {
	impl.Logger.Debug("offer createDefaults started...")
	// DEVELOPERS NOTE: ALL OF THESE ARE PURPOPSEFULLY HARD-CODE AND ARE TO BE DELETED UPON RELEASE.

	// --- OFFER #1 --- //

	o1ID, _ := primitive.ObjectIDFromHex("65132fcf68e145414eba202e")
	o1, _ := impl.GetByID(ctx, o1ID)
	if o1 == nil || true {
		impl.Logger.Debug("upserting offer `Service - C Capsule`...")
		o := &o_ds.Offer{
			ID:                   o1ID,
			Name:                 "Service - C Capsule",
			Description:          "",
			Price:                100,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyOneTime,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionGrantsComicbookSubmissionServices,
			ServiceType:          c_ds.ServiceTypeCPSCapsule,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			PaymentProcessorName: "Stripe, Inc.",
			StripeProductID:      "prod_Oi642jRR3R0nIf",
			StripePriceID:        "price_1NufsLJ5szlo8iRpJcI5gdmT",
		}
		if err := impl.OfferStorer.Upsert(ctx, o); err != nil {
			return err
		}
	}

	// --- OFFER #2 --- //

	o2ID, _ := primitive.ObjectIDFromHex("6513409753766c5b72987fd6")
	o2, _ := impl.GetByID(ctx, o2ID)
	if o2 == nil || true {
		impl.Logger.Debug("upserting offer `Pedigree`...")
		o := &o_ds.Offer{
			ID:                   o2ID,
			Name:                 "Pedigree",
			Description:          "",
			Price:                50,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyOneTime,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionGrantsComicbookSubmissionServices,
			ServiceType:          c_ds.ServiceTypePedigree,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			PaymentProcessorName: "Stripe, Inc.",
			StripeProductID:      "prod_Oi8F7qg5n8qkrk",
			StripePriceID:        "price_1NuhylJ5szlo8iRpM4xMflgI",
		}
		if err := impl.OfferStorer.Upsert(ctx, o); err != nil {
			return err
		}
	}

	// --- OFFER #3 --- //

	o3ID, _ := primitive.ObjectIDFromHex("651343ccb050f47504eca649")
	o3, _ := impl.GetByID(ctx, o3ID)
	if o3 == nil || true {
		impl.Logger.Debug("upserting offer `C Capsule Signature`...")
		o := &o_ds.Offer{
			ID:                   o3ID,
			Name:                 "C Capsule Signature",
			Description:          "",
			Price:                150,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyOneTime,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionGrantsComicbookSubmissionServices,
			ServiceType:          c_ds.ServiceTypeCPSCapsuleSignatureCollection,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			PaymentProcessorName: "Stripe, Inc.",
			StripeProductID:      "prod_Oi8S8vgeoLEeT1",
			StripePriceID:        "price_1NuiC0J5szlo8iRppsD7PmS3",
		}
		if err := impl.OfferStorer.Upsert(ctx, o); err != nil {
			return err
		}
	}

	// --- OFFER #4 --- //

	o4ID, _ := primitive.ObjectIDFromHex("651344e515213afe17b2d4ba")
	o4, _ := impl.GetByID(ctx, o4ID)
	if o4 == nil || true {
		impl.Logger.Debug("upserting offer `C Capsule - Mint Indie Gem`...")
		o := &o_ds.Offer{
			ID:                   o4ID,
			Name:                 "C Capsule - Mint Indie Gem",
			Description:          "",
			Price:                500,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyOneTime,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionGrantsComicbookSubmissionServices,
			ServiceType:          c_ds.ServiceTypeCPSCapsuleIndieMintGem,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			PaymentProcessorName: "Stripe, Inc.",
			StripeProductID:      "prod_Oi8Xo4NRaaSwlE",
			StripePriceID:        "price_1NuiGXJ5szlo8iRpcl6oIxAi",
		}
		if err := impl.OfferStorer.Upsert(ctx, o); err != nil {
			return err
		}
	}

	// --- OFFER #5 --- //

	o5ID, _ := primitive.ObjectIDFromHex("6513459ccc447d773c08d9d2")
	o5, _ := impl.GetByID(ctx, o5ID)
	if o5 == nil || true {
		impl.Logger.Debug("upserting offer `C Capsule - U Grade`...")
		o := &o_ds.Offer{
			ID:                   o5ID,
			Name:                 "C Capsule - U Grade",
			Description:          "",
			Price:                125,
			PriceCurrency:        "CAD",
			PayFrequency:         o_ds.PayFrequencyOneTime,
			Status:               o_ds.StatusActive,
			Type:                 o_ds.OfferTypeService,
			BusinessFunction:     o_ds.BusinessFunctionGrantsComicbookSubmissionServices,
			ServiceType:          c_ds.ServiceTypeCPSCapsuleYouGrade,
			CreatedAt:            time.Now(),
			ModifiedAt:           time.Now(),
			PaymentProcessorName: "Stripe, Inc.",
			StripeProductID:      "prod_Oi8aPnlSuArGtb",
			StripePriceID:        "price_1NuiJUJ5szlo8iRpZtAGBfuJ",
		}
		if err := impl.OfferStorer.Upsert(ctx, o); err != nil {
			return err
		}
	}

	impl.Logger.Debug("offer createDefaults finished")
	return nil
}
