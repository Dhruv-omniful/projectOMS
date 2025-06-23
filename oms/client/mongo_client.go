package client

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/omniful/go_commons/config"
	"github.com/omniful/go_commons/log"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/dhruv/oms/model"
)

var mongoLogger = log.DefaultLogger()

// GetMongoClient returns a Mongo client with config values
func GetMongoClient(ctx context.Context) (*mongo.Client, error) {
	uri := config.GetString(ctx, "mongodb.uri")
	opts := options.Client().ApplyURI(uri)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		mongoLogger.Errorf("❌ Mongo Connect error: %v", err)
		return nil, err
	}

	if err := client.Ping(ctx, nil); err != nil {
		mongoLogger.Errorf("❌ Mongo Ping error: %v", err)
		return nil, err
	}

	return client, nil
}

// GetOrdersCollection returns the orders collection
func GetOrdersCollection(ctx context.Context) (*mongo.Collection, error) {
	client, err := GetMongoClient(ctx)
	if err != nil {
		return nil, err
	}

	dbName := config.GetString(ctx, "mongodb.database")
	return client.Database(dbName).Collection("orders"), nil
}

// SaveOrder inserts an order into MongoDB
func SaveOrder(ctx context.Context, o *model.Order) error {
	coll, err := GetOrdersCollection(ctx)
	if err != nil {
		mongoLogger.Errorf("❌ Failed to get orders collection: %v", err)
		return err
	}

	o.ID = uuid.NewString()
	o.Status = "on_hold"
	o.CreatedAt = time.Now().UTC()

	_, err = coll.InsertOne(ctx, o)
	if err != nil {
		mongoLogger.Errorf("❌ Mongo InsertOne error: %v", err)
		return err
	}

	mongoLogger.Infof("✅ Order saved: %+v", o)
	return nil
}
