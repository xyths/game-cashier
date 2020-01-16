package node

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/xyths/game-cashier/cmd/utils"
	"github.com/xyths/game-cashier/mongo"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
)

type Node struct {
	db *mongo.Database

	Addr   string
	Engine *gin.Engine
}

func (n *Node) Init(ctx context.Context, configFilename string) {
	c, err := utils.ParseConfig(configFilename)
	if err != nil {
		log.Fatalf("config format error: %s", err)
	}
	n.Addr = c.Listen
	n.initMongo(ctx, c)
	n.initEngine(ctx)
}

func (n *Node) initMongo(ctx context.Context, config utils.Config) {
	clientOpts := options.Client().ApplyURI(config.Mongo.URI)
	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal("Error when connect to mongo:", err)
	}
	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Error when ping to mongo:", err)
	}
	n.db = client.Database(config.Mongo.Database)
}

func (n *Node) initEngine(ctx context.Context) {
	r := gin.Default()

	// Ping test
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Group using gin.BasicAuth() middleware
	// gin.Accounts is a shortcut for map[string]string
	authorized := r.Group("/api/v1", gin.BasicAuth(gin.Accounts{
		"sgseos": "20200116",
	}))
	{
		authorized.GET("/history/:memo", n.getHistory)
	}

	n.Engine = r
	log.Println("Set HTTP router.")
}

func (n *Node) getHistory(c *gin.Context) {
	memo := c.Param("memo")
	start := c.Query("start")
	end := c.Query("end")
	log.Printf("get history for memo %s, time from %s to %s", memo, start, end)

	ma := agent.MongoAgent{Db: n.db}
	if records, err := ma.GetHistory(c, memo, start, end); err == nil {
		c.JSON(http.StatusOK, records)
	} else {
		c.JSON(http.StatusInternalServerError, gin.H{})
	}
}
