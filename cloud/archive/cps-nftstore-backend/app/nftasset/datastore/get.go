package datastore

import (
	"context"

	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (impl NFTAssetStorerImpl) GetByID(ctx context.Context, id primitive.ObjectID) (*NFTAsset, error) {
	filter := bson.M{"_id": id}

	var result NFTAsset
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

func (impl NFTAssetStorerImpl) GetByCID(ctx context.Context, cid string) (*NFTAsset, error) {
	filter := bson.M{"cid": cid}

	var result NFTAsset
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

func (impl NFTAssetStorerImpl) GetAllCIDs(ctx context.Context) ([]string, error) {
	// Define an empty filter to retrieve all documents
	filter := bson.M{}

	// Define the fields to nftmetadata, only include `cid` field
	nftmetadataion := bson.M{"cid": 1, "_id": 0}

	// Find all documents with the given filter and nftmetadataion
	cursor, err := impl.Collection.Find(ctx, filter, options.Find().SetProjection(nftmetadataion))
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
