package mongodb

import (
	"context"
	"log"

	"log/slog"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"

	c "github.com/LuchaComics/monorepo/cloud/comiccoin-authority/config"
)

func NewProvider(appCfg *c.Configuration, logger *slog.Logger) *mongo.Client {
	logger.Debug("storage initializing...")
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(appCfg.DB.URI))
	if err != nil {
		log.Fatal(err)
	}

	// The MongoDB client provides a Ping() method to tell you if a MongoDB database has been found and connected.
	if err := client.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}
	logger.Debug("storage initialized successfully")
	return client
}
