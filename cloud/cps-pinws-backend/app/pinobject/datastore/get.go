package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

func (impl PinObjectStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*PinObject, error) {
	filter := bson.M{"_id": id}

	var result PinObject
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

func (impl PinObjectStorerImpl) GetByRequestID(ctx context.Context, rid primitive.ObjectID) (*PinObject, error) {
	filter := bson.M{"requestid": rid}

	var result PinObject
	err := impl.Collection.FindOne(ctx, filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// This error means your query did not match any documents.
			return nil, nil
		}
		impl.Logger.Error("database get by request id error", slog.Any("error", err))
		return nil, err
	}
	return &result, nil
}
