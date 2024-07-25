package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl ReceiptStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*Receipt, error) {
	filter := bson.M{"_id": id}

	var result Receipt
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

func (impl ReceiptStorerImpl) GetByName(ctx context.Context, name string) (*Receipt, error) {
	filter := bson.M{"name": name}

	var result Receipt
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

func (impl ReceiptStorerImpl) GetByPaymentProcessorReceiptID(ctx context.Context, paymentProcessorReceiptID string) (*Receipt, error) {
	filter := bson.M{"payment_processor_receipt_id": paymentProcessorReceiptID}

	var result Receipt
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

func (impl ReceiptStorerImpl) GetByPaymentProcessorPurchaseID(ctx context.Context, paymentProcessorPurchaseID string) (*Receipt, error) {
	filter := bson.M{"payment_processor_purchase_id": paymentProcessorPurchaseID}

	var result Receipt
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
