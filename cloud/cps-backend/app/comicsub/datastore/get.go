package datastore

import (
	"context"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl ComicSubmissionStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*ComicSubmission, error) {
	filter := bson.M{"_id": id}

	var result ComicSubmission
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

func (impl ComicSubmissionStorerImpl) GetByCPSRN(ctx context.Context, cpsrn string) (*ComicSubmission, error) {
	filter := bson.M{"cpsrn": cpsrn}

	var result ComicSubmission
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

func (impl ComicSubmissionStorerImpl) GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*ComicSubmission, error) {
	filter := bson.M{"payment_processor_purchase_id": paymentProcessorPurchaseID}

	var result ComicSubmission
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
