package db

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type DB struct {
	client *mongo.Client
}

func Connect(ctx context.Context, uri string) (*DB, error) {
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		return nil, err
	}

	return &DB{client: client}, nil
}

func (d *DB) Client() *mongo.Client {
	return d.client
}

func (d *DB) Close(ctx context.Context) error {
	return d.client.Disconnect(ctx)
}
