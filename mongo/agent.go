package agent

import (
	"context"
	"github.com/xyths/game-cashier/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"log"
	"time"
)

type MongoAgent struct {
	Db *mongo.Database
}

const transferColl = "transfer"

func (ma *MongoAgent) insertOne(ctx context.Context, coll string, record interface{}) error {
	collection := ma.Db.Collection(coll)

	_, err := collection.InsertOne(ctx, record)
	if err != nil {
		log.Println(err)
	}
	return err
}
func (ma *MongoAgent) insertMany(ctx context.Context, coll string, documents []interface{}) error {
	collection := ma.Db.Collection(coll)

	_, err := collection.InsertMany(ctx, documents)
	if err != nil {
		log.Println(err)
	}
	return err
}

func (ma *MongoAgent) InsertTransferRecord(ctx context.Context, record interface{}) error {
	return ma.insertOne(ctx, transferColl, record)
}

type History struct {
	Tx          string  `json:"txHash" bson:"tx"`
	BlockNumber uint64  `bson:"blockNumer"`
	From        string  `bson:"from"`
	To          string  `bson:"to"`
	Amount      float64 `bson:"amount"`
	Timestamp   string  `bson:"timestamp"`
}

func (ma *MongoAgent) GetHistory(ctx context.Context, memo, start, end string) (records []History, err error) {
	log.Printf("find record from collection %s, memo is %s, time from %s to %s", transferColl, memo, start, end)
	if start != "" {
	}
	if end != "" {

	}
	coll := ma.Db.Collection(transferColl)
	if memo == "" {
		return
	}
	filter := bson.D{
		{"memo", memo},
	}
	hasTime := false
	t := bson.D{}
	if start != "" {
		startDate, err := time.Parse(time.RFC3339, start)
		if err != nil {
			return records, err //inter err
		}
		t = append(t, bson.E{"$gte", startDate})
		hasTime = true
	}
	if end != "" {
		endDate, err := time.Parse(time.RFC3339, end)
		if err != nil {
			return records, err //inter err
		}
		t = append(t, bson.E{"$lte", endDate})
		hasTime = true
	}
	if hasTime {
		filter = append(filter, bson.E{
			"txTime", t,
		})
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return
	}

	if err = cursor.All(ctx, &records); err != nil {
		return
	}
	return
}

func (ma *MongoAgent) NotNotified(ctx context.Context) (records []types.TransferRecord, err error) {
	log.Printf("find record not notified yet from collection %s", transferColl)
	coll := ma.Db.Collection(transferColl)
	filter := bson.D{
		{"notifyTime", ""},
	}

	cursor, err := coll.Find(ctx, filter)
	if err != nil {
		return
	}

	if err = cursor.All(ctx, &records); err != nil {
		return
	}
	return
}

func (ma *MongoAgent) UpdateNotifyTime(ctx context.Context, record types.TransferRecord) error {
	coll := ma.Db.Collection(transferColl)
	_, err := coll.UpdateOne(ctx,
		bson.D{
			{"_id", record.Id},
		},
		bson.D{
			{"$currentDate", bson.D{
				{"notifyTime", true},
			}},
		})

	return err
}
