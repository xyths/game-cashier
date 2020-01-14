package main

import (
	"github.com/xyths/game-cashier/client/cryptolions"
	"github.com/xyths/game-cashier/cmd/utils"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/urfave/cli.v2"
	"log"
	"time"
)

var (
	serveCommand = &cli.Command{
		Action:  serve,
		Name:    "serve",
		Aliases: []string{"s"},
		Usage:   "serve the API for recharge history.",
		Flags: []cli.Flag{
		},
	}
	pullCommand = &cli.Command{
		Action:  pull,
		Name:    "pull",
		Aliases: []string{"p"},
		Usage:   "pull the recharge history to database.",
		Flags: []cli.Flag{
		},
	}
	notifyCommand = &cli.Command{
		Action:  notify,
		Name:    "notify",
		Aliases: []string{"n"},
		Usage:   "pull the recharge history to database.",
		Flags: []cli.Flag{
		},
	}
)

func pull(ctx *cli.Context) error {
	log.Println("pull started")
	config, err := utils.ParseConfig(ctx.String(utils.ConfigFlag.Name))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("network: %s, api server: %s", config.Network, config.Server.ServerType)
	d, err := time.ParseDuration(config.Interval)
	if err != nil {
		log.Fatal(err)
	}
	after := time.Now().Format("2006-01-02T15:04:05.000-07:00")
	cl := cryptolions.Cryptolions{
		Server:  config.Server.CryptoLions.Server,
		Manager: config.Manager,
	}
	clientOpts := options.Client().ApplyURI(config.Mongo.URI)
	client, err := mongo.Connect(ctx.Context, clientOpts)
	if err != nil {
		log.Fatal("Error when connect to mongo:", err)
	}
	// Check the connection
	err = client.Ping(ctx.Context, nil)
	if err != nil {
		log.Fatal("Error when ping to mongo:", err)
	}
	coll := client.Database(config.Mongo.Database).Collection("transfer")
	for {
		select {
		case <-ctx.Context.Done():
			log.Println(ctx.Context.Err())
			return nil
		case <-time.After(d):
			before := time.Now().Format("2006-01-02T15:04:05.000-07:00")
			log.Printf("Pull record from %s to %s", after, before)
			records, err := cl.Pull(ctx.Context, after, before)
			after = before
			if err != nil {
				log.Printf("Error when upload: %s", err)
			}
			if len(records) == 0 {
				continue
			}
			docs := []interface{}{}
			for _, r := range records {
				docs = append(docs, r)
			}

			_, err = coll.InsertMany(ctx.Context, docs)
			if err != nil {
				log.Printf("error when insertMany: %s", err)
				continue
			}

			log.Println("Pull finished.")
		}
	}
	//log.Printf("network: %s, api_key: %s, manager: %s",
	//	config.Server.Network, config.Server.ApiKey, config.Server.Manager)
	//p := new(puller.Puller)
	//
	//if err = p.Init(config.Server.Network, config.Server.ApiKey, config.Server.Manager); err != nil {
	//	log.Fatal(err)
	//}
	//if err = p.Pull(ctx.Context); err != nil {
	//	log.Fatal(err)
	//}
	return nil
}

func notify(ctx *cli.Context) error {
	log.Println("notify started")
	return nil
}
