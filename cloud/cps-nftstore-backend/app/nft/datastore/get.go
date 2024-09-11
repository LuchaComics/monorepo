package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl NFTStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*NFT, error) {
	filter := bson.M{"_id": id}

	var result NFT
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

func (impl NFTStorerImpl) GetByName(ctx context.Context, name string) (*NFT, error) {
	filter := bson.M{"name": name}

	var result NFT
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

func (impl NFTStorerImpl) GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*NFT, error) {
	filter := bson.M{"payment_processor_receipt_id": paymentProcessorPurchaseID}

	var result NFT
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

func (impl NFTStorerImpl) GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*NFT, error) {
	filter := bson.M{"comic_submission_id": comicSubmissionID}

	var result NFT
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
