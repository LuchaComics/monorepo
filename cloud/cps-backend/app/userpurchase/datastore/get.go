package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl UserPurchaseStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*UserPurchase, error) {
	filter := bson.M{"_id": id}

	var result UserPurchase
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserPurchaseStorerImpl) GetByName(ctx context.Context, name string) (*UserPurchase, error) {
	filter := bson.M{"name": name}

	var result UserPurchase
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserPurchaseStorerImpl) GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*UserPurchase, error) {
	filter := bson.M{"payment_processor_receipt_id": paymentProcessorPurchaseID}

	var result UserPurchase
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by `payment_processor_receipt_id` error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}

func (impl UserPurchaseStorerImpl) GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*UserPurchase, error) {
	filter := bson.M{"comic_submission_id": comicSubmissionID}

	var result UserPurchase
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by `comic_submission_id` error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
