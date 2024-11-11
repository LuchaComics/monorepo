package datastore

import (
	"context"
	"log"
	"log/slog"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

func (impl NFTAssetStorerImpl) DeleteByID(ctx context.Context, id primitive.ObjectID) error {
	_, err := impl.Collection.DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		log.Fatal("DeleteOne() ERROR:", err)
	}
	return nil
}

func (impl NFTAssetStorerImpl) DeleteByCID(ctx context.Context, cid primitive.ObjectID) error {
	_, err := impl.Collection.DeleteOne(ctx, bson.M{"cid": cid})
	if err != nil {
		impl.Logger.Error("database failed deletion error",
			slog.Any("error", err))
		return err
	}
	return nil
}
