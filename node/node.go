package node

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/xyths/game-cashier/cmd/utils"
	"github.com/xyths/game-cashier/mongo"
	"github.com/xyths/game-cashier/types"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"net/http"
	"time"
)

type Node struct {
	db *mongo.Database

	Addr   string
	Engine *gin.Engine

	duration  time.Duration
	network   string
	notifyUrl string
}

func (n *Node) Init(ctx context.Context, configFilename string) {
	c, err := utils.ParseConfig(configFilename)
	if err != nil {
		log.Fatalf("config format error: %s", err)
	}
	n.Addr = c.Listen
	n.network = c.Network
	n.notifyUrl = c.Notify
	n.duration, err = time.ParseDuration(c.Interval)
	if err != nil {
		log.Fatal(err)
	}

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
	// RFC3339     = "2006-01-02T15:04:05Z07:00"
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

func (n *Node) Notify(ctx context.Context) error {
	ma := agent.MongoAgent{Db: n.db}
	for {
		select {
		case <-ctx.Done():
			log.Println(ctx.Err())
			return nil
		case <-time.After(n.duration):
			records, err := ma.NotNotified(ctx)
			if err != nil {
				log.Printf("error when get record-not-notified: %s", err)
				break // break select, continue for
			}
			for _, r := range records {
				success, err := n.NotifyOne(ctx, r)
				if err != nil || !success {
					log.Printf("notify error: %s, success: %v", err, success)
					continue

				}

				if err = n.UpdateNotify(ctx, r); err != nil {
					log.Printf("update notify error: %s", err)
					continue
				}
			}
		}
	}
}

func (n *Node) NotifyOne(ctx context.Context, r types.TransferRecord) (success bool, err error) {
	ne := types.NotifyElement{
		Network: n.network,
		Memo:    r.Memo,
		Amount:  r.Amount,
		Tx:      r.Tx,
	}
	jsonValue, _ := json.Marshal(ne)

	resp, err := http.Post(n.notifyUrl, "application/json", bytes.NewBuffer(jsonValue))
	if err != nil {
		log.Printf("error when notify: %s", err)
		return
	}
	defer resp.Body.Close()

	log.Println("Response status:", resp.Status)

	scanner := bufio.NewScanner(resp.Body)
	for i := 0; scanner.Scan() && i < 5; i++ {
		log.Println(scanner.Text())
	}
	err = scanner.Err();
	if err != nil {
		log.Println(err)
		return
	}

	return resp.StatusCode == http.StatusOK, nil
}

func (n *Node) UpdateNotify(ctx context.Context, r types.TransferRecord) error {
	ma := agent.MongoAgent{Db: n.db}
	return ma.UpdateNotifyTime(ctx, r)
}
