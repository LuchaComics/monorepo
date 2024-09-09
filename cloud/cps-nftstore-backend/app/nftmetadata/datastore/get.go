package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl NFTMetadataStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*NFTMetadata, error) {
	filter := bson.M{"_id": id}

	var result NFTMetadata
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

func (impl NFTMetadataStorerImpl) GetByName(ctx context.Context, name string) (*NFTMetadata, error) {
	filter := bson.M{"name": name}

	var result NFTMetadata
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

func (impl NFTMetadataStorerImpl) GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*NFTMetadata, error) {
	filter := bson.M{"payment_processor_receipt_id": paymentProcessorPurchaseID}

	var result NFTMetadata
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

func (impl NFTMetadataStorerImpl) GetByComicSubmissionID(ctx context.Context, comicSubmissionID primitive.ObjectID) (*NFTMetadata, error) {
	filter := bson.M{"comic_submission_id": comicSubmissionID}

	var result NFTMetadata
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
