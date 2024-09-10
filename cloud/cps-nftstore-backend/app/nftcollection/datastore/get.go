package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl NFTCollectionStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*NFTCollection, error) {
	filter := bson.M{"_id": id}

	var result NFTCollection
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

func (impl NFTCollectionStorerImpl) GetByName(ctx context.Context, name string) (*NFTCollection, error) {
	filter := bson.M{"name": name}

	var result NFTCollection
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

func (impl NFTCollectionStorerImpl) GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*NFTCollection, error) {
	filter := bson.M{"payment_processor_receipt_id": paymentProcessorPurchaseID}

	var result NFTCollection
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

func (impl NFTCollectionStorerImpl) GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*NFTCollection, error) {
	filter := bson.M{"comic_submission_id": comicSubmissionID}

	var result NFTCollection
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
