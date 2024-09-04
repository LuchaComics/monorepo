package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (impl PinObjectStorerImpl) GetByCID(ctx context.Context, cid string) (*PinObject, error) {
	filter := bson.M{"cid": cid}

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

func (impl PinObjectStorerImpl) GetAllCIDs(ctx context.Context) ([]string, error) {
	// Define an empty filter to retrieve all documents
	filter := bson.M{}

	// Define the fields to project, only include `cid` field
	projection := bson.M{"cid": 1, "_id": 0}

	// Find all documents with the given filter and projection
	cursor, err := impl.Collection.Find(ctx, filter, options.Find().SetProjection(projection))
	if err != nil {
		impl.Logger.Error("database find error", slog.Any("error", err))
		return nil, err
	}
	defer cursor.Close(ctx)

	var cids []string
	for cursor.Next(ctx) {
		var result struct {
			CID string `bson:"cid"`
		}
		if err := cursor.Decode(&result); err != nil {
			impl.Logger.Error("cursor decode error", slog.Any("error", err))
			return nil, err
		}
		// Append the `cid` value to the result slice
		cids = append(cids, result.CID)
	}

	// Check for errors during cursor iteration
	if err := cursor.Err(); err != nil {
		impl.Logger.Error("cursor iteration error", slog.Any("error", err))
		return nil, err
	}

	return cids, nil
}
