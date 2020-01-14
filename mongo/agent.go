package mongo

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
)

type MongoAgent struct {
	db *mongo.Database
}

const transferColl = "transfer"

func (ma *MongoAgent) insertOne(ctx context.Context, coll string, record interface{}) error {
	collection := ma.db.Collection(coll)

	_, err := collection.InsertOne(ctx, record)
	if err != nil {
		log.Println(err)
	}
	return err
}
func (ma *MongoAgent) insertMany(ctx context.Context, coll string, documents []interface{}) error {
	collection := ma.db.Collection(coll)

	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (ma *MongoAgent) InsertTransferRecord(ctx context.Context, record interface{}) error {
	return ma.insertOne(ctx, transferColl, record)
}
