package main

import (
	"github.com/xyths/game-cashier/client/cryptolions"
	"github.com/xyths/game-cashier/cmd/utils"
	"github.com/xyths/game-cashier/node"
	"github.com/xyths/game-cashier/puller"
	"github.com/xyths/game-cashier/types"
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
			utils.DryRunFlag,
		},
	}
	notifyCommand = &cli.Command{
		Action:  notify,
		Name:    "notify",
		Aliases: []string{"n"},
		Usage:   "pull the recharge history to database.",
		Flags: []cli.Flag{
			utils.DryRunFlag,
		},
	}
	downloadCommand = &cli.Command{
		Action:  download,
		Name:    "download",
		Aliases: []string{"d"},
		Usage:   "manually download the recharge history to database.",
		Flags: []cli.Flag{
			utils.AfterFlag,
			utils.BeforeFlag,
			utils.DryRunFlag,
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
	if config.Server.ServerType == "cryptolions" {
		return pullFromCryptolions(ctx, config)
	} else if config.Server.ServerType == "dfuse" {
		return pullFromDfuse(ctx, config)
	} else {
		log.Printf("unknown server type: %s", config.Server.ServerType)
		return nil
	}
}

func pullFromCryptolions(ctx *cli.Context, config utils.Config) error {
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
	defer client.Disconnect(ctx.Context)
	if err != nil {
		log.Fatal("Error when connect to mongo:", err)
	}
	// Check the connection
	err = client.Ping(ctx.Context, nil)
	if err != nil {
		log.Fatal("Error when ping to mongo:", err)
	}
	coll := client.Database(config.Mongo.Database).Collection(types.TransferCollName)
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
			for _, r := range records {
				// 单条插入，不管成功与否，都会继续下一条
				if _, err = coll.InsertOne(ctx.Context, r); err != nil {
					log.Print(err)
				}
			}

			log.Println("Pull finished.")
		}
	}

	return nil

}

func pullFromDfuse(ctx *cli.Context, config utils.Config) error {
	clientOpts := options.Client().ApplyURI(config.Mongo.URI)
	client, err := mongo.Connect(ctx.Context, clientOpts)
	defer client.Disconnect(ctx.Context)
	if err != nil {
		log.Fatal("Error when connect to mongo:", err)
	}
	// Check the connection
	err = client.Ping(ctx.Context, nil)
	if err != nil {
		log.Fatal("Error when ping to mongo:", err)
	}
	db := client.Database(config.Mongo.Database)
	p := puller.Puller{}
	if err := p.Init(config.Network, config.Server.Dfuse.ApiKey, config.Server.Dfuse.Manager, db); err != nil {
		log.Fatalf("error when init puller object: %s", err)
	}

	return p.Pull(ctx)
}

func notify(ctx *cli.Context) error {
	log.Println("notify started")
	n := node.Node{}
	n.Init(ctx.Context, ctx.String(utils.ConfigFlag.Name))

	return n.Notify(ctx.Context)
}

func download(ctx *cli.Context) error {
	log.Println("manually download started")
	config, err := utils.ParseConfig(ctx.String(utils.ConfigFlag.Name))
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("network: %s, api server: %s", config.Network, config.Server.ServerType)
	after := ctx.String(utils.AfterFlag.Name)
	before := ctx.String(utils.BeforeFlag.Name)
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
	records, err := cl.Pull(ctx.Context, after, before)
	after = before
	if err != nil {
		log.Printf("Error when upload: %s", err)
	}
	if len(records) == 0 {
		return nil
	}
	for _, r := range records {
		// 单条插入，不管成功与否，都会继续下一条
		if _, err = coll.InsertOne(ctx.Context, r); err != nil {
			log.Print(err)
		}
	}

	log.Println("manually download finished")
	return nil
}
